# Spec Doc Linter Skill

`spec-doc-linter` keeps two project-local doc files in sync with the code they describe: per-module `DevAgent.md` and per-domain `CONTEXT.md`. It is a linter, not a reviewer — it surfaces drift and proposes fixes, then writes only after a per-file confirmation.

## What It Lints

- Mechanical drift (Tier 1, static analysis):
  - Renamed symbol: doc references `Foo.bar` but the code no longer has `Foo` or `bar`.
  - Changed signature: doc shows a function signature that no longer matches the code.
  - Module / domain added or removed without a matching `DevAgent.md` / `CONTEXT.md`, or a doc file whose module is gone.
  - Quoted file path no longer exists (`./src/foo/bar.rs` referenced in the doc is gone).
  - New dependency not reflected in the doc's `Dependencies` section.
  - Broken intra-doc links (relative Markdown link points nowhere).
- Judgment drift (Tier 2, LLM):
  - Weakened invariant: doc claims "idempotent" but the code mutates state on second call.
  - Anti-pattern newly violated: doc forbids `import internal/foo` but the code imports it.
  - False idempotency / atomicity / thread-safety claim.
  - Boundary crossed: doc says "domain A must not depend on domain B" but a new import exists.
  - Concurrency claim false (e.g. doc says "thread-safe" but code holds a non-`Send` type across an `.await`).
  - Context map out of date: doc lists dependent contexts but a new one was added.

## Prompt Loading

Always load `prompts/linter.md`, `prompts/drift-checks.md`, `prompts/semantic-review.md`, and `prompts/confirmation.md`. There is no scenario-specific loading; the four are co-required so the agent has the static rules, the LLM-judgment rubric, the output contract, and the per-file confirmation flow in one pass.

## Output Contract

If no drift is found, output exactly:

```text
未发现文档与代码漂移。
```

When drift is found, output a per-file drift report (see `examples/linter-output.md`):

```text
# Drift Report

## <relative/path/to/devagent.md>
- <file:line>  <DriftModeId>  <evidence>  →  <suggested fix>
- ...

## <relative/path/to/context.md>
- ...
```

Mutations are gated by `prompts/confirmation.md`: the agent shows the proposed diff for a single file, asks `Apply? [y/n/q]`, then either writes the patch (`y`), leaves the file untouched and moves on (`n`), or stops the run (`q`). A fix-all mode is intentionally not provided.

## Files

```text
.
├── SKILL.md
├── README.md
├── prompts/
│   ├── linter.md             # core protocol; doc conventions + drift catalog
│   ├── drift-checks.md       # Tier-1 static check rules (extendable)
│   ├── semantic-review.md    # Tier-2 LLM-judgment rubric
│   └── confirmation.md       # per-file diff display + y/n/q contract
└── examples/
    ├── bad-devagent.md
    ├── good-devagent.md
    ├── bad-context.md
    ├── good-context.md
    ├── linter-output.md      # canonical drift report sample
    └── pr-diff-example.diff
```
