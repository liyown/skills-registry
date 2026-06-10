# Incident Call Trace Skill

`incident-call-trace` walks the call path from a user-facing entry point to the persistence layer or external system during a live incident or while writing a postmortem. It is the distributed-system counterpart of `spec-doc-linter`: where the linter keeps `DevAgent.md` / `CONTEXT.md` honest, this skill keeps the call-graph honest for an incident.

## What It Traces

- A user-facing entry point (HTTP route, queue consumer, cron, CLI).
- Every hop from the entry point to persistence or an external system.
- Per-hop latency budget vs. observed.
- Per-team paging list (the team that owns each hop).
- Customer-impact framing (how many customers / requests / regions are affected).

## Discovery Tools

- **`codegraph_explore`** — primary one-shot walk.
- **`codegraph_callers` / `codegraph_callees`** — per-hop walk.
- **`codegraph_search`** — locate the failing symbol.
- **`codegraph_node`** — full source of the failing function.
- **Fallback** — log search + source reading.

## Output Contract

A per-hop trace grouped by call site. Above the trace: the **primary incident locus** (first hop that degraded). Below: the per-team paging list + customer-impact framing.

## Related Skills

- **`spec-doc-linter`** — keeps `DevAgent.md` / `CONTEXT.md` honest; this skill is the runtime counterpart.
- **`observability-coverage-audit`** (planned v0.6.1) — the post-incident coverage audit.
- **`goal-driven-development`** — orchestrates the spec-to-code flow; this skill is the on-call counterpart.

## Files

```text
.
├── SKILL.md
├── README.md
├── prompts/
│   └── call-trace.md       # the per-hop trace framework
└── examples/
    ├── bad-trace.md
    ├── good-trace.md
    └── trace-output.md     # canonical "what the agent should emit" sample
```
