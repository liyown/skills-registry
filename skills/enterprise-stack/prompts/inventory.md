# Enterprise Stack Inventory

> See also: ../SKILL.md, ../../docs/rfcs/2026-06-10-enterprise-flow-skills.md

One line per skill, with the install command. The router (../SKILL.md) decides *which* skill to load; this file is the "what is its name and how do I get it" reference.

## Shipped Today (v0.6.1)

| Skill | One-line summary | Install |
| --- | --- | --- |
| `java-code-reviewer` | Java / Spring / MyBatis / Kafka / Reactor production-risk review. | `npx skills add liyown/superpower-enterprise --skill java-code-reviewer` |
| `react-code-reviewer` | React / TypeScript / Next.js / Vite production-risk review. | `npx skills add liyown/superpower-enterprise --skill react-code-reviewer` |
| `go-code-reviewer` | Go / gRPC / sqlx / GORM production-risk review. | `npx skills add liyown/superpower-enterprise --skill go-code-reviewer` |
| `python-code-reviewer` | Python / asyncio / SQLAlchemy / FastAPI production-risk review. | `npx skills add liyown/superpower-enterprise --skill python-code-reviewer` |
| `node-code-reviewer` | Node.js / Express / Fastify / Koa / Hono / Prisma / TypeORM / Sequelize / Knex production-risk review. | `npx skills add liyown/superpower-enterprise --skill node-code-reviewer` |
| `goal-driven-development` | CodeGraph-assisted spec → code orchestration. | `npx skills add liyown/superpower-enterprise --skill goal-driven-development` |
| `project-knowledge-capture` | Durable post-impl knowledge into `docs/knowledge/`. | `npx skills add liyown/superpower-enterprise --skill project-knowledge-capture` |
| `spec-doc-linter` | DevAgent.md / CONTEXT.md drift detection with per-file y/n/q. | `npx skills add liyown/superpower-enterprise --skill spec-doc-linter` |
| `enterprise-stack` | Single install point that inventories the whole collection. | `npx skills add liyown/superpower-enterprise --skill enterprise-stack` |
| `data-lineage-trace` | PII / PHI / regulated-field lineage end to end. | `npx skills add liyown/superpower-enterprise --skill data-lineage-trace` |
| `code-ownership-impact` | Cross-team blast radius + CODEOWNERS for shared-kernel changes. | `npx skills add liyown/superpower-enterprise --skill code-ownership-impact` |
| `incident-call-trace` | End-to-end call tracing for incident response and postmortems. | `npx skills add liyown/superpower-enterprise --skill incident-call-trace` |
| `progressive-rollout-checklist` | Rollout-readiness, flag-coverage, rollback-path verification. | `npx skills add liyown/superpower-enterprise --skill progressive-rollout-checklist` |
| `compliance-control-walk` | Evidence collection for SOX / ISO / GDPR / HIPAA / PCI controls. | `npx skills add liyown/superpower-enterprise --skill compliance-control-walk` |
| `observability-coverage-audit` | Gap audit: which error paths lack tests, alerts, or runbooks. | `npx skills add liyown/superpower-enterprise --skill observability-coverage-audit` |
| `cross-team-rfc-draft` | Cross-team RFC drafting with named approvers and migration windows. | `npx skills add liyown/superpower-enterprise --skill cross-team-rfc-draft` |
| `pr-authoring` | PR 准备: 描述、reviewer 名单、verification 证据、跨仓库引用. | `npx skills add liyown/superpower-enterprise --skill pr-authoring` |

Plus the collection-level install (all of the above):

```sh
npx skills add liyown/superpower-enterprise
```

## Combined Install (one block, idempotent)

```sh
# superpowers
npx skills add obra/superpowers

# superpower-enterprise — install one or more of the 17 shipped skills
npx skills add liyown/superpower-enterprise --skill java-code-reviewer
npx skills add liyown/superpower-enterprise --skill react-code-reviewer
npx skills add liyown/superpower-enterprise --skill go-code-reviewer
npx skills add liyown/superpower-enterprise --skill python-code-reviewer
npx skills add liyown/superpower-enterprise --skill node-code-reviewer
npx skills add liyown/superpower-enterprise --skill spec-doc-linter
npx skills add liyown/superpower-enterprise --skill goal-driven-development
npx skills add liyown/superpower-enterprise --skill project-knowledge-capture
npx skills add liyown/superpower-enterprise --skill enterprise-stack
npx skills add liyown/superpower-enterprise --skill data-lineage-trace
npx skills add liyown/superpower-enterprise --skill code-ownership-impact
npx skills add liyown/superpower-enterprise --skill incident-call-trace
npx skills add liyown/superpower-enterprise --skill progressive-rollout-checklist
npx skills add liyown/superpower-enterprise --skill compliance-control-walk
npx skills add liyown/superpower-enterprise --skill observability-coverage-audit
npx skills add liyown/superpower-enterprise --skill cross-team-rfc-draft
npx skills add liyown/superpower-enterprise --skill pr-authoring

# or all at once
npx skills add liyown/superpower-enterprise
```

Re-run any of the lines above to update an installed skill in place.

## Pre-v0.6.0 Migration

Releases v0.3.0 – v0.5.0 were published under the original repo path `liyown/skills-registry`. The repo was renamed to `liyown/superpower-enterprise` in v0.6.0. If your existing install commands reference `liyown/skills-registry`, update them to `liyown/superpower-enterprise` — the installed skill folders are identical, only the URL changed. GitHub's rename redirect handles `#v0.3.0` – `#v0.5.0` pin URLs automatically.

v0.6.1+ install lines use the new path; no migration needed.
