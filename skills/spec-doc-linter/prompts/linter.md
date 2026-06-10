# Spec Doc Linter Core Protocol

> See also: prompts/drift-checks.md, prompts/semantic-review.md, prompts/confirmation.md

## Goal

Detect drift between two project-local doc files and the code they describe; propose a per-file patch; write only after per-file y/n/q confirmation.

The two doc conventions this skill teaches:

- `DevAgent.md` — per module. Sections (in order):
  `## Overview`, `## Public API`, `## Invariants`,
  `## Anti-patterns`, `## Dependencies`, `## File Map`.
- `CONTEXT.md` — per domain folder. Sections (in order):
  `## Bounded Context`, `## Owns`, `## Does Not Own`,
  `## Upstream Contexts`, `## Downstream Contexts`,
  `## Critical Invariants`.

A doc that is missing a required section is itself a drift finding (mode `STRUCT-MISSING-SECTION`).

## Discovery Order

1. Find candidate doc files:
   - `rg -l '^## (Overview|Public API|Invariants|Anti-patterns|Dependencies|File Map)$' --glob 'DevAgent.md'`
   - `rg -l '^## (Bounded Context|Owns|Does Not Own|Upstream Contexts|Downstream Contexts|Critical Invariants)$' --glob 'CONTEXT.md'`
2. For each `DevAgent.md`, identify the module root (parent directory) and the language/toolchain (look for `package.json`, `Cargo.toml`, `go.mod`, `pyproject.toml`, `pom.xml`, `build.gradle`, `*.csproj`).
3. For each `CONTEXT.md`, identify the domain folder (parent directory) and enumerate sibling/child code roots.
4. Build a target list of (doc file, scope root, toolchain).

## Drift Mode Catalog

### Class A — mechanical (Tier 1, see prompts/drift-checks.md)

| ID | Mode | Detection |
| --- | --- | --- |
| A.1 | Renamed symbol | `Foo.bar` reference in doc vs `rg` for `Foo` and `bar` |
| A.2 | Changed signature | parsed signature vs current AST/regex |
| A.3 | Module/domain added or removed | doc files vs file tree |
| A.4 | Quoted file path missing | `./path` in doc vs filesystem |
| A.5 | New dep not in doc | manifest diff vs doc `Dependencies` |
| A.6 | Broken intra-doc link | relative Markdown link resolution |

### Class B — judgment (Tier 2, see prompts/semantic-review.md)

| ID | Mode | Detection |
| --- | --- | --- |
| B.1 | Weakened invariant | doc claim vs code behaviour on re-call |
| B.2 | Anti-pattern newly violated | doc forbids X, code does X |
| B.3 | False idempotency/atomicity/thread-safety claim | doc claim vs code shape |
| B.4 | Boundary crossed | doc dependency rule vs new import |
| B.5 | Concurrency claim false | doc claim vs Send/Sync shape |
| B.6 | Context map out of date | doc context list vs actual imports |

## Report Format

Per file, in order encountered:

```text
## <relative/path>

- <file:line>  <ModeId>  <one-line evidence>  →  <suggested fix>
- ...
```

If a file has no findings, omit it.

## "All Clean" Sentinel

When no drift is found in any scanned doc, output exactly:

```text
未发现文档与代码漂移。
```

(Do not reuse the reviewer sentinel `未发现明确高风险问题。` — it has a different meaning and a different consumer set.)

## Confirmation Gate

After the report, do not write. Invoke `prompts/confirmation.md` for the per-file y/n/q loop.

## Fallback

If the toolchain probe (step 1.2) cannot identify a manifest, treat the file as `unknown-toolchain` and run only the language-agnostic checks: A.3, A.4, A.6, B.2, B.4, B.6. The final report must include the line:

```text
Toolchain unrecognised; only language-agnostic drift checks ran.
```
