# Cross-Team Blast Radius — Bad Packet

## Diff

`git diff --name-only main..HEAD`:

```
common-utils/src/main/java/com/example/common/OrderId.java
services/orders/OrderService.java
services/billing/InvoiceService.java
services/payments/PaymentService.java
```

## The Packet (as the author would otherwise send it)

> Hi all, here's the PR for the `OrderId` refactor. Please review when you can. Tagging the orders team.

## What's wrong with this packet

- **No team attribution.** "Tagging the orders team" is ambiguous — is it the author's team, or the order-fulfillment team, or the order-pricing team? `CODEOWNERS` resolves this in one lookup.
- **No blast-radius for `common-utils`.** `common-utils` is a shared kernel; renaming `OrderId` breaks every consumer. The packet does not name the consumers.
- **No per-team deadlines.** The packet says "when you can", which is the same as no deadline. The teams' release-train windows are different.
- **No ack-channel per team.** "Hi all" is one broadcast; the per-team ack needs a different channel.
- **No shared-kernel subscriber list.** The diff touches `common-utils`; every team that imports from it is a downstream approver. None are named.
