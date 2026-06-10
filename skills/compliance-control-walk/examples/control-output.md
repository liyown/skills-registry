# Patient Audit-Log Control — Evidence Pack

## Control rows

### HIPAA 164.312(b) — audit controls

| Cell | Value |
|---|---|
| Regulator frame + control ID | HIPAA `164.312(b)` |
| Canonical enforcement symbol | `interceptors/AuditLogInterceptor.audit` (`interceptors/AuditLogInterceptor.java:18`) |
| Per-site evidence | `PatientService.java:42` (create), `:71` (update), `:104` (delete) — interceptor invoked at the start of each method |
| Per-bypass findings | `BYPASS-1` (recompute-billing.sql:18), `BYPASS-2` (PatientExportJob.exportDailySnapshot) |
| Data-lineage cross-check | `postgres://patients.patients` covered; `snowflake://analytics.patients_fact` NOT covered; backups exempt |
| Evidence-collector | `bin/audit-log query --entity=patient --since=...` |

## Per-bypass findings

- **`BYPASS-1`** `scripts/maintenance/recompute-billing.sql:18` — direct `UPDATE patients SET billing_address = ...`. Severity: **High**. Suggested fix: route through `PatientService.recomputeBilling` (new method that invokes the interceptor), or add a one-time re-walk script that emits audit entries.
- **`BYPASS-2`** `reports/PatientExportJob.exportDailySnapshot` (`reports/PatientExportJob.java:42`) — warehouse write of the column-level snapshot. Severity: **Medium**. Suggested fix: invoke `AuditLogInterceptor.audit('patient.export', id, snapshot)` before the warehouse write.

## Data-lineage cross-check

- `postgres://patients.patients` — covered (3 enforcement sites)
- `snowflake://analytics.patients_fact` — NOT covered (BYPASS-2)
- `backups/s3://patients-daily` — exempt (disk-level copy, not in-band)

## Evidence-collector list

- `bin/audit-log query --entity=patient --since=2026-06-01 --until=2026-06-30` — enumerates every audit entry; count must match write count within 0.1%.
- `bin/audit-log query --action=patient.recompute_billing --since=2026-06-01` — zero entries expected; absence confirms BYPASS-1.
- `bin/audit-log query --action=patient.export --since=2026-06-01` — zero entries expected; absence confirms BYPASS-2.

## Verdict

HIPAA `164.312(b)` audit control: **PARTIALLY COVERED** (in-band write paths covered; one maintenance script and one warehouse export uncovered).
