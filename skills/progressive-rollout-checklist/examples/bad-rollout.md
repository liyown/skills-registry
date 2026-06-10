# Pricing Engine Rollout — Bad Packet

## Change

Flipping the default of `use_new_pricing_engine` from `false` to `true`.

## The Packet (as the author would otherwise send it)

> The flag is in place; the test passes; ready to ship.

## What's wrong with this packet

- **No read sites listed.** "The flag is in place" is a claim, not a list. A flag with one read site is wired up; a flag with ten read sites where one is missed is silently off in production.
- **No write sites listed.** A flag can be flipped by an admin endpoint, a cohort expander, or a kill switch. If the author does not name every writer, the flip can be triggered by an unexpected path.
- **No cohort definition.** A flag that flips globally with no cohort is a flag without a rollout — it is a big-bang change, which is exactly what the flag mechanism is supposed to prevent.
- **No SLO gate.** A flag with no metric-based gate is a flag that will be expanded by gut feel.
- **No kill switch.** Without a kill switch, the only way to stop the change is to roll back the deploy, which is slow and lossy.
- **No rollback hook.** The flip is reversible by flipping the flag back, but the underlying data (the new price computations stored in the DB) may not be.
