# PR Body 模板

> See also: ../SKILL.md

PR body 是 8 个 section 的 Markdown 文档，按固定顺序。每节有 checklist，缺一节就是 `PR-INCOMPLETE` finding。

## The 8 sections

### 1. Summary

1-3 段。这个 PR 改什么、为什么改、影响范围。

- [ ] 改动的一句话总结。
- [ ] 为什么改（动机 / 关联的 issue / RFC）。
- [ ] 影响范围（哪些 service / 哪些 consumer / 哪些 team）。

### 2. Type of change

Checkbox list（按 org 的 PR 模板）。本 collection 的 8 个 reviewer 都不需要 checkbox，但 `pr-authoring` 生成的 PR body 会被其他 repo 的 PR 模板消费。

- [ ] New skill under `skills/<name>/`
- [ ] New scenario prompt
- [ ] New `bad-*` / `good-*` example pair
- [ ] Documentation
- [ ] CI / script
- [ ] Bug fix
- [ ] Refactor

### 3. Commits

`git log --oneline main..HEAD` 的输出 + 一句话 per commit。

### 4. Files changed

`git diff --stat main..HEAD` 的输出 + 关键文件的路径 + 行号链接。表格比纯文本清晰。

### 5. Reviewers

从 `CODEOWNERS` 解析出的 per-team reviewer 名单。每个 reviewer 一行：`@org/<team> — @<handle> — 关联的 path`。

- [ ] 至少一个 reviewer（per-team）。
- [ ] 跨 team 的 PR 每个 team 都有 reviewer。
- [ ] 共享 kernel 的改动列出所有 subscriber team 的 reviewer。

### 6. Verification evidence

每类至少一项：unit / integration / e2e / manual / metrics。

- [ ] Unit tests: `<test class>.<test method>` (passing)
- [ ] Integration tests: `<test class>.<test method>` (passing)
- [ ] e2e tests: `<test name>` (passing)
- [ ] Manual smoke: `<what you did> — <link to screenshot or recording>`
- [ ] CI: <link to the CI run>
- [ ] Metrics: <link to the dashboard or query>

### 7. Linked docs

关联的 issue / RFC / 设计文档的链接，按 org 格式。

- [ ] Issue: `#<id>` or full URL
- [ ] RFC: `docs/rfcs/...` or full URL
- [ ] Design doc: `<org's design doc URL>`
- [ ] 关联的 PR（if this PR is a follow-up）

### 8. Risk + rollback

这个 PR 的风险点 + 回滚路径。

- [ ] Risk: `<what could go wrong>`
- [ ] Rollback: `<how to revert>` (e.g. `git revert <commit>`, `POST /admin/flags/<name>/off`)
- [ ] 关联的 runbook: `<runbook path>` (if any)
- [ ] On-call: `<on-call schedule>` (if this PR can page)

## 中文 + 企业 CI 兼容的 formatting

- **Code blocks** 用 ` ``` ` 围起来；行内 code 用 `` ` ``。
- **路径** 用 `` `path/to/file.java:42` `` 格式。
- **Checkbox** 用 `- [ ]`（跟 GitHub 渲染兼容）。
- **Reviewer 名单** 用 `@handle` 格式，GitHub 会自动解析。
- **Commit message** 用 conventional commits。
- **不要用 emoji** 替代文字（CI 的 grep 会断）。

## The "all clean" sentinel

The skill is not a linter. A PR body where every section has content and every reviewer is named is the "no finding" case. There is no canonical "all clean" line; the PR body itself is the output.
