# Compliance Control Walk Skill

`compliance-control-walk` walks a regulatory control across the running system and produces an auditor-ready evidence pack. It is the regulator-frame counterpart of `spec-doc-linter`: where the linter keeps `DevAgent.md` / `CONTEXT.md` honest, this skill keeps the controls honest in the running system.

## What It Walks

- A regulator frame (SOX, ISO 27001, SOC 2, GDPR, HIPAA, PCI).
- A specific control name within the frame.
- The canonical enforcement symbol.
- Every enforcement site.
- Every bypass (every site that touches the protected data but does not go through the canonical enforcement).
- A data-lineage cross-check (which sinks the control covers).

## Discovery Tools

- **`codegraph_explore`** — one-shot control walk.
- **`codegraph_search`** — locate the control symbol.
- **`codegraph_callers` / `codegraph_callees`** — every enforcement site, every bypass.
- **`codegraph_files`** — path-to-control-owner map.
- **Fallback** — per-feature source inspection.

## Output Contract

A per-control evidence pack: regulator frame + control ID, canonical enforcement symbol, per-site evidence, per-bypass finding, data-lineage cross-check, evidence-collector.

## Related Skills

- **`spec-doc-linter`** — keeps `DevAgent.md` / `CONTEXT.md` honest; this skill is the runtime counterpart.
- **`data-lineage-trace`** — the upstream input.
- **`goal-driven-development`** — orchestrates the spec-to-code flow; this skill is the regulator-frame counterpart.

## Files

```text
.
├── SKILL.md
├── README.md
├── prompts/
│   └── control-walk.md       # the regulator-frame + control-name framework
└── examples/
    ├── bad-control.md
    ├── good-control.md
    └── control-output.md     # canonical "what the agent should emit" sample
```
