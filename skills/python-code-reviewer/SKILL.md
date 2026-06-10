---
name: python-code-reviewer
description: Review Python backend services, pull requests, diffs, and snippets for production-risk bugs in asyncio, async/await, GIL, exception handling, type hints, SQLAlchemy/Django ORM, FastAPI/aiohttp, dependencies, security, and microservice architecture — focusing on capital-loss, data-corruption, authorization-bypass, and incident-causing defects. Use for asyncio event-loop blocking, requests-inside-async, pickle.loads, SQLAlchemy N+1, and incident-blocking PR review. Do NOT use for: PEP 8 style nits, type-hint pedantry, or non-incident review.
metadata:
  short-description: Evidence-driven Python backend review
---

# Python Code Reviewer

Review Python backend code for real production risk. Keep the entrypoint small and load scenario prompts only when the code needs them.

## Required Loading

Always load `prompts/reviewer.md`.

Load additional prompts only when relevant:

- `prompts/async-reviewer.md`: asyncio, event loop, blocking calls, task groups, cancellation, GIL, shared state.
- `prompts/error-reviewer.md`: exception types, raise from, swallowing exceptions, logging, context managers, `__exit__` exceptions.
- `prompts/sql-reviewer.md`: SQLAlchemy, Django ORM, async sessions, N+1, transactions, connection pool.
- `prompts/web-reviewer.md`: FastAPI, aiohttp, Django, middleware, timeouts, request lifecycle, streaming.
- `prompts/security-reviewer.md`: authorization, tenant isolation, injection, deserialization, secrets, SSRF.

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
