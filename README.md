# Skills Registry

A collection of reusable AI skills installed with the [skills.sh](https://skills.sh)
CLI (`npx skills`). Each skill lives in its own directory under `skills/` and
is discoverable by name.

## Latest Release

**v0.3.9** (2026-06-08). Pin to a specific tag with `#v0.3.9`:

```sh
npx skills add liyown/skills-registry#v0.3.9 --skill java-code-reviewer
```

v0.3.8 → v0.3.9 extends the router-strengthening pattern to the
two meta-skills. `goal-driven-development` now ends with
"minimizing the risk of a reverted merge, a missed dependency,
or a code-review block at PR time". `project-knowledge-capture`
now ends with "so the next developer (or the next session) can
answer 'why was this done' and 'what to be careful about' in
minutes, not days". 1 commit since v0.3.8.

All 5 reviewer frontmatter descriptions end with a production-incident
router clause. The 2 workflow-skill descriptions end with a "Use when
..." clause scoped to the workflow trigger. 4 commits of router tuning
across v0.3.8 and v0.3.9.

See [docs/CHANGELOG.md](./docs/CHANGELOG.md) for the full release history
and per-version highlights. Run `git log v0.3.8..v0.3.9` for the full
diff since the last release.

## Install

```sh
# install one or more skills
npx skills add liyown/skills-registry --skill java-code-reviewer
npx skills add liyown/skills-registry --skill react-code-reviewer
npx skills add liyown/skills-registry --skill go-code-reviewer
npx skills add liyown/skills-registry --skill python-code-reviewer
npx skills add liyown/skills-registry --skill node-code-reviewer
npx skills add liyown/skills-registry --skill goal-driven-development
npx skills add liyown/skills-registry --skill project-knowledge-capture

# install everything
npx skills add liyown/skills-registry
```

`npx skills` reads the `SKILL.md` frontmatter in each subdirectory and copies
the whole skill folder into the consumer's local skills directory
(`~/.claude/skills/<name>/` for Claude Code, etc.).

## Cross-Skill Dependencies

`npx skills` does not auto-install dependencies referenced inside a skill's
body. When a skill says "invoke `<other-skill>` for X", the consumer must
install that other skill separately. Concrete cases in this collection:

- `goal-driven-development` references `java-code-reviewer`,
  `react-code-reviewer`, `go-code-reviewer`, `python-code-reviewer`,
  `node-code-reviewer`, and `project-knowledge-capture` as runtime
  helpers. Install all six alongside it:

  ```sh
  npx skills add liyown/skills-registry \
    --skill goal-driven-development \
    --skill java-code-reviewer \
    --skill react-code-reviewer \
    --skill go-code-reviewer \
    --skill python-code-reviewer \
    --skill node-code-reviewer \
    --skill project-knowledge-capture
  ```

  Install can be repeated; existing skills are updated in place.

## Included Skills

- `java-code-reviewer` — evidence-driven Java backend production-risk review.
- `react-code-reviewer` — React / TypeScript / Next.js frontend production-risk review.
- `go-code-reviewer` — Go backend production-risk review (goroutines, context, errors, sqlx, gRPC).
- `python-code-reviewer` — Python backend production-risk review (asyncio, error handling, SQLAlchemy/Django, security).
- `node-code-reviewer` — Node.js backend production-risk review (async, error handling, Prisma, Express/Fastify, security).
- `goal-driven-development` — CodeGraph-assisted implementation workflow for existing specs/goals.
- `project-knowledge-capture` — durable project knowledge capture into `docs/knowledge/`.

## Skills at a Glance

Pick the skill whose description best matches what you are doing.
The first sentence of each `description` is the router the consumer
uses to decide whether to load the skill.

### Reviewers (load when reviewing a diff)

| Skill | Load when you are … |
| --- | --- |
| `java-code-reviewer` | … reviewing a Java / Spring / MyBatis / Kafka / Reactor change |
| `react-code-reviewer` | … reviewing a React / TypeScript / Next.js / Vite change (frontend, not Node backend) |
| `go-code-reviewer` | … reviewing a Go / gRPC / sqlx / GORM change |
| `python-code-reviewer` | … reviewing a Python / asyncio / SQLAlchemy / FastAPI change |
| `node-code-reviewer` | … reviewing a Node.js / Express / Fastify / Koa / Hono / Prisma / TypeORM / Sequelize / Knex change |

Each reviewer ships a `prompts/reviewer.md` (severity ladder + output
contract) and 6-10 scenario-specific prompts (see the Coverage
Matrix below). Bad/good `.java` / `.tsx` / `.go` / `.py` / `.ts`
example pairs in `examples/` show the minimum fix for the most
common production-risk findings.

### Workflows (load when executing a goal)

| Skill | Load when you are … |
| --- | --- |
| `goal-driven-development` | … turning an existing spec into code with verification, review gates, and knowledge capture |
| `project-knowledge-capture` | … done with a task and want durable project knowledge persisted into `docs/knowledge/` |

These are not reviewers; they are end-to-end workflows that may
invoke the reviewers as helpers. See "Cross-Skill Dependencies"
above for the combined install command for `goal-driven-development`.

## Reviewer Coverage Matrix

Each reviewer skill loads `prompts/reviewer.md` (core protocol) and one or
more scenario-specific prompts. The matrix shows which scenarios each
reviewer covers; cell entries are the scenario prompt file names.

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

A `—` cell means the reviewer has no dedicated scenario prompt for that
category. Add a `<scenario>-reviewer.md` under the relevant `prompts/`
folder and reference it from `SKILL.md` to close a gap.

> Cell entries are file names inside the row's `prompts/` directory; some
> scenario names (`security-reviewer`, `sql-reviewer`, `error-reviewer`,
> `concurrency-reviewer`) exist in multiple skills, each maintained for
> its target language.

> Each prompt in `prompts/` opens with a `> See also:` line listing the
> most relevant sibling prompts in the same skill, so the matrix here
> (skill × scenario index) and the per-prompt cross-references are
> navigable from both sides.

## Layout

```text
.
├── README.md
├── CONTRIBUTING.md
├── LICENSE
└── skills/
    ├── java-code-reviewer/
    ├── react-code-reviewer/
    ├── go-code-reviewer/
    ├── python-code-reviewer/
    ├── node-code-reviewer/
    ├── goal-driven-development/
    └── project-knowledge-capture/
```

Each skill folder contains:

```text
skills/<name>/
├── SKILL.md            # required: frontmatter (name, description) + body
├── README.md           # human-readable description
├── prompts/            # scenario-specific prompt fragments (loaded on demand)
└── examples/           # bad/good code samples and review outputs
```

## Authoring A Skill

See [CONTRIBUTING.md](./CONTRIBUTING.md) for the full contract
(frontmatter schema, `SKILL.md` body conventions, prompts/examples
naming, bad/good pairing rules, quality bar).

Quick version:

1. Create `skills/<name>/` with a `SKILL.md` whose frontmatter
   `name` matches the directory and `description` is ≥ 40 chars.
2. Add `prompts/` (one scenario per file, progressive disclosure) and
   `examples/` (every `bad-*` paired with a `good-*`).
3. Update the "Included Skills" list above and any cross-skill
   dependency blocks that reference the new skill.
4. Open a PR. `npx skills` consumers pull from `main` directly, so the
   next install picks up new skills without version coordination.

## Local Checks

This repo ships text only — there is no build step, no test runner,
and no linter. A skill is correct if it is a self-contained folder
that an agent can copy verbatim. Review the contributor contract in
[CONTRIBUTING.md](./CONTRIBUTING.md) for the frontmatter, prompt, and
example-pairing rules that govern a merge.

## License

MIT
