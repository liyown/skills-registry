---
name: progressive-rollout-checklist
description: "Use before merging a change that ships behind a feature flag, kill switch, canary, dark launch, or progressive rollout. Triggers on: feature flag, kill switch, canary, dark launch, progressive delivery, rollout plan, 'how do I ship this safely', 'what's the rollback plan', 'is the flag wired up', 'who is in the canary cohort', sunset, deprecation, rollback. Do NOT use for: finishing a branch (use obra/superpowers finishing-a-development-branch), or for post-incident analysis (use incident-call-trace). Complementary to goal-driven-development, which orchestrates the spec-to-code flow; progressive-rollout-checklist is the pre-ship counterpart that verifies the flag, the cohort, the SLO gate, and the rollback path before the change goes live."
metadata:
  short-description: Rollout-readiness, flag-coverage, and rollback-path verification
---

# Progressive Rollout Checklist

Before merging a change that ships behind a feature flag, kill switch, canary, dark launch, or progressive rollout, verify the flag is wired up, the cohort is defined, the SLO gate exists, and the rollback path is in place. The output is a rollout packet that refuses to mark the change "ready to ship" until every required field is filled.

## Required Loading

Always load:

- `prompts/rollout-check.md` — the flag / cohort / SLO-gate / rollback framework; the per-flag-row checklist.

## When To Run

- Before merging a change that introduces or flips a feature flag.
- Before flipping a flag's default.
- Before expanding a canary cohort.
- Before sunsetting a public endpoint.
- Before any production rollout where a kill switch is required.

## Discovery Order

1. Identify the flag(s) involved (`rg '<flag-name>'` over the codebase, or `codegraph_search`).
2. For each flag, find every read site and every write site. A flag is "wired up" only if every read site has a default value.
3. For each flag, find the cohort definition (the rule that determines who sees the new behaviour).
4. For each flag, find the SLO gate (the metric that gates the cohort expansion).
5. For each flag, find the kill switch (the action that disables the flag globally).
6. For each flag, find the rollback hook (the action that reverts the change).
7. For each flag, find the sunset date (if the flag is being deprecated).

## Output Contract

The rollout packet is grouped by flag. For each flag, name:

- the flag name
- the read sites and write sites (with path + line)
- the default value and the proposed new default
- the cohort definition
- the SLO gate (which metric, which threshold)
- the kill switch (the action, the on-call who holds it)
- the rollback hook (the action, the owner)
- the sunset date (if any)

A flag that is missing any of these is itself a `ROLL-FLAG-INCOMPLETE` finding. The rollout cannot proceed until the missing field is filled.

## Tools

- **`codegraph_search`** — locate a flag by name.
- **`codegraph_callers` / `codegraph_callees`** — every read and write of the flag.
- **`codegraph_impact`** — radius if the default flips.
- **`codegraph_files`** — route-to-cohort mapping.
- **Fallback** — `rg '<flag-name>'` + per-path `Read` + manual config review.

## Fallback

If CodeGraph is unavailable, the fallback is `rg` over flag names + manual config review. The final report must include the line:

```text
CodeGraph unavailable; rollout gathered by rg/file inspection.
```

## Related Skills

- **`spec-doc-linter`** — keeps `DevAgent.md` / `CONTEXT.md` honest.
- **`code-ownership-impact`** — names the per-team approvers for a cross-team change; this skill adds the rollout-time checks.
- **`goal-driven-development`** — orchestrates the spec-to-code flow; this skill is the pre-ship counterpart.

## Examples

Each `bad-*` doc has a matching `good-*` in `examples/` showing the minimum packet. Read them side by side to calibrate packet depth. `examples/rollout-output.md` is the canonical "what the agent should emit" sample.
