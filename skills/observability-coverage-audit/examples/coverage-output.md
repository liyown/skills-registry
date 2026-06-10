# Payments Coverage Audit — Matrix

**In-scope path:** `services/payments/PaymentService` — every `throw` site in the charge / refund / payout flow.

## Per-site coverage rows

| # | Site | Test | Alert | Runbook | Classification |
|---|---|---|---|---|---|
| 1 | `PaymentService.java:42` (`throw new PaymentDeclinedException`) | `PaymentServiceTest.test_charge_decline` | `alert.payments.decline_rate` | `runbooks/payments.md#decline` | `COV-COMPLETE` |
| 2 | `PaymentService.java:71` (`throw new StripeTimeoutException`) | `PaymentServiceTest.test_charge_stripe_timeout` (real Stripe test mode) | `MISSING` | `MISSING` | `COV-SILENT` |
| 3 | `PaymentService.java:104` (`throw new RefundFailedException`) | `PaymentServiceTest.test_refund_failure` (mocks gateway — does NOT exercise production) | `MISSING` | `runbooks/payments.md#refund-failed` | `COV-NO-TEST-AND-NO-ALERT` |
| 4 | `PaymentService.java:131` (`throw new InsufficientFundsException`) | `MISSING` | `MISSING` | `MISSING` | `COV-SILENT` |
| 5 | `PaymentService.java:158` (`throw new PayoutQueueFullException`) | `PayoutServiceTest.test_payout_queue_full` | `alert.payments.payout_queue_depth` | `MISSING` | `COV-NO-RUNBOOK` |
| 6 | `PaymentService.java:189` (`throw new CurrencyMismatchException`) | `PaymentServiceTest.test_currency_mismatch` | `MISSING` | `runbooks/payments.md#currency-mismatch` | `COV-NO-ALERT` |

## Gap findings (priority order)

1. **`COV-SILENT`** `PaymentService.java:131` (Insufficient funds) — all three missing.
2. **`COV-SILENT`** `PaymentService.java:71` (Stripe timeout) — alert + runbook missing.
3. **`COV-NO-TEST-AND-NO-ALERT`** `PaymentService.java:104` (Refund failure) — test mocks the gateway, no alert.
4. **`COV-NO-RUNBOOK`** `PaymentService.java:158` (Payout queue full) — no runbook section.
5. **`COV-NO-ALERT`** `PaymentService.java:189` (Currency mismatch) — no alert.

## Suggested fixes

1. `PaymentService.java:131` — add unit test on `InsufficientFundsException`; add alert `payments.insufficient_funds_rate`; add runbook `runbooks/payments.md#insufficient-funds`.
2. `PaymentService.java:71` — add real-Stripe-mode test (no mock); add alert `payments.stripe_timeout_rate`; add runbook `runbooks/payments.md#stripe-timeout`.
3. `PaymentService.java:104` — replace mocked test with real-Stripe-mode test; add alert `payments.refund_failed_rate`.
4. `PaymentService.java:158` — add runbook `runbooks/payments.md#payout-queue-full`.
5. `PaymentService.java:189` — add alert `payments.currency_mismatch_rate`.
