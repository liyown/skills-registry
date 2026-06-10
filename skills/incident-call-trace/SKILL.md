---
name: incident-call-trace
description: "Use during a live incident or while writing a postmortem, to walk the call path from a user-facing entry point (HTTP route, queue consumer, cron, CLI) to the persistence layer or external system involved in the failure. Triggers on: incident, outage, postmortem, on-call, SLO breach, 'walk the call path from the route to the DB', 'trace the failing request', 'which service did the timeout come from', runbook, customer impact. Do NOT use for: static code review (use java-code-reviewer, react-code-reviewer, go-code-reviewer, python-code-reviewer, or node-code-reviewer for that), or for local bug debugging (use obra/superpowers systematic-debugging). Complementary to goal-driven-development, which orchestrates the spec-to-code flow; incident-call-trace is the on-call counterpart that walks a live distributed-system failure."
metadata:
  short-description: End-to-end call tracing for incident response and postmortems
---

# Incident Call Trace

Walk the call path from a user-facing entry point to the persistence layer or external system involved in a live failure. The output is a per-hop trace + customer-impact framing + per-team paging list, suitable for an incident channel or a postmortem doc.

## Required Loading

Always load:

- `prompts/call-trace.md` — the per-hop trace framework; the customer-impact framing; the per-team paging list.

## When To Run

- During a live incident when the failing user-facing entry point is known but the failing hop is not.
- During postmortem authoring when the postmortem needs a call-path diagram.
- When an SLO breach is in progress and the team needs to know which downstream call is the new tail-latency contributor.
- When the on-call rotation is being paged and the runbook points to "trace the call path".

## Discovery Order

1. Identify the user-facing entry point (HTTP route, queue consumer name, cron job, CLI command).
2. Walk the call path forward: from the entry point through every middleware, service, and external call. Use `codegraph_explore` (one-shot) or `codegraph_callers` / `codegraph_callees` per hop.
3. At each hop, name the latency budget and the actual observed latency (if logs are available).
4. Identify the first hop that degraded against its budget. That hop is the **primary incident locus**.
5. List every team that owns a hop in the trace, with the on-call paging channel.
6. Compute the customer-impact framing: how many customers / requests / regions are affected.

## Output Contract

The trace is grouped by hop, not by file. For each hop, name:

- the path + line (the call site)
- the team that owns the hop
- the latency budget vs. observed (where available)
- the failure mode (timeout, error rate, saturation, etc.)

Above the trace, name the **primary incident locus** (the first hop that degraded). Below the trace, list the **per-team paging list** and the **customer-impact framing**.

## Tools

- **`codegraph_explore`** — primary one-shot walk from the entry point to the persistence / external system.
- **`codegraph_callers` / `codegraph_callees`** — per-hop walk when a finer-grained view is needed.
- **`codegraph_search`** — locate the failing symbol by name.
- **`codegraph_node`** — full source of the failing function.
- **Fallback** — log search (`rg '<entry-point>'` over log aggregators) + per-hop source reading. The fallback is acceptable but not as fast.

## Fallback

If CodeGraph is unavailable, the fallback is log search + source reading. The final report must include the line:

```text
CodeGraph unavailable; trace gathered by rg/file inspection.
```

## Related Skills

- **`spec-doc-linter`** — keeps `DevAgent.md` / `CONTEXT.md` honest; this skill is the runtime counterpart.
- **`observability-coverage-audit`** (planned v0.6.1) — the post-incident "could this have been detected?" audit; this skill produces the trace, that one produces the coverage matrix.
- **`goal-driven-development`** — orchestrates the spec-to-code flow; this skill is the on-call counterpart.

## Examples

Each `bad-*` doc has a matching `good-*` in `examples/` showing the minimum trace. Read them side by side to calibrate trace depth. `examples/trace-output.md` is the canonical "what the agent should emit" sample.
