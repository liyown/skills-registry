---
name: data-lineage-trace
description: "Use when you need to prove where a piece of data (PII, PHI, payment data, customer identifier) originates, where it flows through the system, where it lands, and who reads it. Triggers on: data lineage, PII flow, GDPR right-to-erasure, data residency, schema change, retention policy, 'where does field X come from', 'who reads PII column Y', 'walk this field end to end', deletion propagation, audit trail, data classification. Do NOT use for: static code review (use java-code-reviewer, react-code-reviewer, go-code-reviewer, python-code-reviewer, or node-code-reviewer for that), or for design-time brainstorming (use obra/superpowers brainstorming). Complementary to goal-driven-development, which orchestrates the spec-to-code flow; data-lineage-trace is the regulated-data counterpart that runs after a change lands."
metadata:
  short-description: PII / PHI / regulated-field lineage end to end
---

# Data Lineage Trace

Trace a regulated data field from producer to sink. Forward walk for ingestion, backward walk for deletion, full graph for classification. The output is an audit-trail artefact, not a review finding.

## Required Loading

Always load:

- `prompts/lineage-walk.md` — the producer/transformer/consumer/retention framework; the four lineage report sections.

## When To Run

- After a schema change that adds or renames a regulated field.
- After a new persistence sink is added (warehouse, analytics, export job).
- When a GDPR / HIPAA / PCI / SOC 2 right-to-erasure request arrives.
- When a regulator asks "show me every site that touches `<field>`".
- When a feature team is about to consume a field and wants to know what it actually is (column vs computed vs joined).

## Discovery Order

1. Locate the field by name (`rg` for the literal, with extension `.ts`/`.java`/`.go`/`.py`/`.sql`/`.proto`).
2. Identify the type / schema definition (`codegraph_search` for the type, or `rg` for the struct / class / table).
3. Trace producers: every site that writes the field. Use `codegraph_callers` (filtered to write-side symbols) or `rg '<field>\s*[:=]'`.
4. Trace transformers: every site that mutates the field between producer and sink.
5. Trace consumers: every reader, every export, every log, every dashboard, every external processor.
6. Classify each sink by retention class (transactional / analytics / log / external) and by residency (region / cross-border).

## Output Contract

The lineage report is grouped by sink, not by file. For each sink, name:

- the producer (which path writes the field there)
- the transformer (if any) and the alias / cast along the way
- the consumer (who reads the field from this sink)
- the retention class
- the residency note (if cross-border)

A field that has no producer in the walked graph is itself a finding (mode `LIN-ORPHAN-READ`): someone is reading a field that nothing in the current codebase writes.

If no drift is found in any scanned doc, output exactly the sentinel:
`未发现文档与代码漂移。` — but **this is the spec-doc-linter sentinel, not a data-lineage sentinel**. The data-lineage skill has its own "no finding" line:

```text
未发现数据血缘漂移。
```

When a finding exists, output:

```text
# Lineage Report

## <field-name>

- producer: <path>:<line>
- transformer: <path>:<line> (alias / cast if any)
- consumer: <path>:<line> OR <external-sink>:<id>
- retention: <transactional | analytics | log | external>
- residency: <region-or-cross-border>
- control: <audit-logged? | encrypted-at-rest? | consent-checked?>
```

## Tools

- `codegraph_explore` for the one-shot field walk when CodeGraph is available.
- `rg` for the field name + per-path `Read` when CodeGraph is not available.
- Type / schema / manifest introspection is per-language: Java `javap`, TypeScript `tsc --noEmit`, Go `go doc`, Python `ast` or `pyright`, SQL via the database migration tool.

## Fallback

If CodeGraph is unavailable, the fallback is `rg` + source reading. The final report must include the line:

```text
CodeGraph unavailable; lineage gathered by rg/file inspection.
```

## Related Skills

- **`spec-doc-linter`** — keeps `DevAgent.md` / `CONTEXT.md` honest. data-lineage-trace is the runtime-data counterpart.
- **`compliance-control-walk`** (planned v0.6.1) — the regulator-frame audit that uses this lineage report as one of its inputs.
- **`goal-driven-development`** — orchestrates the spec-to-code flow; data-lineage-trace is the post-impl counterpart that verifies the new code's field flows match the spec.

## Examples

Each `bad-*` doc has a matching `good-*` in `examples/` showing the minimum sync. Read them side by side to calibrate lineage depth. `examples/lineage-output.md` is the canonical "what the agent should emit" sample.
