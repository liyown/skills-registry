# Superpower Enterprise

The **enterprise augment layer** for [`obra/superpowers`](https://github.com/obra/superpowers). Like `mybatis-plus` is to `mybatis`, this collection adds the production-grade / org-grade / regulator-grade skills that `superpowers` deliberately refuses to assume: code-review depth, spec linting, dev-flow orchestration, and (in upcoming releases) data lineage, CODEOWNERS-aware blast-radius, incident trace, compliance control walks, and progressive-rollout checklists.

Use this collection **alongside** `superpowers`, not instead of it. `superpowers` covers the generic author-craft (TDD, debugging, brainstorming, requesting/receiving code review, finishing a branch). This collection covers the **language-specific production-risk review** (5 reviewers), the **spec-to-code orchestration** (`goal-driven-development`), the **durable post-impl knowledge** (`project-knowledge-capture`), and the **doc drift linter** (`spec-doc-linter`). The 8 enterprise-flow skills proposed in [`docs/rfcs/2026-06-10-enterprise-flow-skills.md`](./docs/rfcs/2026-06-10-enterprise-flow-skills.md) land in v0.6.1+.

## Latest Release

**v0.5.0** (2026-06-10). Pin to a specific tag with `#v0.5.0`:

```sh
npx skills add liyown/superpower-enterprise#v0.5.0 --skill goal-driven-development
```

The v0.5.0 release ships the slim `goal-driven-development` orchestrator rework + a fix to a stale `codegraph.md` prompt (dropped the now-removed `codegraph_context` / `codegraph_trace` tool names; `codegraph_explore` is the primary walk tool now) + the [enterprise-flow-skills RFC](./docs/rfcs/2026-06-10-enterprise-flow-skills.md) listing the 8 skills planned for the next release. No new skills ship in v0.5.0.

> **Pre-v0.6.0 history.** Releases v0.3.0 – v0.5.0 were published under the original repo path `liyown/skills-registry`. They remain reachable on GitHub via the `superpower-enterprise` rename redirect. If your existing install commands reference `liyown/skills-registry`, update them to `liyown/superpower-enterprise` — the installed skill folders are identical, only the URL changed. See [`docs/knowledge/2026-06-10-v0.6.0-rename.md`](./docs/knowledge/2026-06-10-v0.6.0-rename.md) for the rename rationale and migration.

## Install

The combined install pulls `superpowers` and this collection in one line per skill, so a consumer never has to remember which collection a skill came from:

```sh
# install superpowers + the enterprise collection in one shot
npx skills add obra/superpowers
npx skills add liyown/superpower-enterprise --skill java-code-reviewer
npx skills add liyown/superpower-enterprise --skill react-code-reviewer
npx skills add liyown/superpower-enterprise --skill go-code-reviewer
npx skills add liyown/superpower-enterprise --skill python-code-reviewer
npx skills add liyown/superpower-enterprise --skill node-code-reviewer
npx skills add liyown/superpower-enterprise --skill spec-doc-linter
npx skills add liyown/superpower-enterprise --skill goal-driven-development
npx skills add liyown/superpower-enterprise --skill project-knowledge-capture

# install everything
npx skills add liyown/superpower-enterprise
```

`npx skills` reads the `SKILL.md` frontmatter in each subdirectory and copies the whole skill folder into the consumer's local skills directory (`~/.claude/skills/<name>/` for Claude Code, etc.). Existing skills are updated in place on a re-install.

## Cross-Skill Dependencies

`npx skills` does not auto-install dependencies referenced inside a skill's body. When a skill says "invoke `<other-skill>` for X", the consumer must install that other skill separately. Concrete cases in this collection:

- `goal-driven-development` references `java-code-reviewer`, `react-code-reviewer`, `go-code-reviewer`, `python-code-reviewer`, `node-code-reviewer`, and `project-knowledge-capture` as runtime helpers. Install all six alongside it:

  ```sh
  npx skills add liyown/superpower-enterprise \
    --skill goal-driven-development \
    --skill java-code-reviewer \
    --skill react-code-reviewer \
    --skill go-code-reviewer \
    --skill python-code-reviewer \
    --skill node-code-reviewer \
    --skill project-knowledge-capture
  ```

  Plus `obra/superpowers` if the consumer also wants the generic author-craft:

  ```sh
  npx skills add obra/superpowers
  ```

- The 8 enterprise-flow skills proposed in the RFC (data-lineage, code-ownership, incident-call, progressive-rollout, compliance-control, cross-team-RFC, observability-coverage, pr-authoring) will each declare their own `superpowers` peer dependency when they ship. See the RFC for the per-skill peer list.

## Included Skills

- `java-code-reviewer` — evidence-driven Java backend production-risk review.
- `react-code-reviewer` — React / TypeScript / Next.js frontend production-risk review.
- `go-code-reviewer` — Go backend production-risk review (goroutines, context, errors, sqlx, gRPC).
- `python-code-reviewer` — Python backend production-risk review (asyncio, error handling, SQLAlchemy/Django, security).
- `node-code-reviewer` — Node.js backend production-risk review (async, error handling, Prisma, Express/Fastify, security).
- `goal-driven-development` — CodeGraph-assisted implementation workflow for existing specs/goals.
- `project-knowledge-capture` — durable project knowledge capture into `docs/knowledge/`.
- `spec-doc-linter` — DevAgent.md / CONTEXT.md drift detection with per-file auto-sync.

## Skills at a Glance

Pick the skill whose description best matches what you are doing. The first sentence of each `description` is the router the consumer uses to decide whether to load the skill.

### Reviewers (load when reviewing a diff)

| Skill | Load when you are … |
| --- | --- |
| `java-code-reviewer` | … reviewing a Java / Spring / MyBatis / Kafka / Reactor change |
| `react-code-reviewer` | … reviewing a React / TypeScript / Next.js / Vite change (frontend, not Node backend) |
| `go-code-reviewer` | … reviewing a Go / gRPC / sqlx / GORM change |
| `python-code-reviewer` | … reviewing a Python / asyncio / SQLAlchemy / FastAPI change |
| `node-code-reviewer` | … reviewing a Node.js / Express / Fastify / Koa / Hono / Prisma / TypeORM / Sequelize / Knex change |

Each reviewer ships a `prompts/reviewer.md` (severity ladder + output contract) and 6-10 scenario-specific prompts (see the Coverage Matrix below). Bad/good `.java` / `.tsx` / `.go` / `.py` / `.ts` example pairs in `examples/` show the minimum fix for the most common production-risk findings.

### Workflows (load when executing a goal)

| Skill | Load when you are … |
| --- | --- |
| `goal-driven-development` | … turning an existing spec into code with verification, review gates, and knowledge capture |
| `project-knowledge-capture` | … done with a task and want durable project knowledge persisted into `docs/knowledge/` |

### Tools (load when keeping docs in sync)

| Skill | Load when you are … |
| --- | --- |
| `spec-doc-linter` | … syncing a module's DevAgent.md or a domain's CONTEXT.md with the code (Tier-1 static + Tier-2 LLM judgment, per-file y/n/q confirmation) |

These are not reviewers; they are end-to-end workflows that may invoke the reviewers as helpers. See "Cross-Skill Dependencies" above for the combined install command for `goal-driven-development`.

## Reviewer Coverage Matrix

Each reviewer skill loads `prompts/reviewer.md` (core protocol) and one or more scenario-specific prompts. The matrix shows which scenarios each reviewer covers; cell entries are the scenario prompt file names.

| Scenario | java | react | go | python | node |
| --- | --- | --- | --- | --- | --- |
| Framework / runtime | `spring-reviewer.md` | `nextjs-reviewer.md` | `rpc-reviewer.md` | `web-reviewer.md` | `http-reviewer.md` |
| Concurrency / async | `concurrency-reviewer.md`, `reactor-reviewer.md` | — | `concurrency-reviewer.md` | `async-reviewer.md` | `async-reviewer.md` |
| Context propagation | — | — | `context-reviewer.md` | — | — |
| Error handling | — | `error-boundary-reviewer.md` | `error-reviewer.md` | `error-reviewer.md` | `error-reviewer.md` |
| Database / ORM | `mybatis-reviewer.md` | — | `sql-reviewer.md` | `sql-reviewer.md` | `sql-reviewer.md` |
| Caching / messaging | `redis-kafka-reviewer.md` | — | — | — | — |
| Security | `security-reviewer.md` | `security-reviewer.md` | `security-reviewer.md` | `security-reviewer.md` | `security-reviewer.md` |
| Performance / bundle | — | `performance-reviewer.md`, `bundle-reviewer.md` | — | — | — |
| Forms / validation | — | `forms-reviewer.md` | — | — | — |
| State management | — | `state-reviewer.md` | — | — | — |
| Accessibility (a11y) | — | `a11y-reviewer.md` | — | — | — |
| Testing | — | `testing-reviewer.md` | — | — | — |

A `—` cell means the reviewer has no dedicated scenario prompt for that category. Add a `<scenario>-reviewer.md` under the relevant `prompts/` folder and reference it from `SKILL.md` to close a gap.

> Cell entries are file names inside the row's `prompts/` directory; some scenario names (`security-reviewer`, `sql-reviewer`, `error-reviewer`, `concurrency-reviewer`) exist in multiple skills, each maintained for its target language.

> Each prompt in `prompts/` opens with a `> See also:` line listing the most relevant sibling prompts in the same skill, so the matrix here (skill × scenario index) and the per-prompt cross-references are navigable from both sides.

## Relationship to `obra/superpowers`

| Concern | `superpowers` | This collection |
| --- | --- | --- |
| TDD / debugging / brainstorming / writing-plans | `test-driven-development`, `systematic-debugging`, `brainstorming`, `writing-plans` | — (use `superpowers`) |
| Language-specific production-risk review | — | 5 `*-code-reviewer` skills |
| Spec → code orchestration | `executing-plans` | `goal-driven-development` (more end-to-end, less per-step) |
| Post-impl durable knowledge | — | `project-knowledge-capture` |
| Doc / spec drift detection | — | `spec-doc-linter` |
| Receiving / requesting code review | `requesting-code-review`, `receiving-code-review` | — (use `superpowers`); v0.6.1+ adds `pr-authoring` (中文 + 企业 CI) |
| Data lineage / PII / compliance | — | v0.6.1+: `data-lineage-trace`, `compliance-control-walk` |
| CODEOWNERS / cross-team blast radius | — | v0.6.1+: `code-ownership-impact` |
| Incident / postmortem / SLO | — | v0.6.1+: `incident-call-trace` |
| Feature flag / canary / rollout | — | v0.6.1+: `progressive-rollout-checklist` |
| Cross-team RFC | — | v0.6.1+: `cross-team-rfc-draft` |
| Observability coverage | — | v0.6.1+: `observability-coverage-audit` |

The full proposal for the 8 enterprise-flow skills is in [`docs/rfcs/2026-06-10-enterprise-flow-skills.md`](./docs/rfcs/2026-06-10-enterprise-flow-skills.md).

## Authoring A Skill

See [CONTRIBUTING.md](./CONTRIBUTING.md) for the full contract (frontmatter schema, `SKILL.md` body conventions, prompts/examples naming, bad/good pairing rules, quality bar).

Quick version:

1. Create `skills/<name>/` with a `SKILL.md` whose frontmatter `name` matches the directory and `description` is ≥ 40 chars.
2. Add `prompts/` (one scenario per file, progressive disclosure) and `examples/` (every `bad-*` paired with a `good-*`).
3. Update the "Included Skills" list above and any cross-skill dependency blocks that reference the new skill.
4. Open a PR. `npx skills` consumers pull from `main` directly, so the next install picks up new skills without version coordination.

## License

MIT
