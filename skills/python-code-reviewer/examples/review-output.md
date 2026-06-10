# Critical

## 1. Module-level mutable dict is cross-request, cross-worker shared state

Location:
`OrderService.pay`

Problem:
```python
HITS = {}
...
HITS.setdefault(user_id, []).append(time.time())
```
A module-level `dict` is shared across every request and every worker. Under ASGI, requests interleave their `await` points, so writes to `HITS` race without any lock.

Impact:
Rate limiting silently fails. Multi-worker deployments diverge in their per-user counter state, allowing bursty abuse. The dict grows without bound for the lifetime of the process.

Suggestion:
Move the counter to Redis with an atomic sliding-window (`ZADD` + `ZREMRANGEBYSCORE`), or — at minimum — scope it to a single process with an LRU cap and a per-key TTL.

Recommended code:
```python
# Redis-backed sliding window
async def is_allowed(user_id: str) -> bool:
    now = time.time()
    key = f"hits:{user_id}"
    pipe = redis.pipeline()
    pipe.zremrangebyscore(key, 0, now - 60)
    pipe.zadd(key, {str(now): now})
    pipe.zcard(key)
    pipe.expire(key, 60)
    _, _, count, _ = await pipe.execute()
    return count <= 100
```

## 2. f-string interpolated into raw SQL causes injection

Location:
`OrderService.pay`

Problem:
```python
text(f"SELECT * FROM orders WHERE id = {order_id} AND user_id = {user_id}")
```
`text()` does not bind parameters that come from string interpolation. The values are pasted into the SQL string at format time.

Impact:
An attacker who controls `order_id` or `user_id` (e.g. from a path or query parameter) can change SQL semantics, bypass ownership checks, or dump the entire `orders` table.

Suggestion:
Use bound parameters: a parameterised SQL string with `.params(...)`.

Recommended code:
```python
stmt = text("SELECT * FROM orders WHERE id = :id AND user_id = :uid")
result = await session.execute(stmt, {"id": order_id, "uid": user_id})
```

## 3. Synchronous requests inside an async handler blocks the event loop

Location:
`create_order_async`

Problem:
`requests.post(...)` is a blocking call. Calling it from an `async def` handler suspends the event loop worker for the duration of the HTTP round-trip.

Impact:
Every other request served by the same worker is delayed for as long as the downstream call takes. Under load, the service's p99 climbs and eventually every worker is stuck, which is the canonical "asyncio but synchronous I/O" outage.

Suggestion:
Use an async HTTP client. If a sync client is unavoidable, run it through `loop.run_in_executor` so the event loop is released.

Recommended code:
```python
async def create_order_async(payload: dict) -> dict:
    async with httpx.AsyncClient(timeout=5.0) as client:
        resp = await client.post(PAYMENT_URL, json=payload)
        resp.raise_for_status()
        return resp.json()
```

## 4. Status transition has no WHERE-clause guard

Location:
`OrderService.pay`

Problem:
```python
text(f"UPDATE orders SET status = 'PAID' WHERE id = {order_id}")
```
The update does not include `AND status = 'UNPAID'`. Two concurrent requests can both pass the `PAID` check and both call this update.

Impact:
Duplicate deduction and order-state corruption. This is a capital-loss risk.

Suggestion:
Add a status-conditional update and check the affected row count. If 0 rows were updated, the order was already paid by a concurrent caller — return without charging.

Recommended code:
```python
stmt = text("UPDATE orders SET status = 'PAID' WHERE id = :id AND status = 'UNPAID'")
result = await session.execute(stmt, {"id": order_id})
if result.rowcount != 1:
    return  # already paid or does not exist
```

## 5. pickle.loads on untrusted input → arbitrary code execution

Location:
`load_blob`

Problem:
`pickle.loads(blob)` runs arbitrary code at deserialisation time.

Impact:
An attacker who can submit a blob (file upload, message queue, Redis value) executes arbitrary commands as the service user. This is total takeover, not data corruption.

Suggestion:
Replace `pickle` with a data-only format (JSON, MessagePack, protobuf). If you must accept Python objects, use a signed pickle from a trusted source only.

Recommended code:
```python
import json

def load_blob(blob: bytes) -> dict:
    return json.loads(blob)
```

# High

## 1. Fire-and-forget asyncio task loses exceptions and the result

Location:
`fan_out`

Problem:
`asyncio.create_task(coro())` is called without holding the task reference and without `await`. If `coro` raises, the exception is never observed; if the process restarts, the in-flight task is gone.

Impact:
Silent failures: monitoring never sees the error, the user never gets the result, and on graceful shutdown the work is lost.

Suggestion:
Hold the task reference, `await` (or schedule with `asyncio.TaskGroup` in 3.11+), and register a `done_callback` to log any exception that escapes.

Recommended code:
```python
tasks = [asyncio.create_task(coro(arg)) for arg in args]
for t in tasks:
    t.add_done_callback(_log_task_exception)
await asyncio.gather(*tasks)
```

## 2. except Exception returning False hides real failures as "not executed"

Location:
`charge_card`

Problem:
```python
try:
    charge_card(order)
except Exception:
    return False
```
The handler catches everything and returns `False`. The caller cannot distinguish "card was charged but post-processing failed" from "charge was never attempted".

Impact:
Business treats the order as paid while the user was never charged (silent loss) or vice versa (silent double-charge on retry).

Suggestion:
Catch a narrow exception type, log with full context, and re-raise as a domain error that the caller can map to a clear response.

Recommended code:
```python
try:
    charge_card(order)
except PaymentError as e:
    logger.exception("charge failed", extra={"order_id": order.id})
    raise OrderProcessingError("charge failed") from e
```

# Medium

## 1. N+1 in find_orders

Location:
`find_orders`

Problem:
For each id in the input list, the function issues a separate `SELECT` against the database.

Impact:
A list of N ids becomes N round-trips. Database QPS scales with the number of ids per call, latency scales with the same factor, and connection-pool pressure spikes under load.

Suggestion:
Batch the lookup with a single `IN` clause. For SQLAlchemy 2.x async, `session.execute(select(Order).where(Order.id.in_(ids)))` is sufficient.

Recommended code:
```python
async def find_orders(ids: list[int]) -> list[Order]:
    if not ids:
        return []
    stmt = select(Order).where(Order.id.in_(ids))
    result = await session.execute(stmt)
    return list(result.scalars())
```

# Low

## 1. Time-format string uses local timezone implicitly

Location:
`OrderService.pay`

Problem:
`now.strftime("%Y-%m-%d %H:%M:%S")` uses the local timezone of the Python process, which differs between containers and dev laptops.

Impact:
Timestamps in logs and DB rows are inconsistent across environments, which makes incident timelines and audit logs unreliable.

Suggestion:
Render timestamps in UTC explicitly. Use ISO 8601 (`now.isoformat()`) for portability.

Recommended code:
```python
from datetime import datetime, timezone
now = datetime.now(timezone.utc).isoformat()
```
