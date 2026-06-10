# Changelog

All notable changes to this skills registry are documented here.
Versions follow [Semantic Versioning](https://semver.org/). The
source of truth is the git tag; this file mirrors the release
history for offline review.

## Unreleased

- New skill `spec-doc-linter` for keeping `DevAgent.md` (per module) and `CONTEXT.md` (per domain folder) in sync with the code. Detects mechanical drift (renamed symbols, removed files, broken links, missing dependencies) and LLM-judgment drift (weakened invariants, anti-patterns newly violated, false concurrency / idempotency claims, broken domain boundaries, stale context maps). Auto-sync uses per-file y/n/q confirmation, no `--fix-all`. Sentinel: `未发现文档与代码漂移。`.
- Dropped `scripts/{validate.sh, smoke.sh, check-examples.sh, release.sh}` and `.github/workflows/release-readiness.yml`. The repo is text only — there is no build step, no test runner, and no linter. Contributor contract is enforced by review against `CONTRIBUTING.md` rather than by a script.

## 0.3.9 (2026-06-08)

Tag: `v0.3.9`. Commits since v0.3.8: 1.

### Changed

- Router strengthening extended to the two meta-skills.
  `goal-driven-development` description now ends with "minimizing
  the risk of a reverted merge, a missed dependency, or a
  code-review block at PR time". `project-knowledge-capture`
  description now ends with "so the next developer (or the
  next session) can answer 'why was this done' and 'what to be
  careful about' in minutes, not days".
- All 7 skills now end their frontmatter `description` with a
  severity-tier router clause. The router strength was rolled
  out in 3 commits: v0.3.8 for the 5 reviewers, v0.3.9 for the
  2 meta-skills.

### Installation

```sh
npx skills add liyown/skills-registry#v0.3.9 --skill java-code-reviewer
```

## 0.3.8 (2026-06-08)

Tag: `v0.3.8`. Commits since v0.3.7: 2.

### Changed

- All 5 reviewer `SKILL.md` frontmatter `description` fields
  extended with a router clause that names the severity tier:
  "focusing on capital-loss, data-corruption, authorization-
  bypass, and incident-causing defects" (and the frontend
  equivalent "user-data-leak, XSS, authz-bypass" for
  react-code-reviewer). The trigger list now also includes
  "any review where a production incident is the worst case",
  matching the router style documented in CONTRIBUTING.md
  "What Makes A Great Skill" #1.

### Verification

- This release was cut with the new `--bump-mode minor` flag
  in `scripts/release.sh` (v0.3.7 → v0.4.0; patch reset to 0).
  The dry-run preview printed the would-be tag and confirmed
  the patch component reset, before the script applied the
  patch as v0.3.8.

### Installation

```sh
npx skills add liyown/skills-registry#v0.3.8 --skill java-code-reviewer
```

## 0.3.7 (2026-06-08)

Tag: `v0.3.7`. Commits since v0.3.6: 1.

### Changed

- `scripts/release.sh` adds `--bump-mode patch|minor|major`. The
  default is `patch` (preserves prior behaviour). `minor` resets
  the patch component to 0 and increments the minor; `major`
  resets both minor and patch to 0 and increments the major.
  `--bump X.Y.Z` still works and overrides the mode for
  non-monotonic cases. The complete release flag surface is
  now 7 flags.

### Installation

```sh
npx skills add liyown/skills-registry#v0.3.7 --skill java-code-reviewer
```

## 0.3.6 (2026-06-08)

Tag: `v0.3.6`. Commits since v0.3.5: 3.

### Added

- `.github/ISSUE_TEMPLATE/bug_report.md` — structured bug report
  with pre-flight checklist pointing at `validate.sh`,
  `CONTRIBUTING.md`, and `docs/releases/INDEX.md`.
- `.github/ISSUE_TEMPLATE/feature_request.md` — structured
  proposal for new skills or scenario prompts, with a pre-flight
  pointing at the Coverage Matrix and the "What Makes A Great
  Skill" section.
- `.github/ISSUE_TEMPLATE/config.yml` — disables blank issues
  and surfaces the release index, contributing guide, and
  "What Makes A Great Skill" as contact links.
- `.github/PULL_REQUEST_TEMPLATE.md` — 8-item self-check
  covering `validate.sh` exit, bad/good pairing, README sync,
  Coverage Matrix update, and `release.sh --dry-run` preview.

