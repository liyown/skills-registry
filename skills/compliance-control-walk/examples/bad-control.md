# Patient Audit-Log Control — Bad Pack

## Control

HIPAA `164.312(b)` — audit controls. Every write to `Patient` must be audit-logged.

## The Pack (as the auditor would otherwise receive it)

> We audit all writes to Patient via the global audit interceptor.

## What's wrong with this pack

- **No enforcement site list.** "The global audit interceptor" is one symbol, but the auditor needs every path + line that invokes it. A single missing invocation is a compliance finding.
- **No bypass list.** A direct `UPDATE patients SET ...` from a maintenance script is a bypass that the global interceptor does not cover. The pack does not name it.
- **No data-lineage cross-check.** The pack does not say which sinks the audit interceptor covers. If a secondary sink (a warehouse, an export) is not covered, the auditor's question is open.
- **No evidence-collector path.** The auditor cannot re-verify the claim; there is no log-line id, no metrics query, no test command.
- **No per-site code excerpt.** "Audit interceptor" is the name; the auditor needs the smallest code excerpt that proves the interceptor is invoked.
