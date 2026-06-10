# Call Trace Framework

> See also: ../SKILL.md

The framework is a per-hop trace with customer-impact framing and per-team paging list. For each hop, the framework says what to record.

## The 5 cells per hop

1. **Call site** — the path + line that makes the call. The call's argument is named if it carries the failing identifier.
2. **Owning team** — the team that owns this hop. Sourced from the path → team map.
3. **Latency budget vs. observed** — the SLO for this hop vs. the actual latency at the time of the incident. If logs are unavailable, mark "no observed data; using budget as proxy".
4. **Failure mode** — what went wrong: timeout, error rate, saturation, 5xx, partial failure, slow. The failure mode is the *kind* of degradation, not the cause.
5. **Downstream chain** — the next hop the call goes to. The chain is the link to the next row in the trace.

## The 3 trace sections

The final report is structured as three sections, in this order:

1. **Primary incident locus** — the first hop that degraded against its budget. One sentence. This is the row the on-call rotation focuses on first.
2. **Per-hop trace** — every hop from the entry point to the persistence / external system. For each hop, the 5 cells above. The first row that matches the primary locus is highlighted.
3. **Per-team paging list** — every team that owns a hop in the trace, with the on-call paging channel (PagerDuty schedule, Slack channel, or named on-call).
4. **Customer-impact framing** — how many customers / requests / regions are affected. Quantify if possible; if not, name the customers / request classes that are known to be affected.

## The SLO-breach case

When the incident is an SLO breach without a specific error, the trace is a **latency trace** rather than an error trace. The framework still applies; the "failure mode" cell is "p99 above SLO" rather than "5xx".

## The postmortem mode

When the skill runs in postmortem mode, the primary incident locus is known from the incident timeline. The trace is reconstructed from logs and the repo, not from a live run. The framework is the same; the "observed" cell is filled from logs (or marked "no observed data" if the logs are not available).

## The "all clean" sentinel

The skill is not a linter. A trace that finds no failing hop is itself the "no finding" case; the trace itself is the output. There is no canonical "all clean" line; the framework's output is always a trace.
