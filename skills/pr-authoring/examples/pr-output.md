# PR Body — Output

## 1. Summary

本次 PR 修复了 `OrderService.pay` 在并发场景下的双重扣款 bug。

**为什么改：** 5 月 27 日的 incident #INC-2026-514 复盘指出，缺少 `WHERE status = 'UNPAID'` 旧状态条件是双重扣款的根因。

**影响范围：** `services/orders/OrderService.java`、`OrderMapper.java`、新增 4 个测试。受影响的下游：`services/billing/InvoiceService`（informational，不在本 PR 范围内）。

## 2. Type of change

- [x] Bug fix (closes #INC-2026-514)
- [x] Tests
- [ ] Documentation
- [ ] CI / script
- [ ] Refactor

## 3. Commits

```
a1b2c3d fix(orders): add status-conditional update to markPaying
b2c3d4 test(orders): add 4 tests for the concurrent-pay fix
c3d4e5 docs(orders): update DevAgent.md to reflect the new contract
```

## 4. Files changed

```
services/orders/OrderService.java     | 12 +++++++++---
services/orders/OrderMapper.java      |  6 +++---
services/orders/OrderServiceTest.java | 84 +++++++++++++++++++++++++++++++
skills/orders/DevAgent.md             |  8 ++++---
4 files changed, 95 insertions(+), 15 deletions(-)
```

关键文件：
- `services/orders/OrderService.java:42` — `markPaying` 增加 status-conditional update
- `services/orders/OrderMapper.java:18` — 新增 `updateStatus` SQL
- `services/orders/OrderServiceTest.java` — 4 个新测试

## 5. Reviewers

| Team | Reviewer | Path |
|---|---|---|
| `@org/orders-core` | `@alice` | `services/orders/OrderService.java`, `OrderMapper.java`, `OrderServiceTest.java` |
| `@org/billing-core` | `@bob` (informational) | downstream consumer; not blocking |

## 6. Verification evidence

- **Unit:** `OrderServiceTest.test_markPaying_status_conditional` (passing)
- **Unit:** `OrderServiceTest.test_markPaying_concurrent_two_calls` (passing)
- **Unit:** `OrderServiceTest.test_markPaying_already_paid` (passing)
- **Unit:** `OrderServiceTest.test_markPaying_affected_row_zero` (passing)
- **Integration:** `OrderPaymentIT.test_concurrent_pay_via_http` (passing)
- **CI:** [CI run #1234](https://ci.example.com/runs/1234) — all checks green
- **Manual:** locally ran `make test` + `make it` — all green

## 7. Linked docs

- **Issue:** #INC-2026-514
- **RFC:** [RFC-007 auth token format migration](https://github.com/liyown/superpower-enterprise/blob/main/docs/rfcs/2026-06-10-enterprise-flow-skills.md)
- **Design doc:** [internal: orders service contract](https://internal.example.com/docs/orders-contract) §4.2
- **关联的 PR:** none

## 8. Risk + rollback

- **Risk:** low — the abandonment path is the desired behaviour; only consumer-affecting change is a manual re-pay within a 5-minute window is a no-op.
- **Rollback:** `git revert a1b2c3d` reverts the fix and docs. Tests are additive.
- **Runbook:** [runbooks/orders.md#markPaying-rollback](https://internal.example.com/runbooks/orders#markPaying-rollback)
- **On-call:** `@alice` (`@org/orders-core`) — paged if `markPaying.error_rate > 0.5%` post-deploy
