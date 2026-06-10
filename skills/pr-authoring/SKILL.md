---
name: pr-authoring
description: "Use when preparing a pull request for an enterprise repo: organise the PR description, attach the verification evidence, name the required reviewers per CODEOWNERS, surface linked issues / RFCs / design docs, and emit a 中文 + 企业 CI 兼容的 PR body. Triggers on: 写 PR, PR 描述, 整理 commit, 准备合入, code review 请求, 列出 reviewer, 附 verification 证据, link issue, link RFC, 关联设计文档. Do NOT use for: cross-team blast-radius analysis (use code-ownership-impact), or for the long-form RFC (use cross-team-rfc-draft). Complementary to code-ownership-impact, which is the per-PR packet in English + generic; pr-authoring is the 中文 + 企业 CI 实践 + 多 reviewer 列表 variant, emits a PR body shaped for the consumer's existing PR template, and links RFCs / design docs / issues in the format the org uses."
metadata:
  short-description: PR 准备: 描述、reviewer 名单、verification 证据、跨仓库引用
---

# PR Authoring

为企业的仓库准备一个 pull request：把 PR body 组织成 reviewer 一眼能看懂的格式，附上 verification 证据，列出 CODEOWNERS 解析后的 reviewer 名单，把关联的 issue / RFC / 设计文档链接到 org 用的格式上。

## Required Loading

Always load:

- `prompts/pr-body.md` — the 8-section PR body template; the per-section checklist; the 中文 + 企业 CI 兼容的 formatting rules.

## When To Run

- 当你本地有 1+ 个 commit 准备开 PR。
- 当 reviewer 在 PR 上留言让你重新组织 PR body。
- 当你接手别人的 PR 想让它符合 org 的 PR 模板。
- 当一个 PR 跨越多个仓库（monorepo 或 federated repos），需要把每个仓库的引用都列出来。

## Discovery Order

1. 看 `git log` 找出 PR 范围内的 commits。
2. 看 `git diff --stat main..HEAD` 找出改了哪些文件 / 多少行。
3. 解析 `CODEOWNERS` 找出每个改动的 path 对应的 team + named approver。
4. 找出关联的 issue / RFC / 设计文档（从 commit message、PR template、或 org 的 `gh issue list --search`）。
5. 收集 verification 证据（`make test`、`npm test`、CI 链接、manual smoke test 截图、metrics 截图）。
6. 按 org 的 PR 模板填入 8 个 section。

## Output Contract

PR body 是一个 Markdown 文档，8 个 section 按固定顺序：

1. **Summary** — 1-3 段。这个 PR 改什么、为什么改、影响范围。
2. **Type of change** — checkbox list（按 org 的 PR 模板）。
3. **Commits** — `git log --oneline main..HEAD` 的输出。
4. **Files changed** — `git diff --stat main..HEAD` 的输出 + 关键文件的链接。
5. **Reviewers** — 从 CODEOWNERS 解析出的 per-team reviewer 名单，附 `@` ping。
6. **Verification evidence** — unit / integration / e2e / manual / metrics，每类至少一项。
7. **Linked docs** — issue / RFC / 设计文档的链接，按 org 格式。
8. **Risk + rollback** — 这个 PR 的风险点 + 回滚路径。

每个 section 都要有内容。一个 PR body 缺任何一个 section 是 `PR-INCOMPLETE` finding。

## 中文 + 企业 CI 兼容的 formatting

- **Code blocks** 用 ` ``` ` 围起来；行内 code 用 `` ` ``。
- **路径** 用 `` `path/to/file.java:42` `` 的格式。
- **PR description 的 checkbox** 用 `- [ ]`，跟 GitHub 渲染兼容。
- **Reviewer 名单** 用 `@handle` 格式，GitHub 会自动解析。
- **Commit message** 用 conventional commits（`feat:`, `fix:`, `docs:` 等）。
- **不要用 emoji** 替代文字（CI 的 grep 会断）。

## Tools

- **`codegraph_impact`** — blast radius。
- **`codegraph_files`** — path → team map（与 `code-ownership-impact` 共用）。
- **`git log` / `git diff`** — 改动列表。
- **GitHub CLI** — `gh pr create --body-file` 直接从文件读 PR body。
- **Fallback** — 手动 `git log` + `git diff` + 手动解析 CODEOWNERS。

## Fallback

If CodeGraph is unavailable, the fallback is manual `git log` + manual CODEOWNERS parsing. The final report must include the line:

```text
CodeGraph unavailable; pr body gathered by rg/git inspection.
```

## Related Skills

- **`spec-doc-linter`** — keeps `DevAgent.md` / `CONTEXT.md` honest.
- **`code-ownership-impact`** — the per-PR packet in English + generic; this skill is the 中文 + 企业 CI variant.
- **`cross-team-rfc-draft`** — the long-form RFC; this skill is the per-PR counterpart.
- **`goal-driven-development`** — orchestrates the spec-to-code flow; this skill is the PR-prep counterpart.

## Examples

Each `bad-*` doc has a matching `good-*` in `examples/` showing the minimum PR body. Read them side by side to calibrate depth. `examples/pr-output.md` is the canonical "what the agent should emit" sample.
