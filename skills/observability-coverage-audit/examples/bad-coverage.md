# Payments Coverage Audit — Bad Matrix

## In-scope path

`services/payments/PaymentService` — every `throw` site in the charge / refund / payout flow.

## The Matrix (as the engineer would otherwise assemble it)

> We're well-covered. The path has integration tests.

## What's wrong with this matrix

- **No per-site rows.** "Well-covered" is a claim, not a list. A path with 12 `throw` sites and 1 integration test is "covered" by the integration test only if the test exercises every throw.
- **No alert rows.** A path with no alert is silent when it fails; the on-call is not paged.
- **No runbook rows.** A path with no runbook leaves the on-call guessing on the first 10 minutes of an incident.
- **No per-classification findings.** The four findings (`COV-SILENT`, `COV-NO-RUNBOOK`, `COV-NO-ALERT`, `COV-NO-TEST`) are the action items; the matrix does not name any.
- **No suggested fixes.** A matrix without fixes is a status report, not an audit.
