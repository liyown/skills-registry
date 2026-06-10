---
name: java-code-reviewer
description: Review Java backend projects, pull requests, diffs, and snippets for production-risk bugs in Spring Boot, Spring Cloud, MyBatis, MyBatis-Plus, Redis, Kafka, Reactor, WebFlux, transactions, concurrency, security, performance, DDD, and microservice architecture — focusing on capital-loss, data-corruption, authorization-bypass, and incident-causing defects. Use for transaction integrity, Spring proxy correctness, MyBatis wrapper safety, Redis/Kafka duplicate consumption, and incident-blocking PR review. Do NOT use for: code style, formatting, naming, lint, or other non-incident review.
metadata:
  short-description: Evidence-driven Java backend review
---

# Java Code Reviewer

Review Java backend code for real production risk. Keep the entrypoint small and load scenario prompts only when the code needs them.

## Required Loading

Always load `prompts/reviewer.md`.

Load additional prompts only when relevant:

- `prompts/spring-reviewer.md`: Spring Boot, Spring Cloud, transactions, proxy annotations, validation, async, scheduling.
- `prompts/mybatis-reviewer.md`: MyBatis, MyBatis-Plus, SQL, Wrapper, batch operations, indexes, tenant filters.
- `prompts/security-reviewer.md`: authorization, tenant isolation, injection, sensitive data, file/URL/deserialization risks.
- `prompts/concurrency-reviewer.md`: idempotency, duplicate writes, locks, thread pools, ThreadLocal, cache races.
- `prompts/reactor-reviewer.md`: Reactor/WebFlux blocking, subscribe, timeout, retry, backpressure, EventLoop misuse.
- `prompts/redis-kafka-reviewer.md`: Redis cache failure modes, Kafka consumption, ordering, offset, retries, compensation.

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
