# Patient Audit-Log Control — Good Pack

<!-- Good counterpart of bad-control.md.
     Fix 1 (ENFORCEMENT-SITES): every enforcement site named with path + line + excerpt.
     Fix 2 (BYPASS-LIST): every bypass named.
     Fix 3 (LINEAGE-CROSS-CHECK): which sinks are covered, which are not.
     Fix 4 (EVIDENCE-COLLECTOR): the path the auditor can run to re-verify.
     Fix 5 (PER-SITE-EXCERPT): the smallest code excerpt that proves the interceptor fires.
-->

## Control

HIPAA `164.312(b)` — audit controls. Every write to `Patient` must be audit-logged.

## Control row

- **Regulator frame + control ID:** HIPAA `164.312(b)` (audit controls)
- **Canonical enforcement symbol:** `interceptors/AuditLogInterceptor.audit` (`interceptors/AuditLogInterceptor.java:18`)
- **Per-site evidence:**
  - `services/patients/PatientService.createPatient` (`services/patients/PatientService.java:42`) — `AuditLogInterceptor.audit('patient.create', id)` invoked at the start of the method, before persistence.
  - `services/patients/PatientService.updatePatient` (`services/patients/PatientService.java:71`) — `AuditLogInterceptor.audit('patient.update', id, before, after)` invoked with a before/after diff.
  - `services/patients/PatientService.deletePatient` (`services/patients/PatientService.java:104`) — `AuditLogInterceptor.audit('patient.delete', id)` invoked before the delete, with a `deletion_target` marker for the erasure worker.
- **Per-bypass findings:**
  - **BYPASS-1** `scripts/maintenance/recompute-billing.sql:18` — direct `UPDATE patients SET billing_address = ...` for the billing address recompute. The script does not go through `PatientService` and so does not invoke the interceptor. **Suggested fix:** route the recompute through `PatientService.recomputeBilling` (a new method that invokes the interceptor), or add a one-time script that re-walks the affected rows and emits audit entries.
  - **BYPASS-2** `reports/PatientExportJob.exportDailySnapshot` (`reports/PatientExportJob.java:42`) — reads from `patients` but does not write. **Not a bypass** (audit is for writes; reads are out of scope for `164.312(b)`).
- **Data-lineage cross-check:**
  - `postgres://patients.patients` — covered (3 enforcement sites above)
  - `snowflake://analytics.patients_fact` — **NOT COVERED**. The daily export job writes the column-level snapshot to the warehouse without invoking the interceptor. **Suggested fix:** the export job should invoke `AuditLogInterceptor.audit('patient.export', id, snapshot)` before the warehouse write.
  - `backups/s3://patients-daily` — **NOT COVERED** (the backup is a disk-level copy, not a write through the interceptor). The HIPAA control does not require per-row audit on the backup; the control is the in-band write audit.
- **Evidence-collector:**
  - Run `bin/audit-log query --entity=patient --since=2026-06-01 --until=2026-06-30` to enumerate every audit entry for the period. The count must equal the number of writes (from the DB metrics) within 0.1%.
  - For the bypass: run `bin/audit-log query --action=patient.recompute_billing --since=2026-06-01` and confirm zero entries; the absence of entries is the bypass evidence.
  - For the export: run `bin/audit-log query --action=patient.export --since=2026-06-01` and confirm zero entries; the absence of entries is the warehouse-sink gap.

## Per-bypass findings

- `BYPASS-1` — `scripts/maintenance/recompute-billing.sql:18` — direct UPDATE bypass. Severity: High.
- `BYPASS-2` — `reports/PatientExportJob.exportDailySnapshot` warehouse write — not covered by the in-band interceptor. Severity: Medium.

## Verdict

- HIPAA `164.312(b)` audit control: **PARTIALLY COVERED** (in-band write paths covered; one maintenance script and one warehouse export uncovered).
