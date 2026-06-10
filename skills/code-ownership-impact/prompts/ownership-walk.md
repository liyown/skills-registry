# Ownership Walk Framework

> See also: ../SKILL.md

The framework is a per-team review packet. For each team in the diff, the framework says what to record.

## The 5 cells per team

1. **Owned paths in the diff** — the file paths the team owns, sourced from the path → team map.
2. **Named approvers** — the GitHub handles from the `CODEOWNERS` entry for those paths.
3. **Suggested review deadline** — the next release-train slot the team uses, or "ASAP" for hotfixes. If the team has no release train, name the team's weekly review window.
4. **Cross-team impact** — for shared-kernel changes, list the team's owned services that consume the kernel. If a service is the consumer, the team is a downstream approver even if no path in the diff is theirs.
5. **Ack-channel** — where to ping: a thread in the team's chat, a CODEOWNERS `@`-ping, an RFC subscribe, or a calendar slot.

## The 3 review-packet sections

The final packet is structured as three sections, in this order:

1. **Single-team ack** — the author's team. The author is the primary reviewer; the team's secondary reviewer is the code-owner of the largest changed path.
2. **Cross-team approvers** — every team whose CODEOWNERS entry covers a path in the diff. One row per team with the 5 cells above.
3. **Shared-kernel subscribers** — every team whose service consumes a shared-kernel path. One row per service. The team's secondary reviewer ack is required.

## The single-team packet (no boundary crossing)

A diff that touches no shared kernel and no per-team boundary produces a single-team review packet: the author's team, the named approver, the suggested review deadline. That is itself the correct output — not an error.

## The shared-kernel detection rule

A path is shared-kernel if any of:

- It is under a directory named `common/`, `shared/`, `pkg/`, `platform/`, `internal/`, `core/`, `kernels/`.
- It matches the repo's `*.shared-kernel*` glob (some monorepos use this).
- It is in the `CODEOWNERS` `@org/platform` group (a marker that the path is platform-owned, not team-owned).
- A `git log` over the file shows ≥ 3 teams as recent authors (mixed ownership is the signature of a shared kernel).

## The review-order suggestion

For each team, the suggested review order is:

1. Single-team ack (the author's team) — first, fastest.
2. Cross-team approvers for the changed paths — second, in CODEOWNERS order.
3. Shared-kernel subscribers — third, in dependency order (most-dependent first).
4. Release-train-locked teams — last, with a hard deadline tied to the next train slot.

The "release-train-locked" step is the one that prevents a small change from being blocked by a team's quarterly release calendar.

## The "all clean" sentinel

The skill is not a linter; it does not have a single no-finding line. A diff that produces a single-team packet is the no-finding case. The packet itself is the output.
