# Confirmation Flow

> See also: prompts/linter.md

The linter never writes without explicit per-file user confirmation. There is no `--fix-all` flag.

## Loop

For each doc file with at least one Class A or Class B finding (excluding `NEEDS_CONTEXT`):

1. Stage the proposed edit as a unified diff in memory. Do not write the file.
2. Print the diff to the terminal in this exact shape:

   ```text
   --- <relative/path>
   +++ <relative/path>
   @@ -<old-start>,<old-count> +<new-start>,<new-count> @@
   < -old line>
   < -old line>
   > +new line
   > +new line
    context line
   ```

3. Print a one-line summary:

   ```text
   <N> change(s) proposed in <relative/path> (M Class A, K Class B)
   Apply? [y/n/q]:
   ```

4. Read the user's reply:
   - `y` (or `yes`): apply the patch to the file, print `wrote <relative/path>`, move to the next file.
   - `n` (or `no`): print `skipped <relative/path>`, move to the next file.
   - `q` (or `quit`): print `stopped at <relative/path>`, end the run.
   - Anything else: re-print the prompt unchanged.

## Edge Cases

- A file with only `NEEDS_CONTEXT` findings is **not** editable. Print it under a separate `## Needs Context` section in the report and skip.
- If the proposed diff is empty (all findings are `NEEDS_CONTEXT` or were resolved mid-loop), do not prompt; print `no changes for <relative/path>` and move on.
- If the user types `q` mid-loop, the in-flight patch is discarded; nothing is written.
- A failed write (permission, IO) must surface the error and abort the run; the agent must not silently swallow IO failures.
