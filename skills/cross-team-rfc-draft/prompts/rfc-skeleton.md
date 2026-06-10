# RFC Skeleton

> See also: ../SKILL.md

The 9-section RFC template. Every section has a fixed shape; an RFC missing any section is itself a `RFC-INCOMPLETE` finding.

## The 9 sections

1. **Title** — the proposed change in one sentence. Examples: "RFC-001: Migrate auth token format from JWT-v1 to JWT-v2".
2. **Status** — `DRAFT`, `IN-REVIEW`, `APPROVED`, `REJECTED`, `WITHDRAWN`. The default for a fresh draft is `DRAFT`.
3. **Author + team + date** — the author, the author's team, the date the RFC is filed.
4. **Summary** — 1-3 paragraphs. What changes, why now, what the alternatives are.
5. **Motivation** — why the change is needed. The problem the change solves; the data or incidents that justify it.
6. **Detailed design** — the technical design. Interfaces, schemas, code shapes, migration steps. Long enough to be reviewable; not a full implementation.
7. **Consumer-team impact** — one row per consumer team. For each: the named approver, the migration window, the per-consumer cost (small / medium / large), the release-train slot.
8. **Rollout plan** — flag flips, cohort expansions, sunset windows, kill switches, rollback hooks. Mirrors the `progressive-rollout-checklist` output for this change.
9. **Open questions + decision log** — every open question; for closed questions, the decision and the date.

## Per-consumer-team migration-window section

Section 7 (Consumer-team impact) is per-team. For each consumer, the section names:

- **Team** — the consumer's `@org/...` group.
- **Named approver** — the GitHub handle from `CODEOWNERS`.
- **Migration window** — the time between "we tell them" and "we cut them off". Examples: "8 weeks (consumer team has a quarterly release)".
- **Per-consumer cost** — `small` (under 1 dev-week), `medium` (1-4 dev-weeks), `large` (4+ dev-weeks).
- **Release-train slot** — the next train the consumer team uses, with the date.
- **Suggested review order** — the slot in the per-team review order this consumer should ack in.

A consumer missing the migration window or the named approver is itself a `RFC-CONSUMER-INCOMPLETE` finding. The RFC cannot be marked `IN-REVIEW` until every consumer has a row.

## The 4 review-order rules

For each consumer, the suggested review order is:

1. **Highest-dependency first.** The consumer whose service is the deepest in the dependency chain reviews first; the consumer whose service is the shallowest reviews last.
2. **Hotfix window first.** A consumer in a hot-fix window reviews first regardless of dependency order.
3. **Release-train-locked last.** A consumer whose next release is 3+ months out reviews last; the RFC author waits for the next train slot.
4. **Consumer-cost-weighted.** A consumer with `large` per-consumer cost reviews earlier than a consumer with `small` per-consumer cost (because the large one needs more time).

## The "all clean" sentinel

The skill is not a linter. An RFC where every section is filled and every consumer has a row is the "no finding" case. There is no canonical "all clean" line; the RFC itself is the output.
