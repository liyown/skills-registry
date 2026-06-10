# Tier-1 Drift Checks

> See also: prompts/linter.md, prompts/semantic-review.md

## How To Extend

Each check is a self-contained block with a header `### A.N`, a `Detection` recipe, a `False positive` guard, and an `Emits` row. To add a check, append a block. The Tier-1 runner iterates A.1 → A.N in order; new checks at the bottom require no other edits.

## A.1 Renamed symbol

- Detection: extract `Identifier.identifier` and `Identifier.method` patterns from the doc (regex `\b[A-Z][A-Za-z0-9_]*\.[a-z_][A-Za-z0-9_]*\b` and bare `\b[A-Z][A-Za-z0-9_]*\b`). For each, run `rg -l '<Identifier>' <module-root>` and `rg -l '<method>' <module-root>`. If neither matches, emit drift.
- False positive guard: identifiers mentioned in the `## Anti-patterns` section are intentionally cited as things to *avoid*; skip them.
- Emits: `A.1`, file:line, evidence "doc references `Foo.bar`; not found under `<module-root>`", suggested fix "remove the bullet, or update to the current symbol name."

## A.2 Changed signature

- Detection: for each line in `## Public API`, parse a function/method header (language-specific regex catalogue: TS `^\s*(export\s+)?(async\s+)?function\s+\w+\s*\(`, Rust `^\s*(pub\s+)?fn\s+\w+\s*\(`, Go `^\s*func\s+\([^)]*\)\s*\w+\s*\(`, Python `^\s*def\s+\w+\s*\(`, Java `^\s*(public|private|protected)?\s*[\w<>,\s]+\s+\w+\s*\(`). Run `rg` for the same symbol and extract its current signature. Compare normalised (`s/\s+//g`) parameter list and return type. Diff → emit drift.
- False positive: language-agnostic mode treats `def` vs `pub fn` mismatches as not applicable; skip.
- Emits: `A.2`, file:line, evidence "doc shows `fn (x: u32) -> Result<T, E>`; code shows `fn (x: u32, y: u32) -> Result<T, E>`", suggested fix "update the signature line to match the new parameter list."

## A.3 Module/domain added or removed

- Detection: walk the file tree from the repo root (or the consumer's configured scope root). For each directory that contains source files (has a manifest OR ≥ 3 source files in a non-`node_modules`/non-`target`/non-`.git` path), check that a `DevAgent.md` exists one level up. For each `CONTEXT.md`, check that the parent domain folder still exists. Conversely, for each `DevAgent.md`, check that the module root contains source files.
- Emits (doc missing): `A.3a`, doc-path, evidence "module `<path>` has source but no DevAgent.md", fix "create DevAgent.md with the canonical sections."
- Emits (orphan doc): `A.3b`, doc-path, evidence "DevAgent.md exists but module is empty or removed", fix "delete or move the doc."

## A.4 Quoted file path missing

- Detection: extract backticked relative paths from the doc (regex `` `[./][^`]+` `` followed by a file extension). For each, resolve relative to the doc's directory and check the filesystem.
- False positive: URLs, anchor links, and paths with no extension are skipped.
- Emits: `A.4`, file:line, evidence "doc references `./src/foo/bar.rs`; file does not exist", fix "remove the reference or create the file."

## A.5 New dependency not in doc

- Detection: parse the toolchain manifest (`package.json` `dependencies`/`devDependencies`, `Cargo.toml` `[dependencies]`, `go.mod` `require`, `pyproject.toml` `[project.dependencies]`, `pom.xml` `<dependency>`, `build.gradle` `dependencies`). Build a name list. Extract the `## Dependencies` section of the doc and parse its bullets. Emit drift for any manifest entry not mentioned in the doc.
- False positive: build/test/dev-tool dependencies that the consumer has tagged `// dev-only` in the manifest are skipped.
- Emits: `A.5`, file:line, evidence "manifest lists `lodash-es`; doc Dependencies section does not mention it", fix "add a bullet to Dependencies describing what the module uses lodash-es for."

## A.6 Broken intra-doc link

- Detection: parse Markdown links `\[([^\]]+)\]\(([^)]+)\)`. For each target whose scheme is empty or `file`, resolve relative to the doc's directory and check existence.
- False positive: anchor-only links (`#section`) and http(s) URLs are skipped.
- Emits: `A.6`, file:line, evidence "link `[order](./orders.md)` does not resolve", fix "fix the path or create the file."

## Running Order

The runner is `A.1 → A.2 → A.3 → A.4 → A.5 → A.6`. Each check is independent; the order exists only so that symbol-level findings surface before file-level ones.
