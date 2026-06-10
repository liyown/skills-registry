# RFC: Enterprise-Dev-Flow Skills (v0.5.0 follow-up)

Date: 2026-06-10
Status: Proposed (not yet scheduled for a release)
Author: skills-registry maintainer

## Context

The local collection ships 8 skills: 5 language-specific
production-risk code reviewers, `goal-driven-development` (spec →
code), `project-knowledge-capture`, and `spec-doc-linter` (keeps
`DevAgent.md` / `CONTEXT.md` honest). The adjacent collection
`obra/superpowers` ships 14 generic-author-craft skills
(`brainstorming`, `writing-plans`, `executing-plans`,
`test-driven-development`, `systematic-debugging`,
`verification-before-completion`, `requesting-code-review`,
`receiving-code-review`, `finishing-a-development-branch`,
`using-git-worktrees`, `using-superpowers`, `writing-skills`,
plus the subagent-driven / dispatching pair).

The two collections do not collide on skill names. They are
**complementary**, not competing: superpowers covers
generic-author-craft, the local collection covers
language-specific production-risk review + post-implementation
knowledge capture. The two have a gap, however, on the
**enterprise-specific** layer — the org / multi-team / regulated /
on-call / change-governance concerns that superpowers
deliberately refuses to assume and that the local collection
does not yet cover.

This RFC proposes **8 skills** that close that gap. None of
them duplicates a superpowers skill; each names the superpowers
skill it sits next to and the boundary it draws. None of the
8 ships in v0.5.0; the v0.5.0 release ships only the thin
orchestrator (`goal-driven-development`) reworked to point at
this RFC, plus a fix to a stale `codegraph.md` prompt. The 8
skills are scheduled for v0.5.1+.

## Constraints

1. **No overlap with superpowers.** Every proposed skill must
   name a superpowers skill it sits next to and a clear boundary.
2. **CodeGraph is a tool, not a design premise.** Some of the 8
   skills benefit from `@colbymchenry/codegraph` MCP/CLI (a
   `codegraph_explore`-driven call walk is cheaper than `rg` +
   file read for blast-radius work), but every skill has a
   `rg` + file-read fallback path so a consumer without
   CodeGraph installed is not blocked.
3. **Skill is text.** No executable scripts under any new
   skill; the consumer's agent uses `rg` / `Read` / `Grep` /
   `codegraph` MCP directly.
4. **Router-shaped description.** Each frontmatter `description`
   starts with the user-facing trigger, lists the technologies
   or context the skill covers, and includes the verbs a user
   would type. `metadata.short-description` is a 3-6 word noun
   phrase.

## Cross-reference to superpowers

For the 8 skills below, the closest superpowers skill and the
boundary is named in each section. Quick map of the gaps:

| Enterprise concern | Closest superpowers skill | Real gap? |
|---|---|---|
| Multi-team coordination (RFC, CODEOWNERS, shared kernel) | `writing-plans` (single-team plan) | Yes — no CODEOWNERS, no service boundary, no cross-team approver concept |
| Compliance & audit (SOX / ISO / GDPR / HIPAA / PCI) | `verification-before-completion` (test green?) | Yes — verification runs tests; compliance is a *production* claim |
| Change governance (feature flags, canary, sunset) | `finishing-a-development-branch` (merge) | Yes — ends at PR; does not know flags, cohorts, rollback |
| Reliability (SLO / on-call / postmortem / chaos) | `systematic-debugging` (local bug) | Yes — `systematic-debugging` is local; no SLO, no on-call, no postmortem |
| Security engineering (threat model, SBOM, CVE) | `requesting-code-review` (one-team review) | Yes — no threat model, no CVE / SBOM workflow |
| Data governance (PII / lineage / retention) | None | Yes — superpowers has no notion of a *field* as a first-class object |
| Tech debt & cost (license, dependency EOL) | `writing-plans` (deps) | Yes — no cost / license / EOL skill |
| Org context (CODEOWNERS, reorg-safe) | `using-git-worktrees` (workspace) | Yes — worktrees are technical, not org |
| AI governance (model card, prompt-injection) | None | Yes — no LLM-feature governance |
| PR / 集成前准备 (author-side) | `requesting-code-review` (English generic) | Adjacent — local skill is 中文 + 企业 CI + 多 reviewer 列表 |

