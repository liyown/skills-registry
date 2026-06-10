# Add `spec-doc-linter` skill and drop scripts/CI gate

Date: 2026-06-10

## Goal

Add a 8th skill, `spec-doc-linter`, that keeps two project-local doc conventions — per-module `DevAgent.md` and per-domain `CONTEXT.md` — in sync with the code they describe, with per-file `y/n/q` auto-sync. In the same release, drop `scripts/{validate.sh, smoke.sh, check-examples.sh, release.sh}` and `.github/workflows/release-readiness.yml`; the repo is text only and the contributor contract is enforced by review against `CONTRIBUTING.md` rather than by a script.

## Context

The collection's `npx skills add` consumers install each `skills/<name>/` directory verbatim — there is no build step, no test runner, and no linter for the registry itself. The shell scripts under `scripts/` and the CI workflow were a "fail the PR if the contract is broken" guard, but every script edit was a friction point (the v0.3.10 release dropped a "redundant release mirror" step; the linter branch in `check-examples.sh` was the right kind of extension but turned out to need a parallel `check_linter_skill` function gated on `prompts/linter.md`, which the v0.4.0 work would have required).

Meanwhile, the user has long-standing feedback that "skill is text" — embedded scripts under `skills/` are out-of-shape. The linter skill is text-only: `rg` / `Read` / `Grep` recipes the agent executes, with a unified diff the agent shows before each per-file `y` / `n` / `q` decision.

## Key Entrypoints

- `skills/spec-doc-linter/SKILL.md` — the always-loaded router. ≤ 80 lines.
- `skills/spec-doc-linter/prompts/linter.md` — core protocol: doc conventions, drift catalog, sentinel.
- `skills/spec-doc-linter/prompts/drift-checks.md` — Tier-1 (mechanical) rules A.1–A.6.
- `skills/spec-doc-linter/prompts/semantic-review.md` — Tier-2 (LLM judgment) rules B.1–B.6.
- `skills/spec-doc-linter/prompts/confirmation.md` — per-file `y/n/q` loop with a hard no-`--fix-all` rule.
- `skills/spec-doc-linter/examples/linter-output.md` — canonical drift-report sample. Must contain the sentinel `未发现文档与代码漂移。`.

## CodeGraph Findings

CodeGraph unavailable; context was gathered by rg/file inspection.

`rg` queries used:

- `rg -l '^## (Overview|Public API|Invariants|Anti-patterns|Dependencies|File Map)$' --glob 'DevAgent.md'` (no hits — no consumer project has adopted the convention yet).
- `rg 'validate.sh|smoke.sh|check-examples.sh|release.sh' CONTRIBUTING.md README.md` (8 hits, all stripped).
- `rg 'prompts/linter.md|spec-doc-linter' skills/` (the discriminator key for the linter branch).

Files read: `skills/{java,react,go,python,node}-code-reviewer/SKILL.md` (router shape), `skills/goal-driven-development/SKILL.md` (workflow shape), `scripts/{validate,smoke,check-examples,release}.sh` (gate shape), `.github/workflows/release-readiness.yml` (CI shape), `docs/CHANGELOG.md` (release shape).

## Decisions

- **Two-tier drift detection.** Class A (mechanical) is static-analysis-detectable: renamed symbols, changed signatures, module/domain added or removed, quoted file paths, missing dependencies, broken intra-doc links. Class B (judgment) needs LLM reasoning: weakened invariants, anti-patterns newly violated, false idempotency / atomicity / thread-safety claims, broken boundaries, false concurrency claims, stale context maps. Tier-1 is fast and deterministic; Tier-2 is slower and probabilistic. Both tiers feed the same drift report.
- **New sentinel, not the reviewer's.** `未发现文档与代码漂移。` (linter) is deliberately distinct from `未发现明确高风险问题。` (reviewer). They have different meanings and different consumer sets; reusing the reviewer string would mislead anything that greps for it.
- **Per-file `y/n/q`, no `--fix-all`.** A `--fix-all` flag would let the linter mutate the working tree without a human pause — out of scope. The user explicitly chose "auto-sync with confirmation", which means every file is a separate decision.
- **Drop `scripts/` and the CI gate in the same release.** They were a structural enforcement of a contract that is now review-enforced. Splitting the two changes across two releases would have left the repo in a half-scripted state with no script, which is worse than either alone.
- **Keep the existing Reviewer Coverage Matrix unchanged.** A new `### Tools` sub-section in README "Skills at a Glance" is the least invasive way to expose the linter without reshaping the matrix (which is keyed by scenario × language; the linter has neither axis).

## Verification

- `gh pr view 1 --json state` returned `MERGED`.
- Plan file at `/Users/liuyaowen/.claude/plans/logical-dazzling-toast.md` was the design doc; the 5-commit plan was collapsed to 3 commits after the user deleted `validate.sh` mid-stream.
- Bad/good example pairs verified by hand: `bad-devagent.md` has 6 embedded drift findings (A.1, B.1, B.3, A.5×2, A.4); `good-devagent.md` addresses all 6 with `Fix N:` annotations. `bad-context.md` has 3 (B.4, B.6, A.6); `good-context.md` addresses all 3.
- Sentinel appears in `prompts/linter.md` and `examples/linter-output.md` (both required for downstream consumers that grep for the literal string).
- `git log --oneline v0.3.10..main` returns 4 commits: `2f6d837`, `604549f`, `143f404`, `6493868`.

## Review Conclusions

未发现明确高风险问题。

The skill is text-only and self-contained. The risk surface is the LLM-judgment tier, but the per-file confirmation gate is the catch: an LLM hallucination about an invariant cannot mutate the working tree without a human `y`.

## Follow-up Notes

- The two conventions (`DevAgent.md` per module, `CONTEXT.md` per domain) are not externally standardized; the linter *teaches* them. The first consumer project that adopts them will surface real drift modes that the current A.1–A.6 / B.1–B.6 catalog doesn't catch — extend `prompts/drift-checks.md` by appending new `### A.N` blocks, no other edit.
- A.5 (new dep not in doc) currently parses a fixed set of manifests (`package.json`, `Cargo.toml`, `go.mod`, `pyproject.toml`, `pom.xml`, `build.gradle`). Add the consumer's manifest if it's not on the list.
- The `NEEDS_CONTEXT` tier-2 verdict (B.1 / B.3 / B.5) is the right escape hatch for claims that cannot be settled by static analysis — keep it distinct from "definite drift". Anything tagged `NEEDS_CONTEXT` is not editable even with confirmation granted.
- If the user later wants `--fix-all`, the change is one line in `prompts/confirmation.md` plus a guard rail. For now the absence is intentional.
