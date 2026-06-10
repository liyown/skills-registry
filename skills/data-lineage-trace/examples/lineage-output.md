# Lineage Report

## Orders.national_id

### Producers

- `services/orders/OrderService.java:42` — written from the request payload.
- `services/orders/OrderController.java:18` — passes the field through from the HTTP body.

### Transformers

- `services/orders/OrderMapper.java:31` — INSERT into `orders` (column is `BYTEA` via `pgcrypto`).

### Sinks

- **Transactional** — `postgres://orders.orders`
  - producer: `OrderService.java:42`
  - transformer: `OrderMapper.java:31` (BYTEA / encrypted)
  - consumer: `InvoiceService.java:14`
  - retention: `transactional`
  - residency: `eu-west-1`
  - control: `audit-logged` (OrderAuditLog) · `encrypted-at-rest` (pgcrypto) · `consent-checked` (boundary)

- **Analytics** — `snowflake://analytics.orders_fact`
  - producer: `daily export job` (cron)
  - transformer: none (column-level encrypted on ingest)
  - consumer: BI dashboard
  - retention: `analytics`
  - residency: `eu-west-1`
  - control: `encrypted-at-rest` (Snowflake column-level) · `consent-checked` (n/a, analytics is post-consent)

- **Log** — `logback://orders.events`
  - producer: `LogRedactor.java:7` (writes `***` instead of the field)
  - transformer: n/a (redacted at write)
  - consumer: log aggregator
  - retention: `log`
  - residency: `eu-west-1`
  - control: `audit-logged` (n/a for logs) · `encrypted-at-rest` (n/a, redacted)

- **External** — `payments.Stripe.paymentIntents.create`
  - producer: `OrderService.java:42` (called after consent interceptor)
  - transformer: none (Stripe stores the field as-is)
  - consumer: Stripe
  - retention: `external`
  - residency: `us-east-1` (Stripe)
  - control: `audit-logged` (Stripe) · `consent-checked` (boundary) · `cross-border` (eu → us, requires DPA)

### Controls

- `audit-logged`: present on transactional and external. Logs and analytics are exempt.
- `encrypted-at-rest`: present on transactional and analytics. Logs are redacted.
- `consent-checked`: present at the boundary. Analytics is exempt (post-consent rows only).

### Deletion Targets (for right-to-erasure)

- `postgres://orders.orders` — covered by `ErasureWorker.deleteOrder`.
- `snowflake://analytics.orders_fact` — covered by `ErasureWorker.purgeAnalytics`.
- `payments.Stripe` — covered by `ErasureWorker.purgeStripeCustomer` (separate Stripe API call).

### Findings

- `LIN-CROSS-BORDER` — the field crosses `eu-west-1` → `us-east-1` (Stripe). The DPA is in place; the cross-border notice is in the consent record.
- (no other findings)

---

When the walk finds nothing, the lineage report emits exactly:

```text
未发现数据血缘漂移。
```
