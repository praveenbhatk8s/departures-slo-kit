# Runbook: departures-api error budget burn

**Alerts:** DeparturesErrorBudgetFastBurn, DeparturesErrorBudgetSlowBurn, DeparturesTargetDown
**Impact:** riders and station displays may get no departure times.
**Goal:** stop the bleeding first. Find the root cause later.

## 1. Acknowledge (first 2 minutes)
- Acknowledge the page so the secondary knows you have it.
- Open an incident channel: `#inc-YYYYMMDD-departures`.
- FastBurn = assume Sev2 until proven otherwise.

## 2. How bad is it? (5 minutes)
```promql
# Current error ratio (SLO allows 0.001)
job:slo_errors_per_request:ratio_rate5m{job="departures-api"}

# Errors by pod: one bad pod, or all of them?
sum by (pod, code) (rate(http_requests_total{job="departures-api"}[5m]))
```
- **One pod:** delete it. `kubectl -n departures delete pod <pod>`
- **All pods:** likely a deploy, a config change or a dependency. Go to step 3.

## 3. Did something change? (most incidents are changes)
```bash
helm history departures-api -n departures        # recent releases
kubectl -n departures get events --sort-by=.lastTimestamp | tail -20
helm get values departures-api -n departures     # current config (check ERROR_RATE!)
```

## 4. Mitigate: roll back first, debug second
If a release happened in the last few hours, roll back without waiting for proof:
```bash
helm rollback departures-api -n departures       # back to the previous revision
```
Watch `job:slo_errors_per_request:ratio_rate5m` fall back below 0.001.
FastBurn should resolve within about 5 to 10 minutes.

## 5. If there was no change
```bash
kubectl -n departures logs -l app=departures-api --tail=100
kubectl -n departures top pods                   # needs metrics-server
kubectl -n departures describe pods -l app=departures-api
```
- Out of memory or CPU throttling → scale up: `kubectl -n departures scale deploy/departures-api --replicas=4`
- Dependency failing → escalate to its owner and post an update.

## 6. TargetDown: we are blind
Prometheus cannot scrape the app, so the SLO alerts cannot fire.
```bash
kubectl -n departures get pods,svc,servicemonitor
```
Check that the pods run, the Service has endpoints, and the ServiceMonitor port name is `http`.

## 7. Communicate
- Post an update in the incident channel every 30 minutes, even if nothing changed.
- If customers are affected for more than 15 minutes, tell customer support
  what riders see and when you will update next.

## 8. Close
- Resolve the page only after the 5m error ratio stays under 0.001 for 15 minutes.
- Open a postmortem from `postmortems/TEMPLATE.md` if the alert was FastBurn,
  or if the incident spent more than 20% of the budget.

**Escalation:** primary on-call → secondary on-call → engineering manager.