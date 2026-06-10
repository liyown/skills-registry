# Order Service — DevAgent

## Overview

`BadOrderService` owns order creation, payment, and cancellation.
It is the only service that talks to `payments.Charge`.

## Public API

- `BadOrderService.pay(userId: Long, orderId: Long): Unit`
  - Old: signature used to be `pay(userId: Long, orderId: Long): Result<Order, PayError>`.
  - Reads the order, checks the status, charges the user, and marks paid.
- `BadOrderService.batchCancel(ids: List<Long>): Unit`
  - Cancels a list of orders in one call. **Idempotent** — calling
    twice with the same ids is a no-op.

## Invariants

- `pay` is idempotent: a second call with the same `orderId`
  on an already-paid order is a no-op.
- `batchCancel` is atomic: all orders are cancelled or none are.
- The status field only ever transitions
  `UNPAID → PAYING → PAID` and `UNPAID → CANCELLED`. There is
  no `PAYING → UNPAID` rollback in normal flow.

## Anti-patterns

- Do not import `com.example.order.internal.LegacyCharge`. Use
  `payments.Charge` only.

## Dependencies

- `payments` (for the charge call)
- `kafka` (for emitting `OrderPaid` events)

## File Map

- `./src/main/java/com/example/order/BadOrderService.java` —
  the main service.
- `./src/main/java/com/example/order/OrderMapper.java`
- `./src/main/java/com/example/order/AccountMapper.java`
- `./src/main/java/com/example/order/PaymentClient.java`
- `./docs/knowledge/order-history.md` — historical context.
