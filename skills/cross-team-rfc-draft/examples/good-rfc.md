# RFC-007: Migrate Auth Token Format from JWT-v1 to JWT-v2

<!-- Good counterpart of bad-rfc.md.
     Fix 1 (SECTIONS): all 9 sections present.
     Fix 2 (CONSUMER-IMPACT): per-team attribution for the 4 consumers.
     Fix 3 (MIGRATION-WINDOW): per-consumer window with the consumer's calendar.
     Fix 4 (ROLLOUT): flag, cohort, kill switch, rollback all named.
     Fix 5 (OPEN-QUESTIONS): the 3 open questions named with the assignee.
-->

## 1. Title

RFC-007: Migrate auth token format from JWT-v1 to JWT-v2.

## 2. Status

`DRAFT` (filed 2026-06-10).

## 3. Author + team + date

- **Author:** `@alice` (`@org/auth-core`).
- **Date filed:** 2026-06-10.

## 4. Summary

The auth service issues access tokens as JWT-v1 (header `alg: HS256`, claim set v1). We propose migrating to JWT-v2 (header `alg: RS256`, claim set v2 with `tenant_id` and `region`). The migration is a token-format change that affects every service that calls `auth.verify(token)`. The migration is required because (a) HS256 is a symmetric algorithm and the secret rotation cost is high; (b) the v1 claim set does not carry `tenant_id`, which the new billing flow needs.

The alternative — staying on JWT-v1 and adding `tenant_id` as a separate claim — is rejected because it would require a parallel claim-set evolution alongside the format migration, doubling the surface area.

## 5. Motivation

- **HS256 secret rotation cost.** The current HS256 secret is shared across all 47 services that verify tokens. The last rotation took 6 weeks of coordinated deploys; the next rotation (per the security team's annual calendar) is in Q4 2026.
- **v1 claim set is missing `tenant_id`.** The new billing flow needs `tenant_id` in the claim set; today, services fetch `tenant_id` by a separate DB lookup, which is a hot path.
- **Incident #INC-2025-918.** Last quarter's incident showed a service that mis-parsed a v1 token under load; the mis-parse was hidden by the symmetric-key fallback. RS256 with a public-key verify would have surfaced the mis-parse.

## 6. Detailed design

- **Header:** `alg: RS256` (was `alg: HS256`).
- **Key:** asymmetric. The auth service holds the private key (rotated quarterly via KMS); the consumer services hold the public key (cached in-process, refreshed every 5 minutes).
- **Claim set:** adds `tenant_id` and `region`; deprecates `user_org_id` (replaced by `tenant_id`).
- **Verification:** the new `auth.verify(token)` returns a v2-typed result. A `auth.verify_v1(token)` shim is provided for the migration window; it returns the v1 fields and a deprecation warning.
- **Library:** the existing `auth-client` Python / Node / Go packages gain a v2 client. The v1 client is deprecated but functional for the migration window.
- **Compatibility matrix:**

  | Consumer | v1 only | v2 only | dual-mode |
  |---|---|---|---|
  | `services/billing/...` | today | at migration | during |
  | `services/orders/...` | today | at migration | during |
  | `services/fulfillment/...` | today | at migration | during |
  | `services/reports/...` | today | at migration | during |

## 7. Consumer-team impact

### `@org/billing-core`

- **Named approver:** `@bob`
- **Migration window:** 8 weeks (consumer has a quarterly release; the next release is 2026-08-15)
- **Per-consumer cost:** medium (1-3 dev-weeks; the billing service has 6 call sites to `auth.verify`)
- **Release-train slot:** 2026-08-15 (next billing train)
- **Suggested review order:** 1st (deepest dependency; needs the most migration time)

### `@org/orders-core`

- **Named approver:** `@carol`
- **Migration window:** 6 weeks
- **Per-consumer cost:** small (under 1 dev-week; 2 call sites)
- **Release-train slot:** 2026-08-22 (next orders train)
- **Suggested review order:** 2nd

### `@org/fulfillment-core`

- **Named approver:** `@dave`
- **Migration window:** 12 weeks (consumer has a quarterly release; the next release is 2026-09-12)
- **Per-consumer cost:** large (4+ dev-weeks; 12 call sites + a stateful batch job)
- **Release-train slot:** 2026-09-12
- **Suggested review order:** 4th (release-train-locked; review the RFC now, ack in the next train)

### `@org/reports-core`

- **Named approver:** `@eve`
- **Migration window:** 4 weeks
- **Per-consumer cost:** small (1 dev-week; the reports service is read-only against auth)
- **Release-train slot:** 2026-08-01 (next reports batch window)
- **Suggested review order:** 3rd

## 8. Rollout plan

- **Feature flag:** `auth.token_format` — values `v1` (default today), `v2`, `dual`. Flag flips from `v1` to `dual` at the start of the migration window; flips from `dual` to `v2` after the last consumer has migrated; the v1 path is removed in a follow-up.
- **Cohort:** first 100% of internal services, then 100% of external API consumers (via the auth API gateway's content negotiation).
- **SLO gate:** `auth.verify.p99 < 50ms` (today's SLO); `auth.verify.error_rate < 0.1%`.
- **Kill switch:** `POST /admin/flags/auth.token_format/off` (reverts to v1-only).
- **Rollback hook:** the flag flip + a re-deploy of the auth service with the v1 signing path re-enabled.
- **Sunset date:** 2027-01-01 (v1 signing path is removed; v2 is the only path).

## 9. Open questions + decision log

### Open questions

- **Q1 (assignee: `@alice`):** Should we add `region` to the v2 claim set, or keep it as a service-side lookup? — **assignee** `@alice`, due 2026-06-20.
- **Q2 (assignee: `@bob`):** Does the billing flow's audit log need the JWT's `kid` (key id) for forensics? — **assignee** `@bob`, due 2026-06-20.
- **Q3 (assignee: `@dave`):** Can the fulfillment batch job run with a `dual`-mode token, or does it need a deterministic token format per batch? — **assignee** `@dave`, due 2026-06-25.

### Decision log

- **2026-06-10:** Rejected the alternative of staying on JWT-v1 with a separate `tenant_id` claim. Reason: doubles the surface area; the v1 claim set would still need to be maintained for the migration window anyway.
- **2026-06-10:** Chose RS256 over EdDSA. Reason: RS256 has wider library support across the consumer languages (Node, Python, Go, Java, .NET); EdDSA is supported but adds an extra integration risk for the 4 consumer teams.
