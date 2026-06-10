# Data Lineage Trace Skill

`data-lineage-trace` walks a regulated data field from producer to sink and emits a per-sink lineage report. It is the audit-trail counterpart of `spec-doc-linter`: where the linter keeps `DevAgent.md` / `CONTEXT.md` honest, this skill keeps the runtime data flows honest.

## What It Traces

- A named field (column, attribute, JSON key) on a class / struct / table / schema.
- The producer (which site writes it).
- The transformer chain (cast, alias, join, derivation).
- The consumer (reader, export, log, dashboard, external processor).
- The retention class (transactional, analytics, log, external).
- The residency (region, cross-border, regulator-bound).
- The control coverage (audit-logged, encrypted, consent-checked).

## Discovery Tools

- **`codegraph_explore`** — one-shot walk. Asks "how does field X flow" end-to-end and returns a relationship map. The primary tool.
- **`codegraph_search`** — locate the type / column / serializer by name.
- **`codegraph_callers` / `codegraph_callees`** — narrow the walk to write-side or read-side.
- **`codegraph_impact`** — schema-change radius.
- **Fallback** — `rg '<field>\s*[:=]'` and per-path `Read`. Falls back when CodeGraph is unavailable; final report declares the fallback.

## Output Contract

A lineage report grouped by sink, not by file. Each sink entry names producer / transformer / consumer / retention / residency / control. An orphan read (a consumer with no in-graph producer) is itself a `LIN-ORPHAN-READ` finding.

Sentinel (no drift): `未发现数据血缘漂移。`

## Related Skills

- **`spec-doc-linter`** — keeps `DevAgent.md` / `CONTEXT.md` honest. data-lineage-trace is the runtime-data counterpart.
- **`compliance-control-walk`** (planned v0.6.1) — the regulator-frame audit that uses this lineage report.
- **`goal-driven-development`** — orchestrates the spec-to-code flow; data-lineage-trace is the post-impl counterpart.

## Files

```text
.
├── SKILL.md
├── README.md
├── prompts/
│   └── lineage-walk.md       # the producer/transformer/consumer/retention framework
└── examples/
    ├── bad-lineage.md
    ├── good-lineage.md
    └── lineage-output.md     # canonical "what the agent should emit" sample
```
