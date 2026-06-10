---
name: node-code-reviewer
description: Review Node.js backend services, pull requests, diffs, and snippets for production-risk bugs in async/await, event loop, error handling, streams, Prisma/TypeORM/Sequelize/Knex, Fastify/Express/Koa/Hono, timeouts, concurrency, security, and microservice architecture — focusing on capital-loss, data-corruption, authorization-bypass, and incident-causing defects. Use for event-loop blocking by readFileSync, missing await/unhandledRejection, prototype pollution, Prisma $executeRawUnsafe, and incident-blocking PR review. Do NOT use for: ESLint/Prettier style nits, naming, or non-incident review.
metadata:
  short-description: Evidence-driven Node backend review
---

# Node Code Reviewer

Review Node.js backend code for real production risk. Keep the entrypoint small and load scenario prompts only when the code needs them.

## Required Loading

Always load `prompts/reviewer.md`.

Load additional prompts only when relevant:

- `prompts/async-reviewer.md`: async/await, event loop, blocking calls, unhandled rejection, AbortController, worker threads, unref.
- `prompts/error-reviewer.md`: error types, Promise rejection, async stack traces, error middleware order, structured logging.
- `prompts/sql-reviewer.md`: Prisma, TypeORM, Sequelize, Knex, raw SQL, transactions, N+1, connection pool.
- `prompts/http-reviewer.md`: Fastify, Express, Koa, Hono, middleware order, timeouts, body limits, request lifecycle, streaming.
- `prompts/security-reviewer.md`: authorization, tenant isolation, injection, prototype pollution, deserialization, secrets, SSRF.

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
