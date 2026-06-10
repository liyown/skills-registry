# Lineage Walk Framework

> See also: ../SKILL.md

A field's lineage is grouped by **sink**, not by file. For each sink, the framework below says what to record.

## The 6 cells per sink

For every sink that the field reaches, name:

1. **Producer** — the path + line that writes the field to this sink. If multiple producers, name each.
2. **Transformer** — the path + line of any cast, alias, derive, or join that mutates the field between the producer and the sink. If no transformer, name the field unchanged.
3. **Consumer** — the path + line of every reader; or, for an external sink, the system id (e.g. `snowflake://analytics.orders_fact`).
4. **Retention class** — one of: `transactional` (live in the operational DB), `analytics` (warehouse / BI), `log` (log aggregator), `external` (third-party processor, payment gateway, etc.).
5. **Residency** — the region the sink lives in. If cross-border, name both sides.
6. **Control coverage** — which of these are present, named by the file that implements them: `audit-logged`, `encrypted-at-rest`, `consent-checked`. Any sink missing a control it should have is itself a `LIN-CONTROL-GAP` finding.

## The 4 lineage report sections

The final report is structured as four sections, in this order:

1. **Producers** — every site that writes the field. Grouped by language / module.
2. **Transformers** — every site that mutates the field between producer and sink. Grouped by what the transformation does.
3. **Sinks** — every destination the field reaches. For each sink, the 6 cells above.
4. **Controls** — which controls are wired up across the lineage, and which sinks lack a required control. This section is the input to `compliance-control-walk` (planned v0.6.1).

## The orphan-read finding

A field that has no producer in the walked graph is a `LIN-ORPHAN-READ` finding. Some code reads a field that nothing in the current codebase writes — typical of a renamed / dropped column still being read by a dashboard. The finding lists every consumer that reads the orphan.

## The deletion-propagation finding (for GDPR / right-to-erasure)

When a deletion request arrives, the report is the same graph walked **backward**: for the row identifier, name every site that holds a copy. Each site is a `LIN-DELETION-TARGET`. Any site that is not in the erasure worker's coverage list is a `LIN-DELETION-GAP` finding.

## Schema-change radius

When the field's type or name changes, every site in the 4 sections is a `LIN-SCHEMA-BREAK` candidate. Group by impact:
- **Type-incompatible producers** — the producer's write type no longer matches.
- **Type-incompatible consumers** — the consumer's read type no longer matches.
- **Schema-version consumers** — the consumer is on a stale schema version (warehouse, BI, export job).

## The "all clean" sentinel

```text
未发现数据血缘漂移。
```

This is **deliberately distinct** from:
- the reviewer sentinel `未发现明确高风险问题。`
- the spec-doc-linter sentinel `未发现文档与代码漂移。`

A consumer grepping for the data-lineage sentinel knows the graph walk found nothing; a consumer grepping for the spec-doc-linter sentinel knows the doc-code walk found nothing.