### Installation

```sh
npx skills add liyown/skills-registry#v0.3.6 --skill java-code-reviewer
```

## 0.3.5 (2026-06-08)

Tag: `v0.3.5`. Commits since v0.3.4: 1.

### Added

- `docs/releases/INDEX.md` — cross-reference for every release.
  Lists tag, release date, commit count since the previous tag,
  GitHub release URL, and the in-repo mirror file path. Also
  documents the full release flow (the `release.sh` flags and
  the `release-readiness` CI workflow). v0.3.0 and earlier have
  no mirror (pre-script era); v0.3.1+ have mirrors.

### Installation

```sh
npx skills add liyown/skills-registry#v0.3.5 --skill java-code-reviewer
```

## 0.3.4 (2026-06-08)

Tag: `v0.3.4`. Commits since v0.3.3: 2.

### Changed

- `scripts/release.sh` extended with two new flags:
  - `--notes-from <file>` validates the file exists and prints the
    exact `gh release create ...` command to run after the tag is
    pushed. Replaces the brittle inline `gh release create --notes "..."`
    flow that broke on shell-globbed content during the v0.3.3
    release.
  - `--no-publish` creates the tag locally without pushing. Useful
    for staged releases and for rehearsing the dry-run before the
    real push.
  - Both flags are reflected in the `--dry-run` preview so a
    consumer can sanity-check the would-be command before running
    with `--yes`.

### Installation

```sh
npx skills add liyown/skills-registry#v0.3.4 --skill java-code-reviewer
```

## 0.3.3 (2026-06-08)

Tag: `v0.3.3`. Commits since v0.3.2: 5.

### Added

- `.github/workflows/release-readiness.yml` — two-stage CI gate
  on every push / PR to `main`:
  1. `validate` runs `./scripts/validate.sh` (smoke + examples +
     full structural validation).
  2. `release-dry-run` runs `./scripts/release.sh --dry-run` and
     asserts no new tag exists at HEAD, so a regression in the
     release flow cannot pass CI unnoticed.
- `pr-diff-example.diff` added to `node-code-reviewer`,
  `python-code-reviewer`, and `react-code-reviewer`. All five
  reviewers now ship the same example asset surface: whole-file
  `bad-service.*` + `good-service.*`, per-scenario `bad-*` /
  `good-*` pairs, `review-output.md`, and `pr-diff-example.diff`.
- README "Skills at a Glance" routing tables — two tables
  (Reviewers, Workflows) that mirror the first clause of each
  skill's `description` field, so consumers can pick the right
  skill without parsing every frontmatter.

### Changed

- `actions/checkout` bumped from `v4` to `v5` in
  `.github/workflows/release-readiness.yml` to silence the
  Node.js 20 deprecation warning and use Node 24 by default.
- `CONTRIBUTING.md` already has the "What Makes A Great Skill"
  section (added in 0.3.2); no further change here.

### Installation

```sh
npx skills add liyown/skills-registry#v0.3.3 --skill java-code-reviewer
```

## 0.3.2 (2026-06-08)

Tag: `v0.3.2`. Commits since v0.3.1: 7.

### Added

- `react-code-reviewer` expanded from 4 to 10 scenario prompts, closing
  the frontend coverage gap relative to the backend reviewers:
  - `testing-reviewer.md` — Jest, Vitest, React Testing Library,
    Playwright, Cypress; `act()` wrapping, mock scoping, behaviour
    vs implementation assertions.
  - `forms-reviewer.md` — react-hook-form, Formik, Zod, server
    actions; submit disable, double-submit, schema-driven errors.
  - `state-reviewer.md` — Redux Toolkit, Zustand, Jotai, Context;
    selector memoisation, store splitting, SSR hydration.
  - `a11y-reviewer.md` — keyboard, screen reader, ARIA, focus
    management, contrast, live regions, reduced motion.
  - `error-boundary-reviewer.md` — top-level vs route-level
    boundaries, event-handler async errors, fallback UX, log
    forwarding.
  - `bundle-reviewer.md` — tree-shaking, dynamic imports, route-
    level code splitting, icon imports, barrel files.
- 6 paired `bad-*` / `good-*` example `.tsx` files for the new
  react scenarios, with `Fix N:` annotations matching the rest of
  the collection.
