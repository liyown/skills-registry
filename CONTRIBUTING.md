# Contributing

This repository is a collection of skills consumed via the
[`npx skills`](https://skills.sh) CLI. Each skill is a self-contained
folder that gets copied into the consumer's local skills directory
(`~/.claude/skills/<name>/` for Claude Code, etc.). There is no registry
metadata, version coordination, or build step — the folder layout is the
contract.

## Directory Layout

```text
skills/<skill-name>/
├── SKILL.md            # required entrypoint
├── README.md           # human-readable description (optional but recommended)
├── prompts/            # scenario-specific prompt fragments (loaded on demand)
└── examples/           # bad/good code samples and review outputs
```

A new skill must live under its own directory. The directory name is the
canonical skill name and must match `SKILL.md` frontmatter `name`.

## `SKILL.md` Frontmatter

Required fields:

| Field | Type | Constraint |
| --- | --- | --- |
| `name` | string | Must equal the parent directory name. Lowercase, digits, dashes only. |
| `description` | string | ≥ 40 chars. Should describe when to invoke the skill, not just what it is. Include the trigger phrases users are likely to type. |

Optional but recommended:

| Field | Type | Purpose |
| --- | --- | --- |
| `metadata.short-description` | string | One-line summary used by some consumers' listings. |

Frontmatter parsing rules consumers apply:

- The body starts with `---` on its own line and ends with another `---`.
- Quoted strings with `:` or leading whitespace are preserved verbatim.
- Anything below the closing `---` is the body.

## `SKILL.md` Body

The body is what agents see when the skill is invoked. Keep it small and
focused on routing:

1. **Required loading** — which `prompts/*.md` are always loaded
2. **Optional loading** — which `prompts/*.md` to load for specific scenarios
3. **Review/output contract** — the exact output format the agent must produce
4. **Examples pointer** — point to `examples/bad-*` and `examples/good-*` for
   reference cases

Hard rules:

- Use **progressive disclosure**: keep `SKILL.md` under ~80 lines, push detail
  into `prompts/`.
- Each `prompts/<file>.md` covers one scenario. Multiple scenarios in one file
  make selective loading impossible.
- Use the exact no-finding sentence `未发现明确高风险问题。` for any reviewer
  skill that uses the same severity ladder.

## `prompts/` Conventions

- File names describe a single scenario, e.g. `concurrency-reviewer.md`,
  `error-reviewer.md`, `codegraph.md`.
- Each file starts with a one-line H1 naming the scope
  (e.g. `# Java Concurrency Reviewer Prompt`).
- Each file's `Required Checks` / `Fallback` section is the
  authoritative checklist for that scenario.
- If a prompt depends on an external tool that may be missing
  (e.g. CodeGraph), include an explicit **fallback contract**:
  - Trigger condition (how to detect unavailability)
  - Fallback order
  - The exact fallback line the agent must emit
    (e.g. `CodeGraph unavailable; context was gathered by rg/file inspection.`)
  - Stop-loss rules (when to stop trying)

## `examples/` Conventions

Each `bad-<file>` should have a matching `good-<file>` that shows the
minimum fix for every Critical/High finding. The pairing is the
primary teaching asset — a bad example without its good counterpart
is incomplete.

- `examples/bad-*.{java,tsx,ts,go,...}` — annotated with the issue
  category it demonstrates
- `examples/good-*.{java,tsx,ts,go,...}` — minimal fix; top-of-file
  comment lists the issue IDs it addresses
- `examples/review-output.md` — full reviewer output for the bad
  example, formatted exactly as a real agent should produce
- `examples/pr-diff-example.diff` — minimal PR-shaped diff that
  reproduces the bad code; optional but helpful for training
- `examples/workflow-note.md` / `examples/knowledge-note.md` — for
  workflow skills, a sample of the artifact the skill should produce

## Cross-Skill References

A skill may reference other skills by name (e.g.
`java-code-reviewer for Java backend changes`). The consumer must
install referenced skills separately; `npx skills` does not transitively
install dependencies. Document any required peers in the skill body and
in this repository's `README.md` under "Cross-Skill Dependencies".

## Cross-Reference to Other Skills

When a new skill is part of an existing workflow, add it to:

- The root `README.md` "Included Skills" list
- Any other skill that references it by name in its body
- Any "Cross-Skill Dependencies" section that lists install commands

## Proposing a Skill

1. Create the directory and `SKILL.md` first; get the frontmatter and
   "Required Loading" right before adding prompts.
2. Add at least one `examples/bad-<file>` + `examples/good-<file>` pair
   that demonstrates the skill's primary use case.
3. Add `examples/review-output.md` (or equivalent) showing the expected
   output for the bad example.
4. Update the root `README.md` skill list and any cross-skill dependency
   blocks that reference the new skill.
5. Open a PR. The `npx skills` consumer pulls from `main` directly, so
   the next install picks up new skills without version coordination.

## Quality Bar

A skill is ready to merge when:

- `SKILL.md` frontmatter is valid (name matches directory, description ≥ 40 chars)
- `SKILL.md` body is under ~80 lines and references prompts by relative path
- Each `prompts/*.md` covers one scenario and is the only authority for that scenario's checklist
- Every `examples/bad-*` has a matching `examples/good-*`
- `examples/review-output.md` matches the contract declared in `SKILL.md`
- Cross-skill references resolve to skills that exist in this repository
- New skills appear in root `README.md` and in any dependent skill's docs

## What Makes A Great Skill

The "Quality Bar" above is the floor. The ceiling is what separates
a useful skill from one a consumer will reach for again. These
guidelines are distilled from the patterns used by
`anthropics/skills`, `vercel-labs/agent-skills`, and similar
collections.

### 1. Description is a router, not a description

The `description` field is the only signal the consumer uses to
decide whether to load the skill. It is **not** a marketing
sentence about what the skill is — it is a router that says "load
me when the user is doing X". A good description:

- Starts with the user-facing trigger ("Review this Java PR for
  production-risk issues…").
- Lists the technologies, frameworks, or languages the skill
  covers, separated by commas.
- Includes the verbs users actually type: "review", "PR review",
  "diff review", "bug-risk review", "release-blocking inspection".
- Avoids internal jargon that an end user would never say.

A bad description is "A skill for Java code review" — vague,
trigger-free, and indistinguishable from a dozen other skills.

### 2. SKILL.md body is a router, prompts/* are the substance

The skill body is loaded every time the skill is invoked. Anything
longer than ~80 lines is loaded into context repeatedly and
costs the user tokens. Move the substance into `prompts/*.md`,
one scenario per file, with progressive disclosure. A consumer
reading the body should be able to decide which scenario files
to load without scanning the whole skill.

### 3. Bad and good are not optional

A bad example without a good counterpart is half a lesson. The
pair is the primary teaching asset. For every issue a reviewer
prompt calls out, the consumer should be able to find a
`bad-<file>` plus a `good-<file>` that demonstrates the minimum
fix, with `Fix N:` annotations explaining each change.

The two `bad-` / `good-` filenames must match stem-to-stem:

```text
examples/bad-wrapper-injection.java   examples/good-wrapper-injection.java
examples/bad-goroutine-leak.go       examples/good-goroutine-leak.go
```

`scripts/check-examples.sh` previously enforced this; that
script has been removed. The pairing rule still applies — a `bad-`
without a matching `good-` is half a teaching asset. Do not skip
the good side.

### 4. Prompts have fallback contracts

Any prompt that depends on an external tool (CodeGraph, a
specific linter, a package manager) must include an explicit
fallback contract:

- The trigger condition (how to detect unavailability).
- The fallback order (which local means to use).
- The exact fallback line the agent must emit in the final
  report — verbatim, because downstream consumers parse it.
- The stop-loss rule (when to stop trying).

The "CodeGraph unavailable; context was gathered by rg/file
inspection." line is the canonical example. Inventing your own
phrasing here breaks downstream consumers that grep for the
canonical string.

### 5. Review output matches the contract

`examples/review-output.md` is not aspirational — it is the
contract. The format declared in `SKILL.md` (severity ladder,
no-finding sentence, required sections) must match the example
exactly. If a consumer loads only the example to calibrate, the
real output must look like the example.

### 6. Cross-references are explicit, not implicit

When a prompt depends on or is closely related to another prompt
in the same skill, the dependency must be visible at the top of
the file as a `> See also: prompts/x.md, prompts/y.md` blockquote.
Consumers should not have to read `SKILL.md` to re-discover the
Required/Optional Loading structure every time they load a
prompt.

### 7. Install is friction-free

A new consumer should be able to install a single skill with one
`npx skills add owner/repo --skill <name>` command. The root
`README.md` must list the install command for every skill in the
"Included Skills" section, not just the top three. Skills with
unsatisfied peer dependencies (e.g. `goal-driven-development`
needs the reviewers) must spell out the combined install command
in "Cross-Skill Dependencies".

### 8. Release is a single command

Releases should be reproducible from the current state of `main`.
Tag manually with `git tag -a vX.Y.Z -m "..."` and `git push origin
vX.Y.Z`. Do not hand-edit tags after they are pushed.
without the script unless the script is broken. If the script is
broken, fix the script first, then release.

## License

By contributing, you agree your contributions are licensed under the
repository's MIT license.
