# Postmortem: departures-api returned 503 for 30% of requests

**Status:** Final
**Severity:** Sev2 (lab exercise: planned fault injection)
**Date:** 2026-11-01
**Author:** Praveen Bhat
**Incident commander:** Praveen Bhat

> This was a planned game day to test the SLO alerts and the runbook,
> written up as if it were a real incident.

## Summary
A config change set `ERROR_RATE=0.3` on departures-api, so 30% of departure
requests returned 503. DeparturesErrorBudgetFastBurn paged after [ ] minutes.
Following the runbook, a `helm rollback` restored service [ ] minutes after the page.

## Impact
| Measure | Value |
|---|---|
| Customer-visible duration | [ ] minutes |
| Peak error ratio | about 30% (SLO allows 0.1%) |
| Error budget spent | [ ]% of 30-day budget |
| Detected by | Alert (FastBurn) |

## Timeline (UTC)
| Time | Event |
|---|---|
| [ ] | `helm upgrade --set errorRate=0.3` applied (release revision [ ]) |
| [ ] | FastBurn pending |
| [ ] | FastBurn firing |
| [ ] | Acknowledged, runbook opened |
| [ ] | Runbook step 3: `helm history` shows a config change minutes earlier |
| [ ] | `helm rollback` applied |
| [ ] | 5m error ratio back under 0.001 |
| [ ] | Alert resolved |

**Time to detect:** [ ] min **Time to mitigate:** [ ] min **Time to recover:** [ ] min

## What happened
A Helm upgrade changed one value, `errorRate`, with no code change and no review.
The canary and the CI gates only apply to new images, so a config-only change
went straight to all pods. Errors started at once. The 1h burn-rate window
crossed 14.4x after about 3 minutes, and the alert fired 2 minutes later.

## Contributing factors
- Config changes skip the review and the canary that image changes get.
- `--reuse-values` hides what changed: the diff is not visible in the command.
- No alert on config drift: nothing compared running values with git.

## What went well
- The multi-window alert fired within the expected time for a 30% error rate.
- The runbook led straight to the cause: step 3, "did something change?"
- Rollback fixed it without needing to understand the cause first.

## What went badly
- [Add anything that was slow or confusing when you followed the runbook.]

## Where we got lucky
- Only one service was affected, and its rollback was clean.

## Action items
| Action | Type | Owner | Due | Ticket |
|---|---|---|---|---|
| Move Helm values into git, deployed by Argo CD (no manual `helm upgrade`) | Prevent | Praveen | Mini-project | |
| Canary config changes too, with the SLO analysis gate | Prevent | Praveen | Mini-project | |
| Drift report: running values vs git | Detect | Praveen | Mini-project | |
| Add `helm diff` to the runbook before any manual change | Mitigate | Praveen | Done | |

## Lessons
Config is code. A one-line config change can cause the same outage as a bad
release, so it needs the same review, canary and rollback path.