# Payments Coverage Audit — Good Matrix

<!-- Good counterpart of bad-coverage.md.
     Fix 1 (PER-SITE-ROWS): every throw site named with path + line.
     Fix 2 (TEST-COL): test class + function named; mocked tests not counted.
     Fix 3 (ALERT-COL): alert rule id named.
     Fix 4 (RUNBOOK-COL): runbook path + section named.
     Fix 5 (CLASSIFICATION): per-site classification finding emitted.
     Fix 6 (SUGGESTED-FIX): smallest coverage to add named per gap.
-->

## In-scope path

`services/payments/PaymentService` — every `throw` site in the charge / refund / payout flow.

## Per-site coverage rows

| # | Site | Test | Alert | Runbook | Classification |
|---|---|---|---|---|---|
| 1 | `services/payments/PaymentService.java:42` — `throw new PaymentDeclinedException` | `PaymentServiceTest.test_charge_decline` | `alert.payments.decline_rate` (page: `@org/payments-core` on-call) | `runbooks/payments.md#decline` | `COV-COMPLETE` |
| 2 | `services/payments/PaymentService.java:71` — `throw new StripeTimeoutException` | `PaymentServiceTest.test_charge_stripe_timeout` (no mock — uses Stripe test mode) | `MISSING` | `MISSING` | `COV-SILENT` |
| 3 | `services/payments/PaymentService.java:104` — `throw new RefundFailedException` | `PaymentServiceTest.test_refund_failure` (mocks the gateway — does NOT exercise production path) | `MISSING` | `runbooks/payments.md#refund-failed` | `COV-NO-TEST-AND-NO-ALERT` |
| 4 | `services/payments/PaymentService.java:131` — `throw new InsufficientFundsException` | `MISSING` | `MISSING` | `MISSING` | `COV-SILENT` |
| 5 | `services/payments/PaymentService.java:158` — `throw new PayoutQueueFullException` | `PayoutServiceTest.test_payout_queue_full` | `alert.payments.payout_queue_depth` (page: `@org/payments-core` on-call) | `MISSING` | `COV-NO-RUNBOOK` |
| 6 | `services/payments/PaymentService.java:189` — `throw new CurrencyMismatchException` | `PaymentServiceTest.test_currency_mismatch` | `MISSING` | `runbooks/payments.md#currency-mismatch` | `COV-NO-ALERT` |

## Gap findings

- **`COV-SILENT`** `PaymentService.java:71` — Stripe timeout. Both alert and runbook missing.
- **`COV-SILENT`** `PaymentService.java:131` — Insufficient funds. All three missing.
- **`COV-NO-TEST-AND-NO-ALERT`** `PaymentService.java:104` — Refund failure. Test mocks the gateway (does not exercise production); no alert.
- **`COV-NO-RUNBOOK`** `PaymentService.java:158` — Payout queue full. Runbook missing.
- **`COV-NO-ALERT`** `PaymentService.java:189` — Currency mismatch. Alert missing.

## Suggested fixes (priority order)

1. **`PaymentService.java:131` (COV-SILENT)** — add a unit test that exercises the `InsufficientFundsException` site; add an alert on the metric `payments.insufficient_funds_rate`; add a runbook section `runbooks/payments.md#insufficient-funds`.
2. **`PaymentService.java:71` (COV-SILENT)** — add a unit test that exercises the `StripeTimeoutException` site against Stripe test mode (no mock); add an alert on `payments.stripe_timeout_rate`; add a runbook section `runbooks/payments.md#stripe-timeout`.
3. **`PaymentService.java:104` (COV-NO-TEST-AND-NO-ALERT)** — replace the mocked test with a real-Stripe-mode test; add an alert on `payments.refund_failed_rate`.
4. **`PaymentService.java:158` (COV-NO-RUNBOOK)** — add a runbook section `runbooks/payments.md#payout-queue-full`.
5. **`PaymentService.java:189` (COV-NO-ALERT)** — add an alert on `payments.currency_mismatch_rate`.
