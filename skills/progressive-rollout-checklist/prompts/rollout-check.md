# Rollout Check Framework

> See also: ../SKILL.md

The framework is a per-flag rollout packet. For each flag, the framework says what to record.

## The 8 cells per flag

1. **Flag name** — the literal flag identifier.
2. **Read sites** — every path + line that reads the flag. Use `codegraph_callers` or `rg '<flag-name>\s*[!=]=\|is[A-Z]Enabled'`.
3. **Write sites** — every path + line that flips the flag (e.g. the admin endpoint, the cohort expander, the kill switch). Use `codegraph_callees` or `rg '<flag-name>\s*=\s*true\|setFlag'`.
4. **Default value + proposed new default** — the current default the flag ships with, and the value the change proposes to flip to. If the change is a fresh flag, the current default is the value the new code is gated on.
5. **Cohort definition** — the rule that determines who sees the new behaviour. Examples: `tenant_id IN (...)`, `user_id % 100 < 10`, `region == 'us-east-1'`.
6. **SLO gate** — the metric + threshold that gates the cohort expansion. Examples: `error_rate < 0.1%`, `p99 < 800ms`. The SLO gate is the contract that says "if the metric is good, expand; if not, freeze".
7. **Kill switch** — the action that disables the flag globally. Examples: `POST /admin/flags/<name>/off`, `kubectl set env FLAG_NAME=false`. The kill switch has an owner (the on-call or a named service).
8. **Rollback hook** — the action that reverts the change. Distinct from the kill switch: the kill switch turns the flag off; the rollback hook reverts the change's effect (e.g. a migration reversal, a feature revert). The hook has an owner.

A sunset date is added as a 9th cell when the flag is being deprecated. A flag with no sunset date and no deprecation intent is marked "no sunset".

## The 3 rollout-packet sections

The final packet is structured as three sections, in this order:

1. **Flag rows** — every flag in scope, with the 8 (or 9) cells above. A flag missing a cell is highlighted.
2. **Cross-flag dependencies** — if flag A's cohort expansion depends on flag B's SLO gate, name the dependency and the order.
3. **Ship-readiness verdict** — for each flag, `READY` or `BLOCKED: <cell>`. The packet refuses to mark a flag-bearing change as "ready to ship" if any flag is `BLOCKED`.

## The "all clean" sentinel

The skill is not a linter. A packet where every flag row is `READY` is the "no finding" case; the verdict is itself the output. There is no canonical "all clean" line; the framework's output is always a packet.
