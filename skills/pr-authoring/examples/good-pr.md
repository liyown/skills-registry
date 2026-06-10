# PR Body — Good Draft

<!-- Good counterpart of bad-pr.md.
     Fix 1 (SECTIONS): all 8 sections present and filled.
     Fix 2 (REVIEWERS): per-team reviewer list with @handle.
     Fix 3 (VERIFICATION): unit / integration / manual / CI / metrics all listed.
     Fix 4 (LINKED-DOCS): issue + RFC + design doc links.
     Fix 5 (RISK-ROLLBACK): risk + rollback + runbook + on-call all listed.
-->

## 1. Summary

本次 PR 修复了 `OrderService.pay` 在并发场景下的双重扣款 bug。

**为什么改：** 5 月 27 日的 incident #INC-2026-514 复盘指出，缺少 `WHERE status = 'UNPAID'` 旧状态条件是双重扣款的根因。本次 PR 加上这个条件，并把 `markPaying` 改成状态条件 update，affected row count = 0 时直接放弃。

**影响范围：** `services/orders/OrderService.java`、`services/orders/OrderMapper.java`、新增 `OrderServiceTest` 的 4 个测试方法。受影响的下游：`services/billing/InvoiceService`（调用 `OrderService.pay`），不在本 PR 范围内但需通知。

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

`@alice` 是必审；`@bob` 是 informational（billing 是下游 consumer，不是本 PR 范围内的 owner）。

## 6. Verification evidence

- **Unit:** `OrderServiceTest.test_markPaying_status_conditional` (passing) — verifies the new `WHERE status = 'UNPAID'` clause.
- **Unit:** `OrderServiceTest.test_markPaying_concurrent_two_calls` (passing) — verifies the concurrent-double-charge fix.
- **Unit:** `OrderServiceTest.test_markPaying_already_paid` (passing) — verifies no-op on already-paid orders.
- **Unit:** `OrderServiceTest.test_markPaying_affected_row_zero` (passing) — verifies the abandonment path.
- **Integration:** `OrderPaymentIT.test_concurrent_pay_via_http` (passing) — verifies the full HTTP path.
- **CI:** [CI run #1234](https://ci.example.com/runs/1234) — all checks green.
- **Manual:** locally ran `make test` + `make it` — all green. Screenshot: <link to screenshot>.

## 7. Linked docs

- **Issue:** #INC-2026-514 (incident report, internal).
- **RFC:** [RFC-007 auth token format migration](https://github.com/liyown/superpower-enterprise/blob/main/docs/rfcs/2026-06-10-enterprise-flow-skills.md) — not directly related to this PR but the postmortem cites the same audit-log gap.
- **Design doc:** [internal: orders service contract](https://internal.example.com/docs/orders-contract) — section 4.2 (payment state machine) updated to reflect the new contract.
- **关联的 PR:** none (this is the fix; the rollback PR, if any, will be linked here).

## 8. Risk + rollback

- **Risk:** The new `WHERE status = 'UNPAID'` clause changes the order state machine; an order that is in `PAYING` (mid-flight) and then re-pays after the test cleanup window will be abandoned. Risk: **low** — the abandonment path is the desired behaviour per the postmortem; the only consumer-affecting change is that a manual re-pay within a 5-minute window is now a no-op.
- **Rollback:** `git revert a1b2c3d` reverts the fix and the docs. The tests are additive (`test_*` methods) and stay even on revert.
- **Runbook:** [runbooks/orders.md#markPaying-rollback](https://internal.example.com/runbooks/orders#markPaying-rollback) — covers the manual flag flip if a partial rollback is needed.
- **On-call:** `@org/orders-core` on-call (`@alice`) — paged if `markPaying.error_rate > 0.5%` post-deploy.
