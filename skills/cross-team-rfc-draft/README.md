# Cross-Team RFC Draft Skill

`cross-team-rfc-draft` drafts an RFC for a change that crosses team, service, or product boundaries. The output is a review-ready Markdown document with 9 fixed sections.

## What It Drafts

- A 9-section RFC skeleton.
- Per-consumer-team migration windows.
- A release-train-aware review order.
- Named approvers per consumer team.

## Discovery Tools

- **`codegraph_impact`** — every downstream consumer.
- **`codegraph_files`** — path → owner.
- **`codegraph_callers`** — every site that calls into the interface.
- **`codegraph_search`** — locate the interface / contract.
- **Fallback** — manual CODEOWNERS + consumer enumeration.

## Output Contract

A single Markdown document with 9 sections in a fixed order. The consumer-team section is per-team.

## Related Skills

- **`spec-doc-linter`** — keeps `DevAgent.md` / `CONTEXT.md` honest.
- **`code-ownership-impact`** — the per-PR packet; this skill is the long-form RFC.
- **`goal-driven-development`** — orchestrates the spec-to-code flow; this skill is the cross-team counterpart.

## Files

```text
.
├── SKILL.md
├── README.md
├── prompts/
│   └── rfc-skeleton.md       # the 9-section RFC template
└── examples/
    ├── bad-rfc.md
    ├── good-rfc.md
    └── rfc-output.md     # canonical "what the agent should emit" sample
```
