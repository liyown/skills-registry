# Cross-Team Blast Radius — Good Packet

<!-- Good counterpart of bad-ownership.md.
     Fix 1 (TEAM-ATTR): per-team attribution with CODEOWNERS-resolved approvers.
     Fix 2 (KERNEL-BLAST): shared-kernel subscriber list with each consumer's
       owning team.
     Fix 3 (DEADLINE): per-team review deadline tied to the team's release train.
     Fix 4 (ACK-CHANNEL): per-team ack-channel rather than a single broadcast.
-->

## Diff

`git diff --name-only main..HEAD`:

```
common-utils/src/main/java/com/example/common/OrderId.java
services/orders/OrderService.java
services/billing/InvoiceService.java
services/payments/PaymentService.java
```

## The Packet

### 1. Single-team ack (author's team)

- **Team:** `@org/orders-core`
- **Owned paths in the diff:** `services/orders/OrderService.java`
- **Named approver:** `@alice`, `@bob` (per `CODEOWNERS`)
- **Suggested review deadline:** EOD today (the author's team has a daily-review window).
- **Cross-team impact:** none for this team's owned paths.
- **Ack-channel:** thread in `#orders-core-review`.

### 2. Cross-team approvers (per `CODEOWNERS`)

- **Team:** `@org/billing-core`
  - **Owned paths:** `services/billing/InvoiceService.java`
  - **Named approver:** `@carol`
  - **Suggested review deadline:** next billing-train slot — Friday 17:00 local.
  - **Cross-team impact:** consumes `OrderId` (via the orders events stream).
  - **Ack-channel:** thread in `#billing-core-review`.

- **Team:** `@org/payments-core`
  - **Owned paths:** `services/payments/PaymentService.java`
  - **Named approver:** `@dave`
  - **Suggested review deadline:** next payments-train slot — Wednesday 17:00 local.
  - **Cross-team impact:** consumes `OrderId` (passes to Stripe).
  - **Ack-channel:** thread in `#payments-core-review`.

### 3. Shared-kernel subscribers (for `common-utils`)

- **Service:** `services/fulfillment/FulfillmentService.java`
  - **Owning team:** `@org/fulfillment-core`
  - **Named approver:** `@eve`
  - **Suggested review deadline:** next fulfillment-train slot — Monday 10:00 local.
  - **Cross-team impact:** direct import of `OrderId`. Type change breaks the call site.
  - **Ack-channel:** CODEOWNERS `@` ping + thread in `#fulfillment-core-review`.

- **Service:** `services/notifications/NotificationService.java`
  - **Owning team:** `@org/notifications-core`
  - **Named approver:** `@frank`
  - **Suggested review deadline:** next notifications-batch window — Tuesday 14:00 local.
  - **Cross-team impact:** direct import of `OrderId`. Type change breaks the call site.
  - **Ack-channel:** CODEOWNERS `@` ping.

- **Service:** `services/reports/ReportingService.java`
  - **Owning team:** `@org/reports-core`
  - **Named approver:** `@grace`
  - **Suggested review deadline:** ASAP (the report is in a hot-fix window).
  - **Cross-team impact:** direct import of `OrderId`. Type change breaks the report's join.
  - **Ack-channel:** CODEOWNERS `@` ping + thread in `#reports-hotfix`.

### 4. Suggested review order

1. `@org/orders-core` (author's team) — first.
2. `@org/reports-core` (hotfix window) — second.
3. `@org/payments-core` — Wednesday 17:00.
4. `@org/billing-core` — Friday 17:00.
5. `@org/notifications-core` — Tuesday 14:00.
6. `@org/fulfillment-core` — Monday 10:00.
