---
name: enterprise-stack
description: "Use when an agent needs the full enterprise-dev-flow skill set: language-specific code review (java/react/go/python/node), spec-to-code orchestration (goal-driven-development), durable post-impl knowledge (project-knowledge-capture), doc drift linting (spec-doc-linter), and (planned) data lineage / CODEOWNERS / incident trace / compliance / rollout / RFC / observability skills. Triggers on: 'set up enterprise skills', 'install the full stack', 'enterprise-stack', 'give me the augment', 'what skills should I load', 'superpower-enterprise'. Do NOT use for: a single review (load the matching `*-code-reviewer` directly); a generic TDD/debugging/brainstorming request (use obra/superpowers)."
metadata:
  short-description: Single install point for the superpower-enterprise collection
---

# Enterprise Stack

The single install point for the [`liyown/superpower-enterprise`](https://github.com/liyown/superpower-enterprise) collection. This meta-skill is a router: it names every skill the collection ships, the install command for each, and the closest `obra/superpowers` skill for the concerns the collection deliberately does not cover.

## Relationship to superpowers

This collection is the **enterprise augment layer** for [`obra/superpowers`](https://github.com/obra/superpowers). The two are complementary, not competing:

- `superpowers` covers generic author-craft: TDD (`test-driven-development`), debugging (`systematic-debugging`), brainstorming (`brainstorming`), writing plans (`writing-plans`), requesting/receiving code review (`requesting-code-review`, `receiving-code-review`), finishing a branch (`finishing-a-development-branch`), subagent-driven work, using git worktrees, using superpowers itself, and writing skills.
- This collection covers what `superpowers` deliberately refuses to assume: language-specific production-risk review, spec-to-code orchestration, durable post-impl knowledge, doc/spec drift, and the planned enterprise-only skills (data lineage, CODEOWNERS, incident trace, compliance, rollout, RFC, observability, pr-authoring).

Install both. A consumer who loads only `superpowers` is missing production-risk review depth and the org / multi-team / regulated layer; a consumer who loads only this collection is missing the generic author-craft.

## Required Loading

Always load `prompts/inventory.md` — the one-line-per-skill index that names the install command, the trigger, and the superpowers neighbour for every skill in the collection.

## Combined Install

```sh
# superpowers (the base layer)
npx skills add obra/superpowers

# superpower-enterprise (the augment layer — one line per skill)
npx skills add liyown/superpower-enterprise --skill java-code-reviewer
npx skills add liyown/superpower-enterprise --skill react-code-reviewer
npx skills add liyown/superpower-enterprise --skill go-code-reviewer
npx skills add liyown/superpower-enterprise --skill python-code-reviewer
npx skills add liyown/superpower-enterprise --skill node-code-reviewer
npx skills add liyown/superpower-enterprise --skill spec-doc-linter
npx skills add liyown/superpower-enterprise --skill goal-driven-development
npx skills add liyown/superpower-enterprise --skill project-knowledge-capture

# or install everything in the collection
npx skills add liyown/superpower-enterprise
```

Install is idempotent — re-running updates existing skills in place.

## Skills In This Collection (today)

| Skill | Trigger | Superpowers neighbour |
| --- | --- | --- |
| `java-code-reviewer` | reviewing a Java / Spring / MyBatis / Kafka / Reactor change | — |
| `react-code-reviewer` | reviewing a React / TypeScript / Next.js / Vite change | — |
| `go-code-reviewer` | reviewing a Go / gRPC / sqlx / GORM change | — |
| `python-code-reviewer` | reviewing a Python / asyncio / SQLAlchemy / FastAPI change | — |
| `node-code-reviewer` | reviewing a Node.js / Express / Fastify / Koa / Hono / Prisma / TypeORM / Sequelize / Knex change | — |
| `goal-driven-development` | turning an existing spec into code with verification, review gates, and knowledge capture | `executing-plans` |
| `project-knowledge-capture` | done with a task and want durable post-impl knowledge persisted into `docs/knowledge/` | — |
| `spec-doc-linter` | syncing a module's DevAgent.md or a domain's CONTEXT.md with the code | — |

## Planned Skills (v0.6.1+, see `docs/rfcs/2026-06-10-enterprise-flow-skills.md`)

| Planned skill | Trigger | Superpowers neighbour |
| --- | --- | --- |
| `data-lineage-trace` | where does a regulated field come from / who reads it / deletion propagation | (none) |
| `code-ownership-impact` | rename / deprecate / schema change / cross-team / CODEOWNERS | `writing-plans` |
| `incident-call-trace` | incident / outage / postmortem / on-call / SLO breach | `systematic-debugging` |
| `progressive-rollout-checklist` | feature flag / canary / dark launch / rollout / rollback | `finishing-a-development-branch` |
| `compliance-control-walk` | SOX / ISO / SOC 2 / GDPR / HIPAA / PCI / audit / control walkthrough | `verification-before-completion` |
| `cross-team-rfc-draft` | RFC / design doc / cross-team review / dependency request | `brainstorming` → `writing-plans` |
| `observability-coverage-audit` | blind spot / missing alert / coverage audit / gap analysis | `test-driven-development` |
| `pr-authoring` | 写 PR / PR 描述 / 整理 commit / 准备合入 | `requesting-code-review` |

When the planned skills land, this table is updated and the per-skill `npx skills add` line in `prompts/inventory.md` is filled in.

## When To Load This Skill

- An agent is being set up for the first time and the operator wants "everything enterprise".
- A consumer is choosing between two skills that look adjacent and wants a quick answer to "which one?" plus the install command.
- A consumer is writing a new skill in this collection and wants to know which neighbours it must declare cross-references for.
- A consumer wants the up-to-date list of "what is shipped" vs "what is planned" before deciding which release to pin (`#v0.5.0` vs `#v0.6.1`).

## Examples

There are no `bad-*` / `good-*` pairs for this meta-skill. The "output" of this skill is the inventory itself; the consumer reads the table to decide which concrete skill to load. See `prompts/inventory.md` for the per-skill one-liner with the install command.

## See Also

- The root `README.md` for the combined install + the full Relationship to superpowers table.
- `docs/rfcs/2026-06-10-enterprise-flow-skills.md` for the 8 planned skills and the boundary each draws vs superpowers.
- `docs/knowledge/2026-06-10-v0.6.0-rename.md` for the rename rationale and migration from the pre-v0.6.0 `liyown/skills-registry` repo path.
