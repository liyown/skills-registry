# Observability Coverage Audit Skill

`observability-coverage-audit` finds which code paths, error states, or failure modes are exercised by tests, monitors, alerts, or runbooks — and which are silent. The output is a per-site coverage matrix with gap rows.

## What It Audits

- A scoped code path or failure class.
- The existing test coverage per site.
- The existing alert coverage per site.
- The existing runbook coverage per site.
- The classification of the gap (silent / no-runbook / no-alert / no-test).

## Discovery Tools

- **`codegraph_explore`** — one-shot fanout.
- **`codegraph_search`** — locate the error / exception / throw.
- **`codegraph_callers`** — every site that can hit the error.
- **`codegraph_files`** — test file filter.
- **Fallback** — per-feature source inspection.

## Output Contract

A per-site coverage matrix: site (path + line), test coverage, alert coverage, runbook coverage, suggested fix.

## Related Skills

- **`spec-doc-linter`** — keeps `DevAgent.md` / `CONTEXT.md` honest.
- **`incident-call-trace`** — walks the failing path; this skill is the post-incident counterpart.
- **`data-lineage-trace`** — the upstream input.
- **`goal-driven-development`** — orchestrates the spec-to-code flow; this skill is the post-incident / pre-launch counterpart.

## Files

```text
.
├── SKILL.md
├── README.md
├── prompts/
│   └── coverage-audit.md       # the per-site coverage framework
└── examples/
    ├── bad-coverage.md
    ├── good-coverage.md
    └── coverage-output.md     # canonical "what the agent should emit" sample
```
