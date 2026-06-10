# Node Code Reviewer Skill

`node-code-reviewer` is an evidence-driven Node.js backend review skill. It is tuned for production-risk findings, not broad style critique.

## What It Reviews

- Event-loop blocking by sync I/O (readFileSync, JSON.parse on large payloads, sync crypto)
- Missing await / unhandledRejection / fire-and-forget Promises
- Prototype pollution via `Object.assign` / spread merge of user input
- Prisma / TypeORM / Sequelize / Knex raw SQL and `$executeRawUnsafe`
- Express / Fastify / Koa / Hono middleware order, body limits, double-response

## Prompt Loading

Always load `prompts/reviewer.md`. Load scenario prompts only when code evidence requires them.

## Examples

Each bad example has a matching `good-<file>` in `examples/` that shows the minimal fix for every Critical/High finding. Read both side by side when triaging a real diff.
