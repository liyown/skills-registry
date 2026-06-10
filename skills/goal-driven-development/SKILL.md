---
name: goal-driven-development
description: "Use this skill to orchestrate an existing goal or spec through the dev flow. Load when the user asks for end-to-end feature execution, a goal-to-PR pipeline, or any task that crosses spec → implementation → review → knowledge capture. The skill is a thin router: it names the 6 phases (spec intake, blast-radius check, implementation, verification, review gate, knowledge capture) and the skill that does the work for each phase. Do NOT use for: a single review (load the matching `*-code-reviewer` directly), or for a generic author-craft task like TDD / debugging / brainstorming (use obra/superpowers). Complementary to the 17 enterprise-flow skills; for the planned v0.6.1+ skills, the orchestrator names which of them does the work for each phase."
metadata:
  short-description: Spec → code orchestration with per-phase skill routing
---

# Goal Driven Development

Use this skill to orchestrate an existing goal or spec through the dev flow. The skill itself is a thin router — the actual work is done by the language-specific reviewers and the planned enterprise-flow skills in `docs/rfcs/2026-06-10-enterprise-flow-skills.md`.

## Required Loading

Always load:

- `prompts/workflow.md`
- `prompts/codegraph.md`

## Orchestration Contract

The 6 phases below are how a spec becomes code in this collection. Each phase names the skill that does the heavy lifting; this skill's job is to name them, in order, and stop the run if any phase finds an issue that should be raised before continuing.

1. **Spec intake** — read the spec, extract acceptance criteria. (The dedicated `spec-intake` skill is not yet shipped; until it is, this phase uses the prompt shape in `prompts/workflow.md` Phase 1.)
2. **Blast-radius check** — before implementing, find every caller, every consumer, every shared kernel that the change touches. Currently this is done inline; the dedicated `code-ownership-impact` skill is in v0.5.0+ RFC.
3. **Implementation** — minimal scoped change. Do not expand into unrelated refactors.
4. **Verification** — run the nearest tests, builds, or page checks. Record the inability to run if it applies.
5. **Review gate** — invoke the matching `*-code-reviewer` (java / react / go / python / node) for the touched files.
6. **Knowledge capture** — invoke `project-knowledge-capture` to write the durable post-impl note.

If any of `*-code-reviewer` surfaces Critical / High, fix or explicitly mark `需要结合上下文确认` before moving to phase 6. Never silently skip.

## CodeGraph Fallback

If CodeGraph is unavailable, fall back to `rg` and source reading, and **declare** the fallback in the final report using the exact line `CodeGraph unavailable; context was gathered by rg/file inspection.` — never silently downgrade. The current CodeGraph tool list is in `prompts/codegraph.md`; the prompt references `codegraph_explore` as the primary walk tool.

## Related Skills

- **Future enterprise-flow skills** (proposed in `docs/rfcs/2026-06-10-enterprise-flow-skills.md`): `data-lineage-trace`, `code-ownership-impact`, `incident-call-trace`, `progressive-rollout-checklist`, `compliance-control-walk`, `cross-team-rfc-draft`, `observability-coverage-audit`, `pr-authoring`. Each of those sits in a different phase and complements this orchestrator.
- **Already shipped:** `spec-doc-linter` (keeps `DevAgent.md` and `CONTEXT.md` honest after the code has moved); the 5 language-specific reviewers invoked in phase 5.

See `README.md` for install instructions and the combined install command.
