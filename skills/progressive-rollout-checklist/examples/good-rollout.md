# Pricing Engine Rollout — Good Packet

<!-- Good counterpart of bad-rollout.md.
     Fix 1 (READ-SITES): every read site named with path + line.
     Fix 2 (WRITE-SITES): every write site named.
     Fix 3 (COHORT): cohort defined as a specific tenant-id list.
     Fix 4 (SLO-GATE): SLO gate named with metric + threshold.
     Fix 5 (KILL-SWITCH): kill switch action + owner named.
     Fix 6 (ROLLBACK): rollback hook named (data + flag).
-->

## Change

Flipping the default of `use_new_pricing_engine` from `false` to `true`.

## Flag row

- **Flag name:** `use_new_pricing_engine`
- **Read sites:**
  - `services/pricing/PricingService.priceOrder` (`services/pricing/PricingService.java:31`) — boolean switch on the compute path
  - `services/pricing/PricingService.priceBulk` (`services/pricing/PricingService.java:54`) — same switch, bulk variant
  - `admin/dashboards/PricingDashboard` (`admin/dashboards/PricingDashboard.tsx:18`) — admin UI shows the flag status
- **Write sites:**
  - `admin/endpoints/FlagController` (`admin/endpoints/FlagController.java:42`) — admin override
  - `services/cohort/CohortExpander.expand` (`services/cohort/CohortExpander.java:18`) — automated expansion (gated by the SLO below)
  - `services/killswitch/EmergencyStop` (`services/killswitch/EmergencyStop.java:7`) — global off
- **Default value:** `false`
- **Proposed new default:** `true` (after the cohort reaches 100%)
- **Cohort definition:** `tenant_id IN (1, 17, 42, 88, 144)` — five pilot tenants, all on the `enterprise` plan
- **SLO gate:** `error_rate < 0.1% AND p99 < 800ms` over a 1-hour window. The `CohortExpander` reads this metric from the metrics backend; if the gate is open for 1 hour, expand to 10% of tenants; if open for 24 hours, expand to 100%.
- **Kill switch:** `POST /admin/flags/use_new_pricing_engine/off` (handled by `FlagController:42`). Owner: `@org/payments-core` on-call (`@dave`).
- **Rollback hook:** two-part.
  1. **Flag-level:** flip `use_new_pricing_engine` to `false` via the kill switch above.
  2. **Data-level:** re-run the legacy-pricing backfill via `scripts/pricing-backfill.sh --from-pricing-engine=legacy`. This recomputes every price the new engine wrote and overwrites the DB. Owner: `@org/pricing-core` (`@alice`).
- **Sunset date:** 2026-09-01 (the flag is fully removed; the cohort-expansion infrastructure is deleted; the new engine becomes the only path).

## Cross-flag dependencies

- `use_new_pricing_engine` depends on `bulk_pricing_cache_warm` — the bulk-pricing cache must be warm for the new engine to be the bottleneck. If `bulk_pricing_cache_warm` is `false`, the cohort does not expand regardless of the SLO gate.

## Ship-readiness verdict

- `use_new_pricing_engine`: **READY** (all 8 cells filled; cross-flag dependency noted).
