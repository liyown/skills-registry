# Critical

## 1. Payment endpoint missing order ownership and concurrent-state guard leads to cross-user charge and duplicate deduction

Location:
`example.BadOrderService.Pay`

Problem:
The method looks up the order by `orderID` and charges it without verifying `order.UserID == userID` and without a status-conditional update. Concurrent requests can both read `UNPAID` and both proceed to charge.

Impact:
Likely cross-user payment, duplicate deduction, and order-state corruption. This is a capital-loss risk.

Suggestion:
Look up the order by `(orderID, userID)`. Use a unique payment-flow id or a status-conditional update (`UPDATE ... WHERE status = 'UNPAID'`) for the deduction, and check the affected row count. On any failure, define an explicit compensation or refund path.

Recommended code:
```go
var o Order
err := s.db.QueryRowContext(ctx,
    "SELECT id, user_id, amount, status FROM orders WHERE id = ? AND user_id = ?",
    orderID, userID,
).Scan(&o.ID, &o.UserID, &o.Amount, &o.Status)
if errors.Is(err, sql.ErrNoRows) {
    return ErrOrderNotFound
}
if o.Status == "PAID" {
    return nil
}
res, err := s.db.ExecContext(ctx,
    "UPDATE orders SET status = 'PAYING' WHERE id = ? AND status = 'UNPAID'", orderID,
)
if n, _ := res.RowsAffected(); n != 1 {
    return ErrConcurrentUpdate
}
```

## 2. fmt.Sprintf with integer in SQL string can degrade into injection if the field type weakens

Location:
`example.BadOrderService.Pay`

Problem:
`fmt.Sprintf("SELECT ... WHERE id = %d", orderID)` concatenates the id into the SQL string. Today the field is `int64`, but the type can be weakened later, or the parameter can be changed to a string with user input.

Impact:
An attacker who can control the field can change SQL semantics, leak data, perform destructive writes, or bypass permission checks. This is a SQL-injection vulnerability.

Suggestion:
Always use `?` placeholders with prepared statements, regardless of the column's Go type.

Recommended code:
```go
row := s.db.QueryRowContext(ctx, "SELECT ... WHERE id = ?", orderID)
```

# High

## 1. Goroutine outlives the request and is not cancellable

Location:
`example.BadOrderService.Pay`

Problem:
```go
go func() {
    for m := range msgs {
        handleMessage(context.Background(), m)
    }
}()
```
The goroutine is not bound to a `ctx.Done()` and the handler uses `context.Background()`, so the goroutine continues to run after the request is cancelled.

Impact:
Worker goroutines accumulate linearly with traffic, hold downstream resources, and never release. In a long-running process this is a slow, hard-to-see goroutine leak that exhausts memory and file descriptors.

Suggestion:
Pass the request context into the goroutine, and exit on `ctx.Done()`. Close the input channel from the producer side on shutdown.

Recommended code:
```go
go func() {
    for {
        select {
        case <-ctx.Done():
            return
        case m, ok := <-msgs:
            if !ok {
                return
            }
            handleMessage(ctx, m)
        }
    }
}()
```

## 2. fmt.Errorf with %v breaks errors.Is / errors.As

Location:
`example.BadOrderService.Pay`

Problem:
`fmt.Errorf("order paid: %v", sql.ErrNoRows)` uses `%v`, which formats the value but does not chain the error. `errors.Is(err, sql.ErrNoRows)` always returns false for the wrapped error.

Impact:
A query that should resolve to "not found" or "already paid" is reported as an unknown 5xx, degrading the availability metric and breaking any retry / circuit-breaker logic that depends on error identity.

Suggestion:
Use `%w` to wrap; reserve `%v` for non-error values that you want to render into the message.

Recommended code:
```go
return fmt.Errorf("order paid: %w", sql.ErrNoRows)
```

## 3. http.Client has no Timeout and the request path uses log.Fatal

Location:
`example.BadOrderService.Pay`

Problem:
`http.Client{}` is constructed without `Timeout`, and a slow downstream triggers `log.Fatal`, which calls `os.Exit(1)`.

Impact:
A flaky downstream terminates the entire process, blowing the error radius far beyond the failing call. A single slow dependency can take the whole service down.

Suggestion:
Set a realistic per-request timeout. Return an error from the function and let the caller decide whether to degrade, retry, or trip a circuit breaker. Do not call `log.Fatal` from a request path.

Recommended code:
```go
client := &http.Client{Timeout: 5 * time.Second}
resp, err := client.Do(req)
if err != nil {
    return fmt.Errorf("pay provider call: %w", err)
}
defer resp.Body.Close()
```

# Medium

## 1. Mutex copied by value receiver loses its protection

Location:
`example.Counter.Inc`

Problem:
`Counter` uses a value receiver. Each call copies the receiver, including the embedded `sync.Mutex`. The caller's mutex state is no longer the mutex being locked.

Impact:
Concurrent calls to `Inc` lock different copies of the mutex, so the increments race. `go test -race` will flag this; production data will show a `c.n` smaller than the number of increments.

Suggestion:
Use a pointer receiver so the lock and the field are shared. Update the call site to pass `&counter`.

Recommended code:
```go
func (c *Counter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.n++
}
```

# Low

## 1. Error log carries no request-scoped context

Location:
`example.BadOrderService.Pay`

Problem:
`log.Fatal(err)` (and the surrounding logger calls) do not include the `orderID`, `userID`, or request id.

Impact:
When the service starts failing in production, on-call cannot correlate the log line with a specific user or request, and the post-mortem takes much longer than it should.

Suggestion:
Use a structured logger (`slog` / `zap` / `zerolog`) and pass `order_id`, `user_id`, and `request_id` as fields on every log line in the request path.

Recommended code:
```go
slog.Error("pay failed", "order_id", orderID, "user_id", userID, "err", err)
```