The 8 proposed skills cover the rows above. Each is laid out
below with its `description`, archetype, primary tools,
superpowers boundary, and 2-3 enterprise scenarios.

---

## 1. `data-lineage-trace` — Lineage archetype

```yaml
description: "Use when you need to prove where a piece of data (PII, PHI, payment data, customer identifier) originates, where it flows through the system, where it lands, and who reads it. Triggers on: data lineage, PII flow, GDPR right-to-erasure, data residency, schema change, retention policy, 'where does field X come from', 'who reads PII column Y', 'walk this field end to end', deletion propagation, audit trail, data classification."
short-description: "PII / PHI / regulated-field lineage end to end"
```

- **Primary tools.** `codegraph_explore` (one-shot field walk),
  `codegraph_search` (locate the field / type / serializer),
  `codegraph_callers` / `codegraph_callees` (readers / writers),
  `codegraph_impact` (schema-change radius). Falls back to
  `rg` + file read for each.
- **Superpowers boundary.** `systematic-debugging` is
  symptom-shaped and one-machine. `data-lineage-trace` is
  **field-shaped, forward-only, regulated** and emits an
  audit-trail artefact.
- **Scenarios.**
  1. "We added a new column `national_id` — walk every read
     site, every write site, every log, every export job, and
     every downstream consumer." Produces a lineage report
     grouped by: producer / transformer / consumer / external
     sink.
  2. "GDPR right-to-erasure request for user 12345 — list every
     persistence site and every external processor we must
     notify." Produces an erasure checklist grouped by
     retention class.
  3. "Schema migration on `orders.total` — what breaks in the
     analytics warehouse, the BI dashboard, the export cron,
     and the regulator-facing report?" Emits impact + a
     data-residency note if any sink is in a different region.

## 2. `code-ownership-impact` — Ownership + Impact archetype

```yaml
description: "Use before any change that crosses team or service boundaries, to find CODEOWNERS, blast radius, and required reviewers. Triggers on: rename, refactor, deprecate, schema change, API break, cross-team change, 'who owns this', 'who needs to review', 'what does this PR break', CODEOWNERS, service boundary, dependency request, shared kernel."
short-description: "Cross-team blast radius + CODEOWNERS for shared-kernel changes"
```

- **Primary tools.** `codegraph_impact` (radius), `codegraph_files`
  (path → team map), `codegraph_callers` / `codegraph_callees`
  (boundary crossing), `codegraph_search` (locate shared
  kernels). Falls back to `rg CODEOWNERS` + `git log` per
  affected path.
- **Superpowers boundary.** `requesting-code-review` is
  one-team, one-PR, one-reviewer — no CODEOWNERS, no service
  boundary, no shared-kernel concept. `code-ownership-impact`
  produces a **multi-team review packet** with named approvers
  per affected path.
- **Scenarios.**
  1. "I'm renaming `OrderService.place` — list every service
     that calls it, the team that owns each caller, and the
     CODEOWNERS entry per path." Emits: per-team ping list,
     per-path CODEOWNERS, suggested review order.
  2. "I'm extracting `common-utils` into a shared kernel —
     which services need to be told before the package is
     published?" Emits: subscriber list + suggested
     release-train slot.
  3. "I'm changing the wire format of `UserCreatedEvent` — who
     is downstream, and which consumer team owns the migration
     window?" Emits: cross-team RFC recipients + deprecation
     window.

## 3. `incident-call-trace` — Trace archetype

