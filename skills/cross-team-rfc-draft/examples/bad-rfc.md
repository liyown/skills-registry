# Auth Token Format RFC — Bad Draft

## The RFC (as the author would otherwise send it)

> # Migrate to JWT-v2
>
> We're switching auth tokens to JWT-v2. Review and let us know.

## What's wrong with this RFC

- **Missing 7 of 9 sections.** No `Status`, no `Author`, no `Motivation`, no `Detailed design`, no `Consumer-team impact`, no `Rollout plan`, no `Open questions`.
- **No consumer-team attribution.** "Review and let us know" is one broadcast. The 4 consumer teams each have their own reviewer, calendar, and migration cost.
- **No migration window.** Switching token format is a coordinated change; "switch" without a window is a forced migration.
- **No rollout plan.** A token-format change needs a flag flip, a cohort, and a kill switch. None are named.
- **No open questions.** Every RFC has open questions at draft time; pretending there are none is a review smell.
