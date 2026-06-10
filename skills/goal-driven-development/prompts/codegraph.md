# CodeGraph Usage Prompt

> See also: prompts/workflow.md

CodeGraph defaults to `@colbymchenry/codegraph`. It exposes a code knowledge graph through MCP/CLI, suitable for "where is it implemented", "who calls whom", "what does this change affect", and "what is the call path".

The currently-exposed MCP tool list (verify with `codegraph --help` after install):

| Tool | Purpose |
| --- | --- |
| `codegraph_explore` | **Primary tool** — answers "how does X work" / flow questions in one call. Absorbs what older docs called `codegraph_context` and `codegraph_trace`. |
| `codegraph_search` | Find symbols by name across the codebase. |
| `codegraph_callers` | Upstream callers of a function / method. |
| `codegraph_callees` | Downstream callees. |
| `codegraph_impact` | Affected code radius for a symbol change. |
| `codegraph_node` | One symbol's details + full source (returns all overloads for ambiguous names). |
| `codegraph_files` | Indexed file structure (faster than `fs` scan). |
| `codegraph_status` | Index health and statistics. |

> **Note on prior names.** The earlier 6-tool surface (`codegraph_context`, `codegraph_trace`, etc.) is no longer exposed by the MCP server. Treat `codegraph_explore` as the primary "give me context for this goal / trace this path" tool; reserve the named tools above for narrower queries.

## Priority

### 1. CodeGraph MCP (preferred)

First call `codegraph_status` to confirm the index is ready, then:

- `codegraph_explore` — primary walk; one call answers "how does X work" end-to-end and returns a relationship map.
- `codegraph_search` — locate by symbol name.
- `codegraph_callers` / `codegraph_callees` — upstream and downstream.
- `codegraph_impact` — change impact radius.
- `codegraph_node` — full source for one symbol.
- `codegraph_files` — see the structure of indexed files.

### 2. CodeGraph CLI (fallback)

When MCP is unavailable but the CLI works, use:

- `codegraph status` — index health.
- `codegraph query <symbol>` — symbol search.
- `codegraph callers <symbol>` / `codegraph callees <symbol>` — upstream / downstream.
- `codegraph impact <symbol>` — change radius.
- `codegraph affected` — files affected by a planned set of changes.
- `codegraph files` — file tree.
- If the repo is not initialised, suggest `codegraph init -i`.

### 3. Fallback (no CodeGraph)

Only when `codegraph_status` reports unavailability (not installed, not indexed, version mismatch, cannot reach the service) do you fall back to local means. **Fallback has a cost** — a question CodeGraph answers in one query often needs 2-5 manual `rg` / file reads in fallback mode, so only ask questions that actually affect the implementation path.

Fallback order:

1. `rg` (ripgrep) for precise symbol / string / file pattern search.
2. Language service / IDE reference (`tsc --noEmit`, `go doc`, `javap`, JSDoc / PyDoc comments).
3. Infer behaviour from test naming conventions (`*_test.go`, `*.spec.ts`, `OrderServiceTest`).
4. Source-code read of the entry path (1-2 levels of calls from the entry).

**Mandatory declaration**: regardless of the fallback path, after the CodeGraph context phase you must include the following line in the final report:

```text
CodeGraph unavailable; context was gathered by rg/file inspection.
```

This declaration is the only signal a reviewer or reviewer-downstream uses to assess "is there graph-level evidence". **Do not silently downgrade** — even if the gathered information looks sufficient, declare.

**Stop-loss**: if 3 consecutive `rg` queries still cannot find a critical symbol, stop exploring and ask the user about the repo structure or request the user to add context, rather than running 5 more `rg` and overwhelming the context.

### 4. When to Stop Calling CodeGraph

Stop CodeGraph calls and switch to the implementation phase when any of the following holds:

- You have the five elements: entry point, key symbols, call chain, impact radius, test entry.
- The same target returned the same conclusion in 2 consecutive queries.
- Any `codegraph_status` in the past 5 minutes already reported unavailable — go straight to fallback, do not probe again.

## Usage Principles

- Prefer CodeGraph for architecture, call chain, and impact questions over blind `rg`.
- CodeGraph output is navigation evidence, not proof of business rules — for critical risks, go back to the source code or tests to confirm.
- Do not call CodeGraph redundantly for its own sake; stop exploring once you have enough context.
- Per session, at most 2 CodeGraph queries for the same target — beyond that is waste.

## Typical Queries

- Find entry: `codegraph_explore` with goal summary.
- Find callers: `codegraph_callers` on service / handler / function symbol.
- Find impact: `codegraph_impact` on changed symbol.
- Find path: `codegraph_explore` from controller / route to persistence / external call.
- Find test entry: `codegraph_files` filtered with `*test*` + `codegraph_explore` with query "tests for <symbol>".
