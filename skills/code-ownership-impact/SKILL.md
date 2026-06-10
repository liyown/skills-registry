---
name: code-ownership-impact
description: "Use before any change that crosses team or service boundaries, to find CODEOWNERS, blast radius, and required reviewers. Triggers on: rename, refactor, deprecate, schema change, API break, cross-team change, 'who owns this', 'who needs to review', 'what does this PR break', CODEOWNERS, service boundary, dependency request, shared kernel. Do NOT use for: code-review findings (use java-code-reviewer, react-code-reviewer, go-code-reviewer, python-code-reviewer, or node-code-reviewer for that), or for plan decomposition within one team (use obra/superpowers writing-plans). Complementary to goal-driven-development, which orchestrates the spec-to-code flow; code-ownership-impact is the pre-merge counterpart that names the per-team approvers."
metadata:
  short-description: Cross-team blast radius + CODEOWNERS for shared-kernel changes
---

# Code Ownership Impact

Before any change that crosses team or service boundaries, find CODEOWNERS, blast radius, and required reviewers. The output is a multi-team review packet with named approvers per affected path.

## Required Loading

Always load:

- `prompts/ownership-walk.md` — the per-team-ping-list and shared-kernel-detection framework.

## When To Run

- Before merging a rename, refactor, deprecate, schema change, or API break.
- Before opening a PR that touches a shared kernel.
- When the author needs to know "who else needs to review this?".
- When a dependency request is being prepared for another team.

## Discovery Order

1. Walk the diff (`git diff --name-only main..HEAD`).
2. For each changed path, resolve the team that owns the path. Use `codegraph_files` (path → team map) or `rg 'path:' CODEOWNERS` for a manual lookup.
3. For each team, list the named approvers from the `CODEOWNERS` file.
4. Detect shared-kernel membership: any file under `common/`, `shared/`, `pkg/`, `platform/`, or matching a shared-kernel pattern (e.g. `*::common::*`).
5. Compute the blast radius with `codegraph_impact` (or `rg` for the symbol + per-caller `Read`).
6. For each shared-kernel hit, list the subscriber services and their owning teams.
7. Suggest a review order: smaller teams / least-loaded first, larger / release-train-locked last.

## Output Contract

The review packet is grouped by team. For each team, name:

- the path(s) the team owns within the diff
- the named approvers (from `CODEOWNERS`)
- the suggested review deadline (next release-train slot)
- the cross-team impact (for the team's owned services that consume a shared-kernel change)
- the suggested ack-channel (a thread in the team's chat, a CODEOWNERS `@`-ping, or an RFC subscribe)

A diff that touches no shared kernel and no per-team boundary produces a single-team review packet (one team, the author's team) — that is itself the correct output.

## Tools

- **`codegraph_impact`** — primary tool for the blast-radius walk.
- **`codegraph_files`** — path → team map.
- **`codegraph_callers` / `codegraph_callees`** — boundary crossing (who calls into / who is called from the changed symbol).
- **`codegraph_search`** — locate a shared-kernel file by name.
- **Fallback** — `rg 'path:' CODEOWNERS` + `git log -- <path>` for per-path ownership inference. The fallback is acceptable but not as precise.

## Fallback

If CodeGraph is unavailable, the fallback is `rg CODEOWNERS` + `git log` per changed path. The final report must include the line:

```text
CodeGraph unavailable; ownership gathered by rg/file inspection.
```

## Related Skills

- **`spec-doc-linter`** — keeps `DevAgent.md` / `CONTEXT.md` honest; this skill is the human-side counterpart for cross-team change.
- **`cross-team-rfc-draft`** (planned v0.6.1) — the long-form RFC for changes that need a written proposal. code-ownership-impact is the per-PR packet; cross-team-rfc-draft is the proposal.
- **`goal-driven-development`** — orchestrates the spec-to-code flow; code-ownership-impact is the pre-merge counterpart that names approvers.

## Examples

Each `bad-*` doc has a matching `good-*` in `examples/` showing the minimum packet. Read them side by side to calibrate packet depth. `examples/ownership-output.md` is the canonical "what the agent should emit" sample.
