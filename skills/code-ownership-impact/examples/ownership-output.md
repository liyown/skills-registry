# Cross-Team Blast Radius — Review Packet

## Diff

`common-utils/src/main/java/com/example/common/OrderId.java` (renamed + type changed)
`services/orders/OrderService.java`
`services/billing/InvoiceService.java`
`services/payments/PaymentService.java`

## 1. Single-team ack

- **Team:** `@org/orders-core`
- **Owned paths:** `services/orders/OrderService.java`
- **Named approver:** `@alice`, `@bob`
- **Deadline:** EOD today
- **Cross-team impact:** none
- **Ack-channel:** thread in `#orders-core-review`

## 2. Cross-team approvers

- **Team:** `@org/billing-core`
  - **Owned paths:** `services/billing/InvoiceService.java`
  - **Named approver:** `@carol`
  - **Deadline:** Friday 17:00 (next billing-train)
  - **Cross-team impact:** consumes `OrderId` via the orders events stream
  - **Ack-channel:** thread in `#billing-core-review`

- **Team:** `@org/payments-core`
  - **Owned paths:** `services/payments/PaymentService.java`
  - **Named approver:** `@dave`
  - **Deadline:** Wednesday 17:00 (next payments-train)
  - **Cross-team impact:** consumes `OrderId` (passes to Stripe)
  - **Ack-channel:** thread in `#payments-core-review`

## 3. Shared-kernel subscribers (`common-utils`)

- `services/fulfillment/FulfillmentService.java` — `@org/fulfillment-core` / `@eve` / Monday 10:00 / CODEOWNERS `@` ping + thread
- `services/notifications/NotificationService.java` — `@org/notifications-core` / `@frank` / Tuesday 14:00 / CODEOWNERS `@` ping
- `services/reports/ReportingService.java` — `@org/reports-core` / `@grace` / ASAP (hotfix) / CODEOWNERS `@` ping + thread

## 4. Suggested review order

1. `@org/orders-core` (author's team)
2. `@org/reports-core` (hotfix)
3. `@org/payments-core` (Wed 17:00)
4. `@org/billing-core` (Fri 17:00)
5. `@org/notifications-core` (Tue 14:00)
6. `@org/fulfillment-core` (Mon 10:00)
