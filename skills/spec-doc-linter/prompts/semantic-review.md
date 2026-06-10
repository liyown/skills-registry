# Tier-2 Semantic Review

> See also: prompts/linter.md, prompts/drift-checks.md

Tier-2 is LLM-judgment, not regex. For each Class B mode, the agent reads the doc claim, reads the code, and reasons about whether the claim is still true.

## Prompt Skeleton

For each Class B candidate, the agent should construct its own prompt along these lines and answer in-line:

```text
Claim (from doc):    "<verbatim line from doc>"
Location:            <file:line in doc>
Code evidence:       <file:line in code, with a short excerpt>
Question:            Is the claim still true given the code?
Verdict:             YES | NO | NEEDS_CONTEXT
If NO:               suggested fix
```

`NEEDS_CONTEXT` is reserved for cases the LLM cannot decide without runtime evidence (e.g. true thread-safety under load). The final report must list every `NEEDS_CONTEXT` and explicitly ask the user to confirm.

## B.1 Weakened invariant

- Claim shape: any bullet in `## Invariants` or `## Critical Invariants` of the form "always X", "must X", "X is Y".
- Reasoning: re-read the referenced function. Apply it twice in dry-run; does state change between calls? Does an exception path skip the invariant's precondition?
- Emits: `B.1`, evidence "doc claims `deduct` is idempotent; second call decrements again", fix "either implement the idempotency guard or weaken the claim."

## B.2 Anti-pattern newly violated

- Claim shape: bullets in `## Anti-patterns` listing things to avoid.
- Reasoning: for each anti-pattern, `rg` for the forbidden identifier. If a hit exists, emit drift.
- Emits: `B.2`, evidence "doc forbids `import internal/foo`; `internal/foo.go` is imported in `service/handler.go:14`", fix "either remove the import or remove the anti-pattern entry."

## B.3 False idempotency / atomicity / thread-safety

- Claim shape: any `## Invariants` line containing "idempotent", "atomic", "thread-safe", "linearizable", "exactly once".
- Reasoning: the LLM must construct a counter-example call pattern in plain English. If the construction succeeds without contradiction, the claim is false.
- Emits: `B.3`, evidence "doc claims `cancel` is idempotent; calling it on a `CANCELLED` order raises and surfaces as HTTP 500", fix "return success on no-op or weaken the claim."

## B.4 Boundary crossed

- Claim shape: a `## Does Not Own` line that names another domain, paired with an `## Owns` line that excludes it.
- Reasoning: enumerate imports from the scope root. Any import from a forbidden domain → drift.
- Emits: `B.4`, evidence "CONTEXT.md says `billing` does not depend on `shipping`; `billing/` imports `shipping.RateQuote`", fix "either drop the import or rewrite CONTEXT.md to acknowledge the new dependency."

## B.5 Concurrency claim false

- Claim shape: a `## Invariants` line containing "thread-safe", "concurrent", "lock-free", "Send", "Sync".
- Reasoning: for the named type, identify any field whose type is `Rc`, `RefCell`, `MutexGuard<'a, T>` with `'a` non-`'static`, raw pointer, or any other non-`Send`/`Sync` shape. If the type is held across an `.await` or `tokio::spawn` boundary, emit drift.
- Emits: `B.5`, evidence "doc claims `RateLimiter` is `Send`; it holds `RefCell<...>`", fix "switch to `Mutex`/`RwLock` or weaken the claim."

## B.6 Context map out of date

- Claim shape: `## Upstream Contexts` and `## Downstream Contexts` lists.
- Reasoning: enumerate imports from the scope root, group by domain folder of the import target, and compare to the lists.
- Emits: `B.6`, evidence "CONTEXT.md lists `pricing` as a downstream; `billing/` also imports `inventory/Stock`", fix "add `inventory` to Downstream Contexts."

## Calibration

- The agent must distinguish "doc is wrong" from "code is wrong". Drift mode names the *doc* as wrong; the suggested fix is to the doc, not the code. If the code is the actual bug, surface it as a separate `B.x` finding and add a note `code may also be wrong; defer to human review`.
- Findings tagged `NEEDS_CONTEXT` are not auto-synced even if confirmation is granted. They are listed under a separate `## Needs Context` section in the report.
