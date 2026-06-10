# Goal Driven Development Skill

`goal-driven-development` is a workflow skill for implementing existing goals/specs. It does not write the spec.

## Flow

1. Read goal/spec.
2. Use CodeGraph for context and impact analysis. If unavailable, fall back to `rg` and source reading — declare the fallback in the final report using the exact line `CodeGraph unavailable; context was gathered by rg/file inspection.`
3. Implement scoped changes.
4. Verify with tests/builds/checks.
5. Run Java, React, Go, Python, or Node review gates.
6. Capture durable knowledge.

## Dependencies

- `java-code-reviewer`
- `react-code-reviewer`
- `go-code-reviewer`
- `python-code-reviewer`
- `node-code-reviewer`
- `project-knowledge-capture`

## Install Path

When a consumer runs `npx skills add liyown/skills-registry --skill goal-driven-development`, the entire `skills/goal-driven-development/` directory is copied to their local skills directory. The default location is:

```text
~/.claude/skills/goal-driven-development/   # Claude Code
~/.cursor/skills/goal-driven-development/    # Cursor (if applicable)
```

`SKILL.md` is the entrypoint the consumer's agent reads first. The agent then loads `prompts/workflow.md` and `prompts/codegraph.md` as the two always-loaded prompts. Scenario-specific reviewer skills (`java-code-reviewer`, `react-code-reviewer`, `go-code-reviewer`, `python-code-reviewer`, `node-code-reviewer`) and the `project-knowledge-capture` skill are **referenced by name** in `SKILL.md` and are NOT auto-installed. The consumer must install them separately:

```sh
npx skills add liyown/skills-registry \
  --skill goal-driven-development \
  --skill java-code-reviewer \
  --skill react-code-reviewer \
  --skill go-code-reviewer \
  --skill python-code-reviewer \
  --skill node-code-reviewer \
  --skill project-knowledge-capture
```

If the agent tries to invoke `java-code-reviewer` and that skill is not installed locally, the invocation will silently no-op or error in agent-specific ways. The combined install command above is the supported way to use `goal-driven-development`.
