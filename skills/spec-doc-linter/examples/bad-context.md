# Billing Context

## Bounded Context

The `billing` domain is responsible for invoicing and refunds.
It is the source of truth for an order's monetary state.

## Owns

- Order monetary state (status from a billing perspective:
  `UNPAID`, `PAID`, `REFUNDED`).
- The `Invoice` aggregate.

## Does Not Own

- Shipping. `billing` does not depend on `shipping`.
- Tax computation — that is `pricing`'s job.

## Upstream Contexts

- `iam` (for the caller's `userId`)

## Downstream Contexts

- `payments` (sends charge requests to)
- `pricing` (asks for the final amount)

## Critical Invariants

- An order is `PAID` if and only if `payments` has
  acknowledged a charge.
- A refund never goes below the original charge amount.

See [`./related-domains.md`](./related-domains.md) for the
historical decomposition notes.
