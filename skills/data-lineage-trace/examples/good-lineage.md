# Orders National ID — Lineage

<!-- Good counterpart of bad-lineage.md.
     Fix 1 (LIN-ENCRYPT): column is encrypted at rest on every sink.
     Fix 2 (LIN-CONSENT): consent is checked at the HTTP boundary before
       the field is accepted into the request.
     Fix 3 (LIN-DELETION): deletion-propagation list is complete and
       matches the actual sinks.
-->

## Field

`Order.national_id` — the customer's national identifier, captured at order placement, with explicit consent.

## Producers

- `services/orders/OrderService.java:42` — written from the request payload after the consent interceptor at `OrderController.java:18` returns success.
- `services/orders/OrderController.java:18` — passes the field through from the HTTP body once consent is verified.

## Transformers

- `services/orders/OrderMapper.java:31` — INSERT into `orders` table; column is `BYTEA` (encrypted via `pgcrypto`).

## Sinks

- **Transactional** — `postgres://orders.orders`. Consumer: `services/billing/InvoiceService.java:14`. Control: `audit-logged` (yes, via `OrderAuditLog`), `encrypted-at-rest` (yes, `pgcrypto`), `consent-checked` (yes, at the boundary).
- **Analytics** — `snowflake://analytics.orders_fact`. Consumer: the BI dashboard. Control: `audit-logged` (no, BI is read-only), `encrypted-at-rest` (yes, Snowflake column-level encryption), `consent-checked` (no, analytics does not require per-row consent).
- **Log** — `logback://orders.events`. Consumer: the log aggregator. Control: `audit-logged` (n/a for logs), `encrypted-at-rest` (no, logs are redacted — see `services/orders/LogRedactor.java:7` which replaces the field with `***` before the log line is written).
- **External** — `payments.Stripe.paymentIntents.create`. Control: `audit-logged` (yes, via Stripe), `consent-checked` (yes — the consent interceptor at the boundary also gates the Stripe call).

## Controls

- `audit-logged` is present on the transactional sink and the external sink; analytics and logs are exempt.
- `encrypted-at-rest` is present on transactional and analytics; logs are redacted instead.
- `consent-checked` is present at the boundary; analytics is exempt (it operates on already-consented rows).

## Deletion Targets (for GDPR right-to-erasure)

- `postgres://orders.orders` — held by `OrderService`. In the erasure worker? **Yes, in `ErasureWorker.deleteOrder`.**
- `snowflake://analytics.orders_fact` — held by the daily export job. In the erasure worker? **Yes, in `ErasureWorker.purgeAnalytics`.**
- `Stripe` — held by Stripe. Erasure requires a separate Stripe API call; the worker invokes it via `ErasureWorker.purgeStripeCustomer`.
