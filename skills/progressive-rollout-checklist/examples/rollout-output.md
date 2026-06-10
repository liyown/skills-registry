# Pricing Engine Rollout — Packet

## Change

Flipping `use_new_pricing_engine` default from `false` to `true`.

## Flag rows

### `use_new_pricing_engine`

| Cell | Value |
|---|---|
| Flag name | `use_new_pricing_engine` |
| Read sites | `services/pricing/PricingService.java:31`, `services/pricing/PricingService.java:54`, `admin/dashboards/PricingDashboard.tsx:18` |
| Write sites | `admin/endpoints/FlagController.java:42`, `services/cohort/CohortExpander.java:18`, `services/killswitch/EmergencyStop.java:7` |
| Default value | `false` |
| Proposed new default | `true` (post-cohort 100%) |
| Cohort definition | `tenant_id IN (1, 17, 42, 88, 144)` — five pilot `enterprise` tenants |
| SLO gate | `error_rate < 0.1% AND p99 < 800ms` over 1h; expand to 10% at 1h open, 100% at 24h open |
| Kill switch | `POST /admin/flags/use_new_pricing_engine/off`; owner `@dave` (`@org/payments-core` on-call) |
| Rollback hook | (1) flip flag to `false` via kill switch; (2) re-run `scripts/pricing-backfill.sh --from-pricing-engine=legacy`; owner `@alice` (`@org/pricing-core`) |
| Sunset date | 2026-09-01 |

## Cross-flag dependencies

- `use_new_pricing_engine` depends on `bulk_pricing_cache_warm` being `true`. Cohort does not expand otherwise.

## Ship-readiness verdict

- `use_new_pricing_engine`: **READY**
