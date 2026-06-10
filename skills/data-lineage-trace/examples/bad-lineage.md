# Orders National ID — Lineage

## Field

`Order.national_id` — the customer's national identifier, captured at order placement.

## Producers

- `services/orders/OrderService.java:42` — written from the request payload (`createOrderRequest.getNationalId()`).
- `services/orders/OrderController.java:18` — passes the field through from the HTTP body.

## Transformers

- `services/orders/OrderMapper.java:31` — INSERT into `orders` table, no cast.

## Sinks

- **Transactional** — `postgres://orders.orders`. Consumer: `services/billing/InvoiceService.java:14` (reads `Order.national_id` to populate the invoice). Control: `audit-logged` (yes, via `OrderAuditLog`), `encrypted-at-rest` (no — column is plaintext), `consent-checked` (no).
- **Analytics** — `snowflake://analytics.orders_fact`. Consumer: the BI dashboard. Control: `audit-logged` (no), `encrypted-at-rest` (no), `consent-checked` (no).
- **Log** — `logback://orders.events`. Consumer: the log aggregator. Control: `audit-logged` (n/a for logs), `encrypted-at-rest` (no).
- **External** — `payments.Stripe.paymentIntents.create` (Stripe receives the field as part of the customer record). Control: `audit-logged` (yes, via Stripe), `consent-checked` (no — the field is sent without a consent check at the order-create boundary).

## Controls

- `audit-logged` is present only on the transactional sink and the external sink.
- `encrypted-at-rest` is present on no sink — the field is plaintext end to end.
- `consent-checked` is present on no sink — the request payload is accepted unconditionally.

## Deletion Targets (for GDPR right-to-erasure)

- `postgres://orders.orders` — held by `OrderService`. In the erasure worker? **Unknown — not in the worker's coverage list.**
- `snowflake://analytics.orders_fact` — held by the daily export job. In the erasure worker? **Unknown — analytics is not in the worker's list.**
- `Stripe` — held by Stripe. Erasure requires a separate Stripe API call; the worker does not invoke it.
