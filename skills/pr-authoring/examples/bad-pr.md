# PR Body — Bad Draft

## The PR body (as the author would otherwise send it)

> 修了几个 bug，review。

## What's wrong with this PR body

- **Missing 7 of 8 sections.** No Summary, no Type of change, no Commits, no Files changed, no Reviewers, no Verification evidence, no Linked docs.
- **No reviewer 名单.** "Review" is one broadcast; the per-team reviewer should be named with `@handle`.
- **No verification 证据.** "修了几个 bug" is a claim; the reviewer needs unit / integration / e2e / manual / metrics 证据。
- **No linked docs.** 关联的 issue / RFC / 设计文档都缺。
- **No risk + rollback.** 改动可能引入 regression，没有风险点和回滚路径。
