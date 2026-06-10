# PR Authoring Skill

`pr-authoring` 为企业仓库准备一个 pull request 的 body：8 个 section 固定顺序，每节有 checklist，附 reviewer 名单、verification 证据、跨仓库引用、风险 + 回滚路径。

## What It Drafts

- PR body 的 8 个 section。
- 从 `CODEOWNERS` 解析出的 per-team reviewer 名单。
- Verification 证据（unit / integration / e2e / manual / metrics）。
- 关联的 issue / RFC / 设计文档链接。
- 风险 + 回滚路径。

## Discovery Tools

- **`codegraph_impact`** — blast radius。
- **`codegraph_files`** — path → team map。
- **`git log` / `git diff`** — 改动列表。
- **GitHub CLI** — `gh pr create --body-file`。
- **Fallback** — manual `git log` + 手动 `CODEOWNERS` 解析。

## Output Contract

PR body 是 8-section Markdown 文档，固定顺序：Summary / Type of change / Commits / Files changed / Reviewers / Verification evidence / Linked docs / Risk + rollback。

## Related Skills

- **`spec-doc-linter`** — keeps `DevAgent.md` / `CONTEXT.md` honest。
- **`code-ownership-impact`** — 英文 + 通用版；这是 中文 + 企业 CI 变种。
- **`cross-team-rfc-draft`** — 长篇 RFC；这是 PR-prep counterpart。
- **`goal-driven-development`** — 编排 spec → code 流程；这是 PR 准备 counterpart。

## Files

```text
.
├── SKILL.md
├── README.md
├── prompts/
│   └── pr-body.md       # 8-section PR body 模板
└── examples/
    ├── bad-pr.md
    ├── good-pr.md
    └── pr-output.md     # canonical "what the agent should emit" sample
```
