# Drift Report

## src/main/java/com/example/order/DevAgent.md

- L13  A.1  doc references `BadOrderService.batchCancel`; not found under `src/main/java/com/example/order/`  →  rename to `BadOrderService.cancelMany`.
- L11  B.1  doc claims `batchCancel` is idempotent; code mutates state per id  →  drop the idempotency claim.
- L19  B.3  doc claims `batchCancel` is atomic; code runs one UPDATE per id  →  drop the atomicity claim.
- L28  A.5  manifest lists `sentry`; doc Dependencies does not mention it  →  add a bullet under Dependencies.
- L30  A.5  doc lists `kafka`; manifest no longer has `kafka`  →  remove the `kafka` bullet.
- L41  A.4  doc references `./docs/knowledge/order-history.md`; file does not exist  →  remove the reference.

## domains/billing/CONTEXT.md

- L8   B.4  doc says `billing does not depend on shipping`; `billing/` imports `shipping.RateQuote`  →  rewrite the boundary line to acknowledge the read-only rate lookup.
- L22  B.6  doc Downstream Contexts does not list `inventory`; `billing/` imports `inventory/Stock`  →  add `inventory` to Downstream Contexts.
- L30  A.6  link `[related-domains.md](./related-domains.md)` does not resolve  →  remove the link.

# Needs Context

(none)

---

Proposed edit for `src/main/java/com/example/order/DevAgent.md`:

```diff
--- src/main/java/com/example/order/DevAgent.md
+++ src/main/java/com/example/order/DevAgent.md
@@ -10,9 +10,9 @@
 - `BadOrderService.pay(userId: Long, orderId: Long): Unit`
   - Reads the order, checks the status, charges the user, and marks paid.
-- `BadOrderService.batchCancel(ids: List<Long>): Unit`
-  - Cancels a list of orders in one call. **Idempotent** — calling
-    twice with the same ids is a no-op.
+- `BadOrderService.cancelMany(ids: List<Long>): Int`
+  - Returns the number of orders that transitioned from
+    `UNPAID` to `CANCELLED`. Not idempotent; not atomic.
@@ -27,7 +27,7 @@
 - `payments` (for the charge call)
-- `kafka` (for emitting `OrderPaid` events)
+- `sentry` (for emitting `OrderFailed` events)
@@ -38,5 +38,3 @@
 - `./src/main/java/com/example/order/PaymentClient.java`
-- `./docs/knowledge/order-history.md` — historical context.
--
```

3 change(s) proposed in src/main/java/com/example/order/DevAgent.md (3 Class A, 0 Class B)
Apply? [y/n/q]:

---

When no drift is found, the linter emits exactly:

```text
未发现文档与代码漂移。
```
