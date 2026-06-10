---
name: compliance-control-walk
description: "Use when you need to demonstrate that a regulatory control (access, audit logging, encryption, retention, consent) actually exists in the running system — not just on a wiki. Triggers on: SOX, ISO 27001, SOC 2, GDPR, HIPAA, PCI, audit, 'show me the audit log for X', 'is PII encrypted at rest', 'who can read this field', 'is consent checked here', control walkthrough, evidence collection. Do NOT use for: a single test-run verification (use obra/superpowers verification-before-completion), or for static code review (use the matching `*-code-reviewer` for that). Complementary to goal-driven-development, which orchestrates the spec-to-code flow; compliance-control-walk is the regulator-frame counterpart that runs against the production-shape behaviour, not the test suite."
metadata:
  short-description: Evidence collection for SOX / ISO / GDPR / HIPAA / PCI controls
---

# Compliance Control Walk

Walk a regulatory control across the running system and produce an evidence pack. The output is auditor-ready: per-control coverage matrix, per-site call-graph evidence, and a per-bypass finding for any reachable-but-uncontrolled path.

## Required Loading

Always load:

- `prompts/control-walk.md` — the regulator-frame + control-name framework; the per-control evidence pack.

## When To Run

- When an auditor asks for evidence on a specific control ("show me every site that audit-logs writes to `Patient`").
- When a SOC 2 / ISO 27001 / SOX / GDPR / HIPAA / PCI control needs an annual walkthrough.
- When a control was added to a wiki but never verified in the codebase.
- When a control was added to the codebase but the test does not cover a known bypass path.

## Discovery Order

1. Identify the regulator frame (SOX, ISO 27001, SOC 2, GDPR, HIPAA, PCI). Each frame has a different set of canonical controls.
2. Identify the specific control name within the frame. Examples: SOX `ITGC.04` (change management), ISO `A.9.4.1` (information access restriction), GDPR Article 17 (right to erasure), HIPAA `164.312(a)(1)` (access control), PCI `3.4` (PAN protection at rest).
3. For the control, identify the canonical enforcement symbol. Examples: an interceptor, a guard, a middleware, a deletion worker.
4. Find every enforcement site (`codegraph_search` for the symbol, or `rg`).
5. Find every bypass — every site that touches the same data but does not go through the canonical enforcement. A bypass is itself a `COMP-BYPASS` finding.
6. For each enforcement site, capture per-site evidence: a code excerpt + a path + a line number, suitable for the auditor to verify.
7. Cross-check the data-lineage report (if available) to see which sinks the control covers.

## Output Contract

The evidence pack is grouped by control. For each control, name:

- the regulator frame + control ID
- the canonical enforcement symbol
- the per-site evidence (path + line + code excerpt)
- the per-bypass finding (every site that touches the data but does not go through the enforcement)
- the data-lineage cross-check (which sinks the control covers; which sinks are uncovered)
- the evidence-collector (the path the auditor can run to re-verify)

A control with no enforcement site is a `COMP-MISSING` finding. A control with enforcement sites but a reachable bypass is a `COMP-BYPASS` finding.

## Tools

- **`codegraph_explore`** — one-shot control walk when the enforcement symbol is the entry point.
- **`codegraph_search`** — locate the control symbol by name.
- **`codegraph_callers` / `codegraph_callees`** — every enforcement site, every bypass.
- **`codegraph_files`** — path-to-control-owner map.
- **Fallback** — per-feature source inspection: every path that touches the protected data, manual review of the control coverage.

## Fallback

If CodeGraph is unavailable, the fallback is per-feature source inspection. The final report must include the line:

```text
CodeGraph unavailable; control walk gathered by rg/file inspection.
```

## Related Skills

- **`spec-doc-linter`** — keeps `DevAgent.md` / `CONTEXT.md` honest; this skill is the runtime-control counterpart.
- **`data-lineage-trace`** — the upstream input: a control without a data-lineage report is missing the sink list.
- **`goal-driven-development`** — orchestrates the spec-to-code flow; this skill is the regulator-frame counterpart.

## Examples

Each `bad-*` doc has a matching `good-*` in `examples/` showing the minimum evidence pack. Read them side by side to calibrate pack depth. `examples/control-output.md` is the canonical "what the agent should emit" sample.
