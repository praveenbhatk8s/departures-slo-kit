# departures-slo-kit

A small, complete reliability practice for one service: an SLO, multi-window
burn-rate alerts, a tested runbook, and a blameless postmortem.
The service is a toy bus-departure API, because what "reliable" means
depends on the domain: a blank departure board is what riders notice.

## What's inside
| Path | What it is |
|---|---|
| `app/` | Go service: `/departures`, `/healthz`, `/metrics`, with `ERROR_RATE` for fault injection |
| `deploy/chart/` | Helm chart: Deployment, Service, ServiceMonitor, load generator |
| `deploy/monitoring-values.yaml` | kube-prometheus-stack settings |
| `slo/SLO.md` | SLI, target, error budget, error budget policy |
| `slo/rules.yaml` | Recording rules and burn-rate alerts (PrometheusRule) |
| `runbooks/high-burn.md` | What the on-call engineer does when paged |
| `postmortems/` | Template and a postmortem of a real game day |

## How it works
```
app fails on purpose → /metrics shows 5xx → Prometheus scrapes it
→ recording rules compute the error ratio over 5m … 1d
→ burn-rate alert fires with a runbook link → rollback → postmortem
```

## Run it
Needs Docker, k3d, kubectl and Helm.
```bash
make up       # cluster, Prometheus stack, app, SLO rules (about 5 minutes)
make break    # 30% of requests fail
              # FastBurn fires in about 5 minutes
make heal     # or follow the runbook: helm rollback
make down     # delete everything
```
Prometheus: `kubectl -n monitoring port-forward svc/monitoring-kube-prometheus-prometheus 9090:9090`
Alertmanager: `kubectl -n monitoring port-forward svc/monitoring-kube-prometheus-alertmanager 9093:9093`
Grafana: `kubectl -n monitoring port-forward svc/monitoring-grafana 3000:80` (admin / admin)

## Design decisions
- **Page on symptoms, not causes.** The SLI is failed requests. CPU and memory are dashboards, not pages.
- **Multi-window, multi-burn-rate alerts.** The long window proves the problem is real. The short window lets the alert resolve soon after the fix. Bigger outages page faster.
- **Alert on blindness.** `DeparturesTargetDown` fires if Prometheus can't scrape the app, because silent SLO alerts look like health.
- **SLO rules live outside the Helm chart.** They're reviewed separately by product and SRE, and Helm would mangle `{{ $value }}`.
- **Secure by default.** Distroless, non-root (UID 65532), read-only filesystem, no capabilities.

## What I learned
- `runAsNonRoot` needs a numeric UID: Kubernetes cannot verify a user name like `nonroot`.
- A local toolchain (Go 1.26) and the build image (Go 1.23) drifted. CI should read the Go version from `go.mod`.
- My own game day showed that config-only changes skip every safety gate. See the postmortem.

## Next
A ship-safely pipeline: CI gates, Argo Rollouts canary gated on this SLO,
and a customer-managed site that pulls pinned versions.

Based on the multi-window, multi-burn-rate approach in Google's
*Site Reliability Workbook*, chapter "Alerting on SLOs".