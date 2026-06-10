# Critical

## 1. Payment endpoint missing order ownership and idempotency guard leads to cross-user charge and duplicate deduction

Location:
`BadOrderService#doPay`

Problem:
The method looks up the order by `orderId` only — it does not check that the order belongs to the current `userId`, and the `PAID` check followed by the deduction and status update is not protected by a state-conditional update or a unique payment-flow id. Concurrent requests can both read `UNPAID` and both proceed to charge.

Impact:
Likely cross-user payment, duplicate deduction, and order-state corruption. This is a capital-loss risk.

Suggestion:
Look up the order by `(orderId, userId)`. Use a unique payment-flow id or a status-conditional update (`UPDATE ... WHERE status = 'UNPAID'`) for the deduction, and check the affected row count. On any failure, define an explicit compensation or refund path.

Recommended code:
```java
Order order = orderMapper.selectByIdAndUserId(orderId, userId);
if (order == null) {
    throw new BizException("order not found");
}
if ("PAID".equals(order.getStatus())) {
    return;
}

int updated = orderMapper.markPaying(orderId, userId, "UNPAID", "PAYING");
if (updated != 1) {
    return;
}
```

## 2. MyBatis-Plus wrapper.last concatenation with user-controlled sort causes SQL injection

Location:
`BadOrderService#search`

Problem:
`wrapper.last("order by " + sort)` concatenates the user-supplied `sort` argument directly into the SQL fragment. `last()` does not bind parameters; whatever string the caller supplies is appended verbatim.

Impact:
An attacker can craft a sort parameter that changes SQL semantics, leaks data, raises errors, or — depending on DB user and connection rights — executes destructive statements. This is a SQL-injection vulnerability.

Suggestion:
Use a whitelist map from sort key to column name. Do not let arbitrary strings reach the SQL fragment.

Recommended code:
```java
Map<String, String> sortColumns = Map.of(
    "createdTime", "created_time",
    "amount", "amount"
);
String column = sortColumns.getOrDefault(sort, "created_time");
wrapper.orderByDesc(column);
```

# High

## 1. Self-invocation skips the Spring AOP proxy and silently disables @Transactional

Location:
`BadOrderService#pay`

Problem:
`pay()` calls `this.doPay()` directly. Spring AOP proxies do not intercept self-invocations, so the `@Transactional` on `doPay()` is silently bypassed.

Impact:
Deduction, remote payment, and order-state update run outside any transaction boundary. On failure, the user can be charged but the order stays in `UNPAID` — the canonical "money taken, order not updated" intermediate state that requires manual reconciliation.

Suggestion:
Move the transactional method to a separate bean and call it through the Spring proxy, or carry the transaction boundary on the outer public method.

Recommended code:
```java
// In a separate @Service bean
@Service
public class OrderTransactionalOps {
    @Transactional(rollbackFor = Exception.class)
    public void doPay(Order order) { ... }
}

// Caller routes through the proxy
orderTransactionalOps.doPay(order);
```

# Medium

## 1. Order and account objects dereferenced without null check

Location:
`BadOrderService#doPay`

Problem:
`orderMapper.selectById(orderId)` and `accountMapper.selectByUserId(userId)` may both return null. The next lines access `order.getStatus()`, `order.getAmount()`, and `account.getPayToken()` unconditionally.

Impact:
An invalid order id, a soft-deleted account, or dirty data throws NPE, surfacing as HTTP 500 and breaking the user flow.

Suggestion:
Add explicit null checks on the query results and return a typed business error (`order not found`, `account not found`) instead of letting the NPE propagate.

Recommended code:
```java
Order order = orderMapper.selectById(orderId);
if (order == null) {
    throw new BizException("order not found");
}
Account account = accountMapper.selectByUserId(userId);
if (account == null) {
    throw new BizException("account not found");
}
```

# Low

## 1. Per-row query and update inside batch cancel

Location:
`BadOrderService#batchCancel`

Problem:
The batch cancel iterates the order ids and runs one `SELECT` and one `UPDATE` per order.

Impact:
The batch's wall-clock time grows linearly with the number of ids, drags the endpoint's p99 up, and increases connection-pool pressure under load.

Suggestion:
Batch the lookup and update into a single `SELECT ... WHERE id IN (?, ?, ...)` and a single `UPDATE ... WHERE id IN (?, ?, ...)`, or cap the batch size and chunk the work.

Recommended code:
```java
List<Order> orders = orderMapper.selectByIds(ids);
List<Long> cancellable = orders.stream()
    .filter(o -> "UNPAID".equals(o.getStatus()))
    .map(Order::getId)
    .toList();
if (!cancellable.isEmpty()) {
    orderMapper.batchCancel(cancellable);
}
```
