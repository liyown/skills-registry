# Goal Driven Development Workflow

> See also: prompts/codegraph.md


## Goal

Turn an existing spec or goal into a verifiable code change, and leave reusable project knowledge when done. Do not write the spec; consume it as the source of truth.

## Phases

1. **Goal intake**
   - Read the user-supplied spec, goal, task description, or issue.
   - Extract success criteria, explicit out-of-scope items, affected modules, verifiable results.
   - If the goal lacks an actionable entry, first search the repo for the closest implementation point; if still uncertain, ask.
   - If the goal does not provide verifiable acceptance criteria, record the gap and ask before implementing. Acceptance-criteria template:

     ```text
     - Given <precondition>
       When <action>
       Then <observable result>
     - Out of scope: <explicit non-goals>
     - Affected modules: <files / packages / services>
     - Verifiable: <test name, build target, screenshot, API call>
     ```

2. **CodeGraph context**
   - Get structural context via `prompts/codegraph.md`: prefer MCP, fall back to CLI, then to `rg` / file reading.
   - Output internal working assumptions: entry point, key symbols, call chain, impact radius, test entry.
   - CodeGraph is not the sole source of business rule truth; business rules are still governed by spec and code behaviour.
   - When falling back, the final report must include `CodeGraph unavailable; context was gathered by rg/file inspection.` — never silently downgrade.

3. **Implementation**
   - Implement the minimal viable change to meet the goal.
   - Prefer existing architecture, tools, test frameworks, and naming conventions.
   - When adding a new abstraction, it must reduce real complexity or match an existing pattern.
   - Default change-size budget: ≤ ~200 lines of diff (excluding generated files, fixtures, and lockfiles). If the goal cannot be met within the budget, split the goal before implementing, or record the exception in the final report with the reason. The budget is a soft cap, not a gate; the goal's acceptance criteria are the gate.

4. **Verification**
   - Run the tests, type checks, builds, or page verifications closest to the change.
   - If tests cannot run, record the command, failure reason, and the alternative verification used.
   - Frontend visual changes need browser or screenshot verification; backend behaviour changes need unit / integration / API-level verification.

5. **Review gate**
   - Java/Spring/MyBatis/Redis/Kafka/Reactor changes invoke `java-code-reviewer`.
   - React/TypeScript/Next.js/Vite changes invoke `react-code-reviewer`.
   - Go/gRPC/sqlx changes invoke `go-code-reviewer`.
   - Python/asyncio/SQLAlchemy/Django/FastAPI changes invoke `python-code-reviewer`.
   - Node.js/Express/Fastify/Prisma/TypeORM changes invoke `node-code-reviewer`.
   - Critical/High must be fixed or explicitly marked `需要结合上下文确认` — never silently skipped.

6. **Knowledge capture**
   - Invoke `project-knowledge-capture`.
   - Capture only stable knowledge: entry point, decisions, constraints, tests, review conclusions.
   - Do not capture chat transcripts, failed attempts, secrets, customer data, or production data.

## Completion Criteria

- Implementation aligns with the spec.
- Code evidence supports every key implementation path.
- When CodeGraph is unavailable, the final report declares it explicitly.
- Related verification has run, or the inability to run has been recorded.
- Code review has been performed.
- Knowledge capture has been written, or the inapplicability is explained.

## Anti-example

- Read the spec, then change code without impact analysis.
- Tests fail but the change is claimed complete.
- Review gate surfaces Critical/High and is recorded but not addressed.
- Knowledge capture becomes a chat summary.

## Positive Example

```markdown
Goal: add a payment-status filter to the order list.
CodeGraph: locate `OrderController#list`, `OrderQueryService#list`, `OrderMapper#selectPage`.
Implementation: add a `status` parameter, propagate to the query object, and append a non-empty condition in the Mapper.
Verification: run `OrderQueryServiceTest` and the mapper SQL tests.
Review: java-code-reviewer / react-code-reviewer / go-code-reviewer / python-code-reviewer / node-code-reviewer reports no high-risk issues.
Capture: docs/knowledge/2026-06-05-order-status-filter.md.
```

Positive example for the fallback scenario:

```markdown
CodeGraph unavailable; context was gathered by rg/file inspection.
Location: rg -n 'OrderController|OrderQueryService' src/main/java followed by manual reading of `OrderController#list`, `OrderQueryService#list`, `OrderMapper#selectPage`.
```
