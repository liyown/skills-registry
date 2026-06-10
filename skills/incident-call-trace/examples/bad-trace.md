# Checkout 5xx Spike — Bad Trace

## Incident

**When:** 2026-06-08 14:02 local.
**Symptom:** checkout success rate dropped from 99.5% to 92%.
**Entry point:** `POST /api/v2/checkout`.

## The Trace (as the on-call would otherwise write it)

> The 5xx is in checkout. The payments service is probably slow.

## What's wrong with this trace

- **No per-hop walk.** "The payments service is probably slow" is a guess, not a trace. The on-call has not walked from `POST /api/v2/checkout` through auth, cart, payment-gateway, and persistence.
- **No primary incident locus.** The trace does not name the first hop that degraded. The on-call cannot tell whether the failure is in the gateway, the persistence layer, or the load balancer.
- **No latency budget vs. observed.** "Slow" is qualitative; a budget (e.g. "payment-gateway p99 ≤ 800ms") is quantitative.
- **No per-team paging list.** Every hop has an owning team; the trace does not name any of them.
- **No customer-impact framing.** "5xx" is a rate; the customer impact (how many customers, which regions, which payment methods) is not stated.
