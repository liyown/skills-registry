# Code Ownership Impact Skill

`code-ownership-impact` walks a diff and emits a per-team review packet with named approvers. It is the cross-team counterpart of `spec-doc-linter`: where the linter keeps `DevAgent.md` / `CONTEXT.md` honest, this skill keeps cross-team change accountable.

## What It Computes

- A path → team map for the diff.
- A named approver list per team, from `CODEOWNERS`.
- A shared-kernel membership list (any changed path under `common/`, `shared/`, `pkg/`, `platform/`, or a shared-kernel pattern).
- A blast radius for each shared-kernel change (subscriber services and their owning teams).
- A suggested review order that respects each team's release train.

## Discovery Tools

- **`codegraph_impact`** — primary tool for the blast-radius walk.
- **`codegraph_files`** — path → team map.
- **`codegraph_callers` / `codegraph_callees`** — boundary crossing.
- **`codegraph_search`** — locate a shared-kernel file by name.
- **Fallback** — `rg 'path:' CODEOWNERS` + `git log -- <path>` for per-path ownership inference.

## Output Contract

A review packet grouped by team. For each team, the packet names: paths owned, named approvers, suggested review deadline, cross-team impact, suggested ack-channel.

## Related Skills

- **`spec-doc-linter`** — keeps `DevAgent.md` / `CONTEXT.md` honest; this is the human-side counterpart.
- **`cross-team-rfc-draft`** (planned v0.6.1) — the long-form RFC; code-ownership-impact is the per-PR packet.
- **`goal-driven-development`** — orchestrates the spec-to-code flow; this skill is the pre-merge counterpart.

## Files

```text
.
├── SKILL.md
├── README.md
├── prompts/
│   └── ownership-walk.md       # the per-team-ping-list framework
└── examples/
    ├── bad-ownership.md
    ├── good-ownership.md
    └── ownership-output.md     # canonical "what the agent should emit" sample
```
