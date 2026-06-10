# Control Walk Framework

> See also: ../SKILL.md

The framework is a per-control evidence pack. For each control, the framework says what to record.

## The 6 cells per control

1. **Regulator frame + control ID** — the frame (SOX, ISO 27001, SOC 2, GDPR, HIPAA, PCI) and the canonical control identifier. Examples: SOX `ITGC.04`, ISO `A.9.4.1`, GDPR Article 17, HIPAA `164.312(a)(1)`, PCI `3.4`.
2. **Canonical enforcement symbol** — the symbol the control is implemented in. Examples: an interceptor, a guard, a middleware, a deletion worker.
3. **Per-site evidence** — for every enforcement site, a path + line + code excerpt. The excerpt is the smallest one that proves the control is invoked. Suitable for the auditor to re-verify by hand.
4. **Per-bypass finding** — every site that touches the protected data but does not go through the canonical enforcement. Each bypass is a finding with the path + line + a one-sentence explanation.
5. **Data-lineage cross-check** — which sinks the control covers; which sinks are uncovered. The cross-check is the link to the upstream `data-lineage-trace` output.
6. **Evidence-collector** — the path the auditor can run to re-verify the evidence. Examples: a log line id, a metrics query, a database query, a unit-test command.

## The 4 evidence-pack sections

The final pack is structured as four sections, in this order:

1. **Control rows** — every control in the auditor's scope, with the 6 cells above. A control with no enforcement site is highlighted as `COMP-MISSING`. A control with a reachable bypass is highlighted as `COMP-BYPASS`.
2. **Per-bypass findings** — every bypass, with the path + line + a one-sentence explanation. Each finding has a suggested fix.
3. **Data-lineage cross-check** — for each sink in the data-lineage report, name the controls that cover it and the controls that do not.
4. **Evidence-collector list** — for each control, the path the auditor can run to re-verify.

## The 5 regulator-frame cheat-sheet

| Frame | Canonical control ID | Enforcement symbol examples |
|---|---|---|
| SOX | `ITGC.04` (change management), `ITGC.05` (logical access) | PR approver, role-based access |
| ISO 27001 | `A.9.4.1` (access restriction), `A.10.1.1` (crypto) | ACL middleware, encryption interceptor |
| SOC 2 | `CC6.1` (logical access), `CC7.2` (monitoring) | RBAC, audit-log middleware |
| GDPR | Article 17 (right to erasure), Article 32 (security) | Erasure worker, encryption at rest |
| HIPAA | `164.312(a)(1)` (access control), `164.312(b)` (audit) | RBAC + audit-log middleware |
| PCI | `3.4` (PAN protection at rest), `4.1` (transmission) | Tokenisation at ingress, TLS |

## The "all clean" sentinel

The skill is not a linter. An evidence pack where every control row has enforcement sites and no bypass is the "no finding" case; the pack itself is the output. There is no canonical "all clean" line.
