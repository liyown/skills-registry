# Checkout 5xx Spike — Good Trace

<!-- Good counterpart of bad-trace.md.
     Fix 1 (HOP-WALK): per-hop walk from the entry point to persistence.
     Fix 2 (PRIMARY-LOCUS): the first hop that degraded is named.
     Fix 3 (BUDGET-VS-OBSERVED): latency budget vs. observed per hop.
     Fix 4 (PAGING-LIST): per-team paging list.
     Fix 5 (CUSTOMER-IMPACT): customer-impact framing.
-->

## Incident

**When:** 2026-06-08 14:02 local.
**Symptom:** checkout success rate dropped from 99.5% to 92%.
**Entry point:** `POST /api/v2/checkout`.

## 1. Primary incident locus

The first hop that degraded is **`services/payments/PaymentService.chargeToStripe`** at `services/payments/PaymentService.java:42`. Latency budget 800ms, observed 6.2s at 14:02–14:18. The error rate from this hop is 7.4%, accounting for the 7.5pp drop in checkout success.

## 2. Per-hop trace

- **`controllers/CheckoutController.checkout` — `controllers/CheckoutController.java:18`**
  - Owning team: `@org/checkout-core`
  - Latency budget: 200ms
  - Observed: 180ms (no degradation)
  - Failure mode: none
  - Downstream: `services/cart/CartService.priceCart` at `:24`

- **`services/cart/CartService.priceCart` — `services/cart/CartService.java:24`**
  - Owning team: `@org/cart-core`
  - Latency budget: 200ms
  - Observed: 190ms (no degradation)
  - Failure mode: none
  - Downstream: `services/payments/PaymentService.chargeToStripe` at `:42`

- **`services/payments/PaymentService.chargeToStripe` — `services/payments/PaymentService.java:42` ← PRIMARY LOCUS**
  - Owning team: `@org/payments-core`
  - Latency budget: 800ms
  - Observed: 6.2s at 14:02–14:18 (7.7× budget)
  - Failure mode: timeout (calls to Stripe timeout at 5s; the SDK's default)
  - Downstream: `external/Stripe.paymentIntents.create`

- **`external/Stripe.paymentIntents.create`**
  - Owning team: `@external/stripe`
  - Latency budget: 5000ms (SDK default)
  - Observed: 5000ms (timeout, not a Stripe slowdown)
  - Failure mode: timeout
  - Downstream: none (terminal)

- **`services/payments/PaymentService.markPaid` — `services/payments/PaymentService.java:71`**
  - Owning team: `@org/payments-core`
  - Latency budget: 100ms
  - Observed: 95ms (no degradation; the call is reached only on success, and the failure is upstream)
  - Failure mode: none
  - Downstream: `repos/OrderRepository.updateStatus` at `:18`

- **`repos/OrderRepository.updateStatus` — `repos/OrderRepository.java:18`**
  - Owning team: `@org/payments-core`
  - Latency budget: 50ms
  - Observed: 45ms (no degradation)
  - Failure mode: none
  - Downstream: none (terminal)

## 3. Per-team paging list

- `@org/payments-core` — on-call: `@dave` (PagerDuty schedule `payments-oncall`)
- `@org/checkout-core` — informational, not paged; the failure is upstream
- `@org/cart-core` — informational, not paged
- `@external/stripe` — vendor support, ticket opened

## 4. Customer-impact framing

- 7.5pp drop in checkout success rate over 16 minutes.
- Approximately 1,200 failed checkouts in the window.
- Region: `us-east-1` (the only Stripe region this service routes to).
- Customer classes affected: all customers using `card` payment; `wallet` and `bank` are unaffected (different payment-method code path).
- Customer-visible: 5xx response; cart state preserved; no charge attempted (good).
