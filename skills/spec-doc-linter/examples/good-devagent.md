# Order Service — DevAgent

<!-- Good counterpart of bad-devagent.md.
     Fix 1 (A.1): renamed batchCancel → cancelMany, signature updated.
     Fix 2 (B.1): idempotency/atomicity claims dropped to per-row truth.
     Fix 3 (A.5): kafka removed, sentry added, accounting note added.
     Fix 4 (A.4): broken ./docs/knowledge/order-history.md reference removed.
-->

## Overview

`BadOrderService` owns order creation, payment, and cancellation.
It is the only service that talks to `payments.Charge`.

## Public API

- `BadOrderService.pay(userId: Long, orderId: Long): Unit`
- `BadOrderService.cancelMany(ids: List<Long>): Int`
  - Returns the number of orders that transitioned from
    `UNPAID` to `CANCELLED`. Not idempotent; not atomic; runs
    one `UPDATE ... WHERE status = 'UNPAID'` per id.

## Invariants

- `pay` is idempotent: a second call with the same `orderId`
  on an already-paid order is a no-op.
- The status field only ever transitions
  `UNPAID → PAYING → PAID` and `UNPAID → CANCELLED`. There is
  no `PAYING → UNPAID` rollback in normal flow.

## Anti-patterns

- Do not import `com.example.order.internal.LegacyCharge`. Use
  `payments.Charge` only.

## Dependencies

- `payments` (for the charge call)
- `sentry` (for emitting `OrderFailed` events to the
  observability backend; replaces the old `kafka` topic)

## File Map

- `./src/main/java/com/example/order/BadOrderService.java` —
  the main service.
- `./src/main/java/com/example/order/OrderMapper.java`
- `./src/main/java/com/example/order/AccountMapper.java`
- `./src/main/java/com/example/order/PaymentClient.java`
