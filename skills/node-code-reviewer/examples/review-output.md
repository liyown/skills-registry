# Critical

## 1. Object.assign + JSON.parse allows prototype pollution

Location:
`mergeConfig`

Problem:
```ts
return Object.assign(defaults, JSON.parse(input));
```
An attacker can pass `{"__proto__":{"isAdmin":true}}`. `Object.assign` writes the `__proto__` key into the target, polluting its prototype.

Impact:
Every object in the process inherits `isAdmin = true`. Any later `if (user.isAdmin)` check passes for every user, including unauthenticated ones. This is an authentication-bypass vulnerability.

Suggestion:
Reject `__proto__`, `constructor`, and `prototype` keys before merging. Use `Object.create(null)` for the target, or define properties with `Object.defineProperty` so they do not walk the prototype chain.

Recommended code:
```ts
function safeMerge<T extends object>(target: T, patch: Record<string, unknown>): T {
  for (const k of Object.keys(patch)) {
    if (k === '__proto__' || k === 'constructor' || k === 'prototype') continue;
    (target as any)[k] = patch[k];
  }
  return target;
}
```

## 2. fs.readFileSync in the request path blocks the event loop

Location:
`OrderService.loadConfig`

Problem:
`fs.readFileSync` is a synchronous call. On the event loop, it blocks the worker for the duration of the read — including disk wait.

Impact:
Other requests on the same worker are stalled. Under any meaningful concurrency, p99 climbs sharply and the worker appears "frozen" for the duration of the file read.

Suggestion:
Use the async file API. For large files, stream rather than read all at once.

Recommended code:
```ts
async loadConfig(): Promise<Config> {
  const raw = await fs.promises.readFile(CONFIG_PATH, 'utf8');
  return JSON.parse(raw);
}
```

## 3. Status update has no WHERE-clause guard and uses $executeRawUnsafe

Location:
`OrderService.pay`

Problem:
```ts
prisma.$executeRawUnsafe(
  `UPDATE orders SET status = 'PAID' WHERE id = ${orderId}`,
)
```
Two issues: the SQL is built with template-string interpolation, and the update has no `AND status = 'UNPAID'` guard. `$executeRawUnsafe` does not bind parameters.

Impact:
Concurrent requests can both pass the `PAID` check and both call the update — duplicate charge. If `orderId` is ever derived from user input, this is also a SQL-injection vulnerability.

Suggestion:
Use `$executeRaw` (the safe variant) with positional parameters and add a status-conditional update. Check the returned row count.

Recommended code:
```ts
const updated = await prisma.$executeRaw(
  sql`UPDATE orders SET status = 'PAID' WHERE id = ${orderId} AND status = 'UNPAID'`,
);
if (updated !== 1) {
  return; // already paid or order vanished
}
```

# High

## 1. Lost await turns an exception into an unhandled rejection

Location:
`backgroundJob`

Problem:
```ts
doWork(payload);
```
The call is not awaited. If `doWork` throws, the rejection is unhandled; depending on the Node version and `unhandledRejection` configuration, the process can crash.

Impact:
Background work fails silently from the caller's point of view. Under strict `unhandledRejection` policy (the default in modern Node), the process exits and takes the entire service down.

Suggestion:
`await` the call inside an async handler, or chain `.catch` to log + record the failure. Register a process-level handler as a last-resort safety net.

Recommended code:
```ts
await doWork(payload);
```

## 2. Missing return after next(err) in try/catch causes double response

Location:
`/pay` route

Problem:
```ts
try { await payService.charge(req.body); } catch (e) { next(e); }
res.json({ ok: true });
```
After `next(e)` the function does not `return`, so `res.json({ ok: true })` still runs.

Impact:
Express logs "Cannot set headers after they are sent". The response state is undefined; in the worst case, the user sees a 200 success while the charge actually failed, which can lead to a retry and a double-charge.

Suggestion:
Add `return` before `next(e)`, and likewise `return res.json(...)` (or restructure with early returns).

Recommended code:
```ts
try {
  await payService.charge(req.body);
  return res.json({ ok: true });
} catch (e) {
  return next(e);
}
```

# Medium

## 1. Module-level mutable Map shared across requests and workers

Location:
`HITS`

Problem:
A module-level `Map<number, number[]>` is loaded once per process. Under cluster mode, each worker has its own copy, but inside one worker every request shares the same map.

Impact:
The rate limit is per-process, not global. Memory grows linearly with traffic. The "rate limit" is not enforcing the intended policy.

Suggestion:
Move the counter to a shared store (Redis sorted-set with `ZADD` / `ZREMRANGEBYSCORE`, or a leaky-bucket in a shared KV). At minimum, cap the map size with an LRU policy and a per-key TTL.

Recommended code:
```ts
// Use ioredis sliding window
async function isAllowed(userId: number): Promise<boolean> {
  const now = Date.now();
  const key = `hits:${userId}`;
  const pipe = redis.pipeline();
  pipe.zremrangebyscore(key, 0, now - 60_000);
  pipe.zadd(key, now, String(now));
  pipe.zcard(key);
  pipe.expire(key, 60);
  const [, , count] = await pipe.exec();
  return count <= 100;
}
```

# Low

## 1. console.log left in production code

Location:
`OrderService.pay`

Problem:
`console.log("paid", orderId)` is left in the request path. The arguments include the raw `orderId`, which may be sensitive depending on what an order id encodes.

Impact:
Operational noise in stdout; potential PII / order-correlation leakage into log-aggregation systems that are not access-controlled for production data.

Suggestion:
Use a structured logger with explicit field allow-listing, and remove the debug log before shipping.

Recommended code:
```ts
logger.info({ order_id: orderId }, 'order paid');
```
