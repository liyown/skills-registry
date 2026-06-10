---
name: go-code-reviewer
description: Review Go backend services, pull requests, diffs, and snippets for production-risk bugs in goroutines, channels, context propagation, error handling, defer/panic, sqlx/GORM, gRPC/HTTP servers, generics, memory, and microservice architecture — focusing on capital-loss, data-corruption, authorization-bypass, and incident-causing defects. Use for goroutine lifetime, errgroup.Wait missing, context cancellation, errors.Is/%w chains, and incident-blocking PR review. Do NOT use for: gofmt / golangci-lint style nits, naming, or non-incident review.
metadata:
  short-description: Evidence-driven Go backend review
---

# Go Code Reviewer

Review Go backend code for real production risk. Keep the entrypoint small and load scenario prompts only when the code needs them.

## Required Loading

Always load `prompts/reviewer.md`.

Load additional prompts only when relevant:

- `prompts/concurrency-reviewer.md`: goroutines, channels, WaitGroup, errgroup, race conditions, leaks, sync primitives.
- `prompts/context-reviewer.md`: context propagation, cancellation, deadlines, request-scoped values, background vs TODO.
- `prompts/error-reviewer.md`: error wrapping, sentinels, %v vs %w, panic recovery, defer error handling.
- `prompts/sql-reviewer.md`: sqlx, GORM, database/sql, prepared statements, transactions, N+1, connection pool.
- `prompts/rpc-reviewer.md`: gRPC, HTTP servers, middleware, timeouts, retries, status codes, streaming.
- `prompts/security-reviewer.md`: authorization, tenant isolation, injection, sensitive logs, file/URL/SSRF, secrets.

## Review Contract

- Find concrete bugs, not generic advice.
- Bind every finding to code evidence and an execution path.
- Mark uncertain findings as `需要结合上下文确认`.
- Do not output style comments unless they hide a real defect.
- If no clear high-risk issue is found, output exactly:

```text
未发现明确高风险问题。
```

Use the severity and output contract from `prompts/reviewer.md`.

## Examples

Each bad example has a matching `good-<file>` in this same `examples/`
directory that shows the minimal fix for every Critical/High finding. Read
both side by side when triaging a real diff.
