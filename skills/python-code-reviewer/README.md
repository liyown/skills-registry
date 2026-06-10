# Python Code Reviewer Skill

`python-code-reviewer` is an evidence-driven Python backend review skill. It is tuned for production-risk findings, not broad style critique.

## What It Reviews

- asyncio / event-loop blocking, async/await misuse, exception swallowing
- SQLAlchemy / Django ORM N+1, transaction boundaries, async session lifecycle
- Pickle / unsafe deserialization, dependency confusion, request smuggling
- Concurrency primitives, GIL, task groups, cancellation, shared state

## Prompt Loading

Always load `prompts/reviewer.md`. Load scenario prompts only when code evidence requires them.

## Examples

Each bad example has a matching `good-<file>` in `examples/` that shows the minimal fix for every Critical/High finding. Read both side by side when triaging a real diff.