```yaml
description: "Use during a live incident or while writing a postmortem, to walk the call path from a user-facing entry point (HTTP route, queue consumer, cron, CLI) to the persistence layer or external system involved in the failure. Triggers on: incident, outage, postmortem, on-call, SLO breach, 'walk the call path from the route to the DB', 'trace the failing request', 'which service did the timeout come from', runbook, customer impact."
short-description: "End-to-end call tracing for incident response and postmortems"
```

- **Primary tools.** `codegraph_explore` (full trace),
  `codegraph_callers` / `codegraph_callees` (each hop),
  `codegraph_search` (locate the failing symbol),
  `codegraph_node` (full source of the failing function).
  Falls back to log search + source reading.
- **Superpowers boundary.** `systematic-debugging` is local.
  `incident-call-trace` is **distributed-system** — given a
  customer-facing symptom and a stack of services, walk the
  path that led to the breach and emit a customer-impact +
  on-call-handoff artefact.
- **Scenarios.**
  1. "Checkout 5xx spike started at 14:02 — walk from the
     public route through auth, cart, payment-gateway, and
     persistence; identify the first hop that degraded."
  2. "Postmortem for incident #INC-9021 — produce the call-path
     diagram from `POST /api/v2/refund` to the third-party
     processor, including every retry and timeout boundary."
  3. "SLO breach on `search.p99 > 800ms` — which downstream
     call is the new tail-latency contributor?" Emits a
     delta vs. last green build.

## 4. `progressive-rollout-checklist` — Impact + Change-governance archetype

```yaml
description: "Use before merging a change that ships behind a feature flag, kill switch, canary, dark launch, or progressive rollout. Triggers on: feature flag, kill switch, canary, dark launch, progressive delivery, rollout plan, 'how do I ship this safely', 'what's the rollback plan', 'is the flag wired up', 'who is in the canary cohort', sunset, deprecation, rollback."
short-description: "Rollout-readiness, flag-coverage, and rollback-path verification"
```

- **Primary tools.** `codegraph_search` (locate the flag),
  `codegraph_callers` / `codegraph_callees` (every read / write
  of the flag), `codegraph_impact` (radius if the default
  flips), `codegraph_files` (route-to-cohort mapping). Falls
  back to `rg` over flag names + manual config review.
- **Superpowers boundary.** `finishing-a-development-branch` is
  merge / PR — it does not know flags, cohorts, canary,
  rollback, or sunset. This skill emits a **rollout plan +
  flag-coverage report + rollback path** and refuses to mark
  a flag-bearing change as "ready to ship" until the report
  is signed off.
- **Scenarios.**
  1. "I'm flipping the default of `use_new_pricing_engine` from
     off to on — list every caller, every default value, every
     test that hard-codes the old behaviour, and the rollback
     hook." Emits a rollout packet.
  2. "I'm adding a new payment provider behind a flag — which
     canary cohort and which kill switch is wired up? What
     metric gates the canary expansion?" Emits a canary plan
     with named SLO gate.
  3. "I'm sunsetting `v1/orders` — produce the deprecation
     window, the per-consumer warning, and the date the route
     will be removed." Emits a sunset calendar.

## 5. `compliance-control-walk` — Impact + Compliance archetype

```yaml
description: "Use when you need to demonstrate that a regulatory control (access, audit logging, encryption, retention, consent) actually exists in the running system — not just on a wiki. Triggers on: SOX, ISO 27001, SOC 2, GDPR, HIPAA, PCI, audit, 'show me the audit log for X', 'is PII encrypted at rest', 'who can read this field', 'is consent checked here', control walkthrough, evidence collection."
short-description: "Evidence collection for SOX / ISO / GDPR / HIPAA / PCI controls"
```

- **Primary tools.** `codegraph_explore` (one-shot control
  walk), `codegraph_search` (locate the control symbol —
  interceptor, guard, middleware), `codegraph_callers` /
  `codegraph_callees` (every enforcement site, every bypass),
  `codegraph_files` (path-to-control-owner map). Falls back to
  per-feature source inspection.