- `docs:` "What Makes A Great Skill" section in `CONTRIBUTING.md`
  capturing the patterns used by `anthropics/skills` and
  `vercel-labs/agent-skills` (description as router, progressive
  disclosure, mandatory bad/good pairs, fallback contracts, etc.).

### Changed

- `react-code-reviewer/SKILL.md` frontmatter `description` now lists
  the additional scenario keywords so the consumer routes correctly.
- `docs/CHANGELOG.md` re-created (deleted during the shadcn drop)
  with the cumulative v0.2.0 → v0.3.2 history.
- README "Reviewer Coverage Matrix" expanded from 7 to 11 rows to
  host the new react scenarios.

### Installation

```sh
npx skills add liyown/skills-registry#v0.3.2 --skill java-code-reviewer
```

## 0.3.1 (2026-06-08)

Tag: `v0.3.1`. Commits since v0.3.0: 4.

### Added

- `scripts/release.sh` — tag-driven release flow. Reads `git describe`,
  bumps the patch component, runs `./scripts/validate.sh`, creates
  an annotated tag, and pushes both the branch and the tag to
  origin. Supports `--dry-run` (preview without mutating) and
  `--bump X.Y.Z` (override the auto-bumped version).
- `docs:` Reviewer Coverage Matrix footer links to the per-prompt
  `> See also:` lines so consumers can navigate the matrix index
  and the in-prompt cross-references from either side.
- `docs:` README "Latest Release" section now points at the current
  tag and includes a `git log` hint for the full diff.
- `docs/releases/v0.3.1.md` — in-repo mirror of the GitHub release
  notes (so reviewers do not need GitHub API access to inspect the
  release).

### Installation

```sh
npx skills add liyown/skills-registry#v0.3.1 --skill java-code-reviewer
```

### Verification at tag time

- `./scripts/validate.sh` exit=0:
  - `smoke check passed`
  - `examples check passed`
  - `validation passed`
- `./scripts/release.sh --dry-run` prints the would-be actions
  without creating a tag; verified both default patch bump and
  explicit `--bump 0.4.0`.

## 0.3.0 (2026-06-08)

Tag: `v0.3.0`. Released: <https://github.com/liyown/skills-registry/releases/tag/v0.3.0>.

### Added

- Drop the shadcn-style registry infrastructure (`registry.json`,
  per-skill `manifest.json` / `registry.json` / `agents/openai.yaml`,
  `scripts/build.sh` / `scripts/validate.sh` /
  `scripts/validate-registry.mjs`, `docs/adding-skills.md`,
  `.github/workflows/build.yml`, `package.json`) in favour of a
  pure `npx skills` collection.
- `go-code-reviewer` — evidence-driven Go backend review (7
  prompts, 12 examples).
- `python-code-reviewer` — Python backend review (6 prompts, 11
  examples).
- `node-code-reviewer` — Node.js backend review (6 prompts, 11
  examples).
- `scripts/smoke.sh`, `scripts/check-examples.sh`,
  `scripts/validate.sh` — offline structural assertions (no
  external toolchains).
- `CONTRIBUTING.md` — frontmatter schema, body conventions, prompt
  /example naming, bad/good pairing rules, fallback contracts,
  quality bar.
- `docs/CHANGELOG.md` — this file.
- Per-prompt `> See also:` cross-references in all 33 prompts
  across 7 skills.

### Changed

- All skill `SKILL.md` / `README.md` / `prompts/*.md` / `examples/`
  bodies unified to English (was previously mixed Chinese/English).
  Three contract phrases retained in their original language:
  `未发现明确高风险问题。`, `需要结合上下文确认`,
  `CodeGraph unavailable; context was gathered by rg/file inspection.`
- `goal-driven-development` now lists `go-code-reviewer`,
  `python-code-reviewer`, `node-code-reviewer` alongside the
  existing `java-code-reviewer` and `react-code-reviewer`.
- Root `README.md` rewritten for the `npx skills` install model
  with a Cross-Skill Dependencies section and a Reviewer Coverage
  Matrix.

### Installation

```sh
npx skills add liyown/skills-registry#v0.3.0 --skill java-code-reviewer
```

## 0.2.0

Initial shadcn-compatible registry release. Four skills:
`java-code-reviewer`, `react-code-reviewer`,
`goal-driven-development`, `project-knowledge-capture`. Single
`SKILL.md` entrypoint per skill, manifest-driven file lists,
nested `registry.json` per skill. Released prior to the move to
`npx skills`.
