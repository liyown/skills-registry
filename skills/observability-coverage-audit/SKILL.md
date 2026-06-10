---
name: observability-coverage-audit
description: "Use to find which code paths, error states, or failure modes are exercised by tests, monitors, alerts, or runbooks — and which are silent. Triggers on: blind spot, missing alert, 'do we have a test for this error path', 'do we alert on this', 'what happens if X fails', 'where is the runbook for Y', coverage audit, dead-letter, timeout path, partial failure, gap analysis. Do NOT use for: writing a test for a known behaviour (use obra/superpowers test-driven-development), or for root-cause analysis of a known symptom (use obra/superpowers systematic-debugging). Complementary to incident-call-trace, which walks the failing call path during an incident; observability-coverage-audit is the post-incident 'could this have been detected?' audit, the pre-launch 'is this new code path observable?' check, and the resilience-engineering gap report."
metadata:
  short-description: Gap audit: which error paths lack tests, alerts, or runbooks
---

# Observability Coverage Audit

Find which code paths, error states, or failure modes are exercised by tests, monitors, alerts, or runbooks — and which are silent. The output is a per-site coverage matrix with gap rows for the uncovered sites.

## Required Loading

Always load:

- `prompts/coverage-audit.md` — the per-site coverage framework; the per-control matrix.

## When To Run

- After an incident, before closing the postmortem: "could this have been detected earlier?".
- Before launching a new code path: "is this path observable end to end?".
- Quarterly as part of the resilience-engineering review.
- When a SLO breach keeps recurring despite the obvious fix: "is the test that 'fixes' it actually testing the production path?".

## Discovery Order

1. Identify the in-scope code path or failure class. Examples: the `payments.charge` path; every `throw` in the service; every external HTTP call.
2. For each site, find the existing coverage: test files (`*_test.go`, `*.spec.ts`, etc.), alert rules (in the metrics backend), runbook entries (in the runbook repo), dashboards.
3. For each site, ask: does the existing coverage actually exercise the production path? A test that mocks the dependency is not coverage of the dependency.
4. For each gap, classify: missing test, missing alert, missing runbook, missing all three. The classification is the suggested fix.

## Output Contract

The coverage matrix is grouped by site. For each site, name:

- the site (path + line)
- the test coverage (existing test name, or `MISSING`)
- the alert coverage (existing alert name, or `MISSING`)
- the runbook coverage (existing runbook path, or `MISSING`)
- the suggested fix (which coverage to add)

A site missing all three is a `COV-SILENT` finding. A site missing only the runbook is a `COV-NO-RUNBOOK` finding. The classification is the suggested-fix priority.

## Tools

- **`codegraph_explore`** — one-shot fanout when the in-scope path is a feature.
- **`codegraph_search`** — locate the error / exception / throw site by name.
- **`codegraph_callers`** — every site that can hit the error.
- **`codegraph_files`** — test file filter.
- **Fallback** — per-feature source inspection + manual test/alert/runbook search.

## Fallback

If CodeGraph is unavailable, the fallback is per-feature source inspection. The final report must include the line:

```text
CodeGraph unavailable; coverage gathered by rg/file inspection.
```

## Related Skills

- **`spec-doc-linter`** — keeps `DevAgent.md` / `CONTEXT.md` honest.
- **`incident-call-trace`** — walks the failing path during an incident; this skill is the post-incident 'could this have been detected?' audit.
- **`data-lineage-trace`** — the upstream input: a coverage gap without a data-lineage report is missing the sink list.
- **`goal-driven-development`** — orchestrates the spec-to-code flow; this skill is the post-incident / pre-launch counterpart.

## Examples

Each `bad-*` doc has a matching `good-*` in `examples/` showing the minimum matrix. Read them side by side to calibrate depth. `examples/coverage-output.md` is the canonical "what the agent should emit" sample.
