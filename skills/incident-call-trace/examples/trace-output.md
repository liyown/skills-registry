# Checkout 5xx Spike — Trace

**Incident:** 2026-06-08 14:02 local. Checkout success rate 99.5% → 92% over 16 minutes.
**Entry point:** `POST /api/v2/checkout`

## 1. Primary incident locus

`services/payments/PaymentService.chargeToStripe` at `services/payments/PaymentService.java:42`. Latency budget 800ms, observed 6.2s. Error rate 7.4%.

## 2. Per-hop trace

| # | Call site | Owning team | Budget | Observed | Failure mode | Downstream |
|---|---|---|---|---|---|---|
| 1 | `controllers/CheckoutController.checkout` (`CheckoutController.java:18`) | `@org/checkout-core` | 200ms | 180ms | none | `CartService.priceCart:24` |
| 2 | `services/cart/CartService.priceCart` (`CartService.java:24`) | `@org/cart-core` | 200ms | 190ms | none | `PaymentService.chargeToStripe:42` |
| 3 | `services/payments/PaymentService.chargeToStripe` (`PaymentService.java:42`) | `@org/payments-core` | 800ms | 6.2s | timeout (Stripe 5s default) | `external/Stripe.paymentIntents.create` |
| 4 | `external/Stripe.paymentIntents.create` | `@external/stripe` | 5000ms | 5000ms | timeout (SDK default) | none |
| 5 | `services/payments/PaymentService.markPaid` (`PaymentService.java:71`) | `@org/payments-core` | 100ms | 95ms | none (not reached) | `OrderRepository.updateStatus:18` |
| 6 | `repos/OrderRepository.updateStatus` (`OrderRepository.java:18`) | `@org/payments-core` | 50ms | 45ms | none (not reached) | none |

## 3. Per-team paging list

- **`@org/payments-core`** — page: `@dave` (PagerDuty schedule `payments-oncall`). Primary responder; owns the failing hop.
- `@org/checkout-core` — informational, not paged; the failure is upstream.
- `@org/cart-core` — informational, not paged.
- `@external/stripe` — vendor support, ticket opened.

## 4. Customer-impact framing

- 7.5pp drop in checkout success rate over 16 minutes (14:02–14:18).
- ~1,200 failed checkouts in the window.
- Region: `us-east-1` only.
- Customer classes: `card` payment only; `wallet` and `bank` unaffected.
- Customer-visible: 5xx response; cart state preserved; no charge attempted.

## Follow-up

- Hand off to postmortem authoring. The `observability-coverage-audit` skill (planned v0.6.1) will check whether the timeout has a circuit breaker, a retry, and an alert; if any is missing, the gap is a follow-up action.