- **Superpowers boundary.** `verification-before-completion`
  runs *tests*; compliance is a *production* claim about
  *running* behaviour, and a control can pass tests but still
  be bypassed at runtime by a forgotten `if (env != 'prod')`
  path. This skill walks the graph of where the control is
  enforced and where it is bypassed.
- **Scenarios.**
  1. "Auditor wants evidence that every write to `Patient` is
     audit-logged. Find every code path that mutates a
     `Patient` and confirm the audit log is invoked on each —
     flag any path that doesn't."
  2. "GDPR right-to-erasure control — find every persistence
     sink and confirm the erasure job touches it. Flag any
     sink that is reachable but not covered by the erasure
     worker."
  3. "PCI scope reduction — list every code path that handles
     a primary account number and confirm tokenisation happens
     at the ingress."

## 6. `cross-team-rfc-draft` — Ownership + Multi-team archetype

```yaml
description: "Use when a planned change crosses team, service, or product boundaries and needs an RFC, design-doc review, or dependency request before implementation. Triggers on: RFC, design doc, cross-team review, dependency request, shared-kernel change, deprecation proposal, vendor swap, 'I need sign-off from team X', 'this changes a public interface', 'I want to deprecate endpoint Y'."
short-description: "Cross-team RFC drafting with named approvers and migration windows"
```

- **Primary tools.** `codegraph_impact` (who is downstream),
  `codegraph_files` (path → owner), `codegraph_callers`
  (subscribers), `codegraph_search` (locate the interface /
  contract). Falls back to manual CODEOWNERS + consumer
  enumeration.
- **Superpowers boundary.** `brainstorming` produces a spec
  for one team; `writing-plans` decomposes for one team.
  Neither produces a document that must be **reviewed and
  signed off by N other teams**, each with their own calendar
  and release train.
- **Scenarios.**
  1. "I want to change the auth token format — draft an RFC,
     list the N teams that consume it, suggest the review
     order, and pre-fill the migration-window section for each
     consumer team."
  2. "I'm proposing to extract a shared `event-bus` package —
     draft the RFC and pre-fill the dependency-request section
     per consumer service."
  3. "I'm proposing to deprecate the public `GET /v1/users`
     endpoint — draft the RFC with the deprecation window, the
     per-consumer notification list, and the suggested sunset
     date."

## 7. `observability-coverage-audit` — Coverage + Reliability archetype

```yaml
description: "Use to find which code paths, error states, or failure modes are exercised by tests, monitors, alerts, or runbooks — and which are silent. Triggers on: blind spot, missing alert, 'do we have a test for this error path', 'do we alert on this', 'what happens if X fails', 'where is the runbook for Y', coverage audit, dead-letter, timeout path, partial failure, gap analysis."
short-description: "Gap audit: which error paths lack tests, alerts, or runbooks"
```

- **Primary tools.** `codegraph_explore` (one-shot fanout),
  `codegraph_search` (locate the error / exception / throw),
  `codegraph_callers` (every site that can hit it),
  `codegraph_files` (test file filter). Falls back to
  per-feature file inspection.
- **Superpowers boundary.** `test-driven-development` writes
  a test for a known behaviour. `systematic-debugging` finds
  a root cause for a known symptom. Both assume a *known*
  failure shape. `observability-coverage-audit` is a
  **gap-detector** — given a code path or a failure class,
  list every site and then ask "for this site, is there a
  test, an alert, a runbook, a metric? If not, why is it OK?".
- **Scenarios.**
  1. "Find every `throw` in the payments path and report: is
     there a unit test, an integration test, an alert, and a
     runbook for each?" Emits a coverage matrix with gap rows.
  2. "Find every `catch (e)` that swallows the error — is each
     one logged? Is each one alerted on?" Emits a
     silent-failure list.
  3. "Find every external HTTP call — is each one wrapped in a
     circuit breaker, a timeout, and a retry? Which is
     missing?" Emits a resilience gap list, ready to feed
     `progressive-rollout-checklist`.

