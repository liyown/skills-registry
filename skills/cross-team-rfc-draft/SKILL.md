---
name: cross-team-rfc-draft
description: "Use when a planned change crosses team, service, or product boundaries and needs an RFC, design-doc review, or dependency request before implementation. Triggers on: RFC, design doc, cross-team review, dependency request, shared-kernel change, deprecation proposal, vendor swap, 'I need sign-off from team X', 'this changes a public interface', 'I want to deprecate endpoint Y'. Do NOT use for: design-time brainstorming for a single team (use obra/superpowers brainstorming, then obra/superpowers writing-plans), or for a per-PR review packet (use code-ownership-impact). Complementary to code-ownership-impact, which is the per-PR packet; cross-team-rfc-draft is the long-form RFC that has to be reviewed and signed off by N other teams, each with their own calendar and release train."
metadata:
  short-description: Cross-team RFC drafting with named approvers and migration windows
---

# Cross-Team RFC Draft

Draft an RFC for a change that crosses team, service, or product boundaries. The output is a review-ready RFC skeleton with named approvers, per-consumer-team migration windows, and a release-train-aware review order.

## Required Loading

Always load:

- `prompts/rfc-skeleton.md` — the 9-section RFC template; the per-consumer-team migration-window section.

## When To Run

- Before any change that crosses a public interface (auth token format, event schema, public API contract).
- Before any deprecation proposal (an endpoint, a package, a feature flag).
- Before any shared-kernel extraction (a package, a service, a CI template).
- Before any vendor swap (Stripe → Adyen, Datadog → Honeycomb, etc.).
- When a team has been blocked by another team's calendar and needs a written proposal to push the conversation past chat threads.

## Discovery Order

1. Identify the proposed change (one sentence: what changes).
2. Identify the author's team.
3. Identify every consumer team (every team whose service / product consumes the changed interface, the changed package, or the changed endpoint).
4. For each consumer team, find the next release-train slot.
5. Find the consumer-side deprecation window (the time between "we tell them" and "we cut them off").
6. Find the per-consumer-team migration cost (small / medium / large) — usually from a past similar migration.
7. Draft the 9 sections in order, with the per-consumer-team migration window section filled for each consumer.

## Output Contract

The RFC is a single Markdown document with 9 sections in a fixed order. Each section has a fixed shape. The consumer-team section is per-team: for each consumer, the migration window, the named approver, the suggested review order.

## Tools

- **`codegraph_impact`** — every downstream consumer of the changed interface.
- **`codegraph_files`** — path → owner (for the per-consumer-team attribution).
- **`codegraph_callers`** — every site that calls into the interface.
- **`codegraph_search`** — locate the interface / contract.
- **Fallback** — manual CODEOWNERS + consumer enumeration.

## Fallback

If CodeGraph is unavailable, the fallback is manual CODEOWNERS + consumer enumeration. The final report must include the line:

```text
CodeGraph unavailable; rfc gathered by rg/file inspection.
```

## Related Skills

- **`spec-doc-linter`** — keeps `DevAgent.md` / `CONTEXT.md` honest.
- **`code-ownership-impact`** — the per-PR packet; this skill is the long-form RFC.
- **`goal-driven-development`** — orchestrates the spec-to-code flow; this skill is the cross-team counterpart.

## Examples

Each `bad-*` doc has a matching `good-*` in `examples/` showing the minimum RFC. Read them side by side to calibrate depth. `examples/rfc-output.md` is the canonical "what the agent should emit" sample.
