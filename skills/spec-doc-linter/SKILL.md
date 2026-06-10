---
name: spec-doc-linter
description: Lint and auto-sync two project-local doc files — per-module DevAgent.md and per-domain CONTEXT.md — against the code they describe. Detects mechanical drift (renamed symbols, removed files, broken links, missing dependencies) and LLM-judgment drift (weakened invariants, anti-patterns newly violated, false concurrency/idempotency claims, broken domain boundaries, stale context maps), proposes a per-file patch, and writes only after a per-file y/n/q confirmation. Use when the user says "sync docs with code", "lint the spec", "refresh DevAgent.md", "update CONTEXT.md", "spec is stale", or "doc/code drift". Do NOT use for: writing new specs from scratch, formatting/lint of code style, or reviewer-style production-risk review (use java-code-reviewer, react-code-reviewer, go-code-reviewer, python-code-reviewer, or node-code-reviewer for that). Complementary to goal-driven-development, which captures decisions into docs/knowledge/ — spec-doc-linter keeps DevAgent.md and CONTEXT.md honest after the code has moved.
metadata:
  short-description: Doc/code drift linter with auto-sync
---

# Spec Doc Linter

Sync two project-local doc conventions with the code they describe. Detect drift, propose a per-file patch, write only after per-file y/n/q confirmation.

## Required Loading

Always load:

- `prompts/linter.md` — core protocol: doc conventions, drift catalog, sentinel.
- `prompts/drift-checks.md` — Tier-1 static rules (extendable; one block per A.N).
- `prompts/semantic-review.md` — Tier-2 LLM-judgment rubric for Class B modes.
- `prompts/confirmation.md` — per-file diff display + y/n/q contract.

## When To Run

- After a non-trivial merge, before opening a PR.
- When the user says the spec or CONTEXT feels stale.
- When a module's DevAgent.md or a domain's CONTEXT.md was last touched months ago and the module's code has churned since.

## Two Doc Conventions This Skill Teaches

- `DevAgent.md` lives at a module's root (e.g. `src/main/java/com/example/order/DevAgent.md`). Sections (in order): `## Overview`, `## Public API`, `## Invariants`, `## Anti-patterns`, `## Dependencies`, `## File Map`.
- `CONTEXT.md` lives at a domain folder's root (e.g. `domains/billing/CONTEXT.md`). Sections (in order): `## Bounded Context`, `## Owns`, `## Does Not Own`, `## Upstream Contexts`, `## Downstream Contexts`, `## Critical Invariants`.

A doc that is missing a required section is itself a drift finding (mode `STRUCT-MISSING-SECTION`).

`prompts/linter.md` defines the canonical section lists and parsing rules.

## Output Contract

- If no drift is found, output exactly the sentinel:
  `未发现文档与代码漂移。`
- When drift is found, output a drift report grouped by file; each finding carries `file:line`, drift-mode id, evidence, suggested fix. See the canonical sample in `examples/linter-output.md`.
- Mutations are gated by `prompts/confirmation.md`. No `--fix-all` flag: every file is a separate y/n/q prompt.

## Examples

Each `bad-*` doc has a matching `good-*` in `examples/` showing the minimum sync. Read them side by side to calibrate drift severity. The canonical "what the agent should emit" sample ships as `examples/linter-output.md`.
