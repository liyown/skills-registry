# Billing Context

<!-- Good counterpart of bad-context.md.
     Fix 1 (B.4): boundary claim rewritten to acknowledge the
       read-only shipping rate lookup.
     Fix 2 (B.6): inventory added to Downstream Contexts.
     Fix 3 (A.6): broken related-domains.md link removed.
-->

## Bounded Context

The `billing` domain is responsible for invoicing and refunds.
It is the source of truth for an order's monetary state.

## Owns

- Order monetary state (status from a billing perspective:
  `UNPAID`, `PAID`, `REFUNDED`).
- The `Invoice` aggregate.

## Does Not Own

- Shipping fulfillment. `billing` reads `shipping` for the
  rate quote at invoice time but does not own the rate.
- Tax computation — that is `pricing`'s job.

## Upstream Contexts

- `iam` (for the caller's `userId`)

## Downstream Contexts

- `payments` (sends charge requests to)
- `pricing` (asks for the final amount)
- `inventory` (reads `Stock` to confirm the order is
  fulfillable before charging)

## Critical Invariants

- An order is `PAID` if and only if `payments` has
  acknowledged a charge.
- A refund never goes below the original charge amount.
