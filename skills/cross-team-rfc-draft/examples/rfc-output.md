# RFC-007: Migrate Auth Token Format from JWT-v1 to JWT-v2

## 1. Title

RFC-007: Migrate auth token format from JWT-v1 to JWT-v2.

## 2. Status

`DRAFT` (filed 2026-06-10).

## 3. Author + team + date

- **Author:** `@alice` (`@org/auth-core`).
- **Date filed:** 2026-06-10.

## 4. Summary

The auth service issues access tokens as JWT-v1 (`alg: HS256`, claim set v1). We propose migrating to JWT-v2 (`alg: RS256`, claim set v2 with `tenant_id` and `region`). The migration is a token-format change that affects every service that calls `auth.verify(token)`. Required by (a) HS256 secret-rotation cost, (b) v1 claim set missing `tenant_id`, (c) incident #INC-2025-918.

Alternative rejected: stay on JWT-v1 with a separate `tenant_id` claim.

## 5. Motivation

- HS256 secret rotation cost (next rotation Q4 2026, 6-week effort).
- v1 claim set missing `tenant_id`; billing flow needs it.
- Incident #INC-2025-918 showed a service mis-parsing a v1 token under load; RS256 would have surfaced the mis-parse.

## 6. Detailed design

- Header `alg: RS256`; asymmetric key; private in auth service (KMS-rotated), public cached in consumers (5-min refresh).
- Claim set adds `tenant_id` + `region`; deprecates `user_org_id`.
- New `auth.verify(token)` returns v2-typed result; `auth.verify_v1(token)` shim provided for migration window.
- `auth-client` Python / Node / Go packages gain a v2 client; v1 client deprecated but functional.

## 7. Consumer-team impact

| Team | Approver | Window | Cost | Release-train | Review order |
|---|---|---|---|---|---|
| `@org/billing-core` | `@bob` | 8w | medium (1-3 dev-wk, 6 call sites) | 2026-08-15 | 1st |
| `@org/orders-core` | `@carol` | 6w | small (<1 dev-wk, 2 call sites) | 2026-08-22 | 2nd |
| `@org/reports-core` | `@eve` | 4w | small (1 dev-wk, read-only) | 2026-08-01 | 3rd |
| `@org/fulfillment-core` | `@dave` | 12w | large (4+ dev-wk, 12 call sites + stateful batch) | 2026-09-12 | 4th |

## 8. Rollout plan

- **Flag:** `auth.token_format` (`v1` → `dual` → `v2` → v1 removed).
- **Cohort:** internal services first, then external API consumers via the auth API gateway's content negotiation.
- **SLO gate:** `auth.verify.p99 < 50ms`, `auth.verify.error_rate < 0.1%`.
- **Kill switch:** `POST /admin/flags/auth.token_format/off` → reverts to v1-only.
- **Rollback hook:** flag flip + re-deploy auth with v1 signing path re-enabled.
- **Sunset date:** 2027-01-01.

## 9. Open questions + decision log

### Open

- **Q1** (`@alice`, due 2026-06-20): Should `region` be in the v2 claim set, or service-side lookup?
- **Q2** (`@bob`, due 2026-06-20): Does the billing audit log need the JWT's `kid` for forensics?
- **Q3** (`@dave`, due 2026-06-25): Can the fulfillment batch job run with a `dual`-mode token?

### Closed

- **2026-06-10** — Rejected "stay on JWT-v1 with separate `tenant_id`": doubles surface area, v1 still maintained.
- **2026-06-10** — Chose RS256 over EdDSA: wider library support across consumer languages.
