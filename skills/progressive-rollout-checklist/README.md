# Progressive Rollout Checklist Skill

`progressive-rollout-checklist` verifies that a feature-flagged change is ready to ship. The output is a per-flag rollout packet. The skill refuses to mark a flag as "ready" if any required field is missing.

## What It Verifies

- The flag is wired up (every read site has a default value).
- The cohort is defined (the rule that determines who sees the new behaviour).
- The SLO gate exists (the metric + threshold that gates the cohort expansion).
- The kill switch is named (the action + on-call).
- The rollback hook is in place (the action + owner).
- The sunset date is set (if the flag is being deprecated).

## Discovery Tools

- **`codegraph_search`** — locate a flag by name.
- **`codegraph_callers` / `codegraph_callees`** — every read / write of the flag.
- **`codegraph_impact`** — radius if the default flips.
- **`codegraph_files`** — route-to-cohort mapping.
- **Fallback** — `rg '<flag-name>'` + per-path `Read`.

## Output Contract

A per-flag rollout packet. Each flag row has 8 cells: name, read sites, write sites, default + proposed default, cohort, SLO gate, kill switch, rollback hook, sunset.

## Related Skills

- **`spec-doc-linter`** — keeps `DevAgent.md` / `CONTEXT.md` honest.
- **`code-ownership-impact`** — names the per-team approvers; this skill adds the rollout-time checks.
- **`goal-driven-development`** — orchestrates the spec-to-code flow; this skill is the pre-ship counterpart.

## Files

```text
.
├── SKILL.md
├── README.md
├── prompts/
│   └── rollout-check.md       # the flag / cohort / SLO-gate / rollback framework
└── examples/
    ├── bad-rollout.md
    ├── good-rollout.md
    └── rollout-output.md     # canonical "what the agent should emit" sample
```