## 8. `pr-authoring` — Author-side PR preparation

```yaml
description: "Use when preparing a pull request for an enterprise repo: organise the PR description, attach the verification evidence, name the required reviewers per CODEOWNERS, surface linked issues / RFCs / design docs, and emit a 中文 + 企业 CI 兼容的 PR body. Triggers on: 写 PR, PR 描述, 整理 commit, 准备合入, code review 请求, 列出 reviewer, 附 verification 证据, link issue, link RFC, 关联设计文档."
short-description: "PR 准备: 描述、reviewer 名单、verification 证据、跨仓库引用"
```

- **Primary tools.** `codegraph_impact` (radius of the diff),
  `codegraph_files` (path → CODEOWNERS), `git log` for the
  changed paths, the issue tracker / RFC doc store. Falls
  back to a single-team `git diff --stat` + `gh pr view`
  style listing.
- **Superpowers boundary.** `requesting-code-review` is the
  English / generic version. `pr-authoring` is the **中文 +
  企业 CI 实践 + 多 reviewer 列表** variant, emits a PR body
  shaped for the consumer's existing PR template, and links
  RFCs / design docs / issues in the format the org uses.
- **Scenarios.**
  1. "This PR changes 4 files across 2 services — generate the
     PR body in the org's template, list the per-path
     CODEOWNERS, attach the verification matrix (unit,
     integration, e2e, manual), and link the RFC that
     approved it."
  2. "I have 6 commits queued locally — squash / reword /
     split them into a review-friendly PR description with
     per-commit intent."
  3. "This PR is a follow-up to a design doc and an RFC — link
     both in the PR body, and pre-fill the cross-team ping
     list for the reviewers who weren't on the original RFC."

## Cross-archetype coverage

| Archetype | Skills |
|---|---|
| Lineage | `data-lineage-trace` |
| Trace | `incident-call-trace` |
| Impact | `code-ownership-impact`, `progressive-rollout-checklist`, `compliance-control-walk` |
| Ownership | `code-ownership-impact`, `cross-team-rfc-draft`, `pr-authoring` |
| Coverage | `observability-coverage-audit` |
| Author-side PR | `pr-authoring` |

## Release plan

- **v0.5.0** (this release, 0 new skills): the thin
  orchestrator rework of `goal-driven-development` +
  `codegraph.md` prompt fix + this RFC + a knowledge note
  recording the decision.
- **v0.5.1 / v0.6.0** (next release): the 8 skills above,
  landed in 2-3 commits per skill with the same shape as
  `spec-doc-linter` (SKILL.md + README + co-required prompts
  + bad/good example pairs + canonical output sample).
  Likely order by leverage: `data-lineage-trace`,
  `code-ownership-impact`, `incident-call-trace`,
  `progressive-rollout-checklist`, `compliance-control-walk`,
  `observability-coverage-audit`, `cross-team-rfc-draft`,
  `pr-authoring`.

## Rejected alternatives

- **Ship all 8 in v0.5.0** (one mega-release). Rejected —
  each skill has its own enterprise context (regulator frame,
  flag system, on-call tool) and the 8 together would make
  v0.5.0 too coarse to roll back if any one is wrong.
- **Add a single `enterprise-dev-flow` mega-skill that
  embeds all 8 archetypes.** Rejected — would violate
  CONTRIBUTING.md "one scenario per prompt" rule and would
  bloat the SKILL.md body well past 80 lines.
- **Defer all 8 and ship a thin orchestrator only.** That is
  what v0.5.0 actually does. The RFC is the parking lot for
  the 8 so the orchestrator can name them as future phases.

## Out of scope

- Per-skill versioning metadata, eval infrastructure, CI.
- Any new skill beyond the 8 listed.
- Changes to the existing 5 `*-code-reviewer` skills.
- Changes to the existing `spec-doc-linter` or
  `project-knowledge-capture` skills.
