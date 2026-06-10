# Coverage Audit Framework

> See also: ../SKILL.md

The framework is a per-site coverage matrix. For each site, the framework says what to record.

## The 4 cells per site

1. **Site** — the path + line of the in-scope code point. Examples: a `throw` site, a `catch` site, an external HTTP call site, a DB write site.
2. **Test coverage** — the existing test that exercises this site, named by file + test class/function. `MISSING` if no test exists. A test that mocks the dependency does not count — the matrix says so explicitly.
3. **Alert coverage** — the existing alert that fires when the site fails in production. Named by alert rule id or by the metrics-backend query. `MISSING` if no alert.
4. **Runbook coverage** — the existing runbook entry that names this site. Named by runbook path + section. `MISSING` if no runbook.

## The 4 classification findings

A site is classified by what is missing:

- **`COV-SILENT`** — missing all three (test, alert, runbook). Highest-priority gap; the site fails in production with no signal.
- **`COV-NO-RUNBOOK`** — has test + alert, no runbook. The on-call is paged, but the runbook does not name this site; the responder's first 10 minutes are wasted.
- **`COV-NO-ALERT`** — has test + runbook, no alert. The test catches the regression; the runbook names the response; but the page does not fire.
- **`COV-NO-TEST`** — has alert + runbook, no test. The site is monitored and runbooked, but the test suite does not catch a regression.

## The 3 matrix sections

The final matrix is structured as three sections, in this order:

1. **Per-site coverage rows** — every site in the in-scope path, with the 4 cells above. Sites with any `MISSING` are highlighted.
2. **Gap findings** — one row per site with a classification (`COV-SILENT`, `COV-NO-RUNBOOK`, `COV-NO-ALERT`, `COV-NO-TEST`). Sorted by classification priority.
3. **Suggested fixes** — for each gap, the smallest coverage to add. Examples: "add a unit test that exercises this `throw` site"; "add an alert on the error-rate metric for this call".

## The "all clean" sentinel

The skill is not a linter. A matrix where every site has test + alert + runbook is the "no finding" case; the matrix itself is the output. There is no canonical "all clean" line.
