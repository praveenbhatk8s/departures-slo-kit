    # SLO: departures-api availability

| Field | Value |
|---|---|
| Service | departures-api |
| Owner | Platform / SRE (Praveen Bhat) |
| User journey | A rider or a station display asks "when is my next bus?" |
| SLI | Share of `/departures` requests that do not return a 5xx |
| Target | 99.9% |
| Window | 30 days, rolling |
| Error budget | 0.1% of requests (about 43 minutes of full outage per 30 days) |
| Source | `http_requests_total{job="departures-api"}`, recorded in `slo/rules.yaml` |

## Why this SLI
A blank departure board is what riders notice. A 5xx means we gave them nothing.
CPU, memory and pod restarts are not in the SLI: they are causes, not symptoms.

## Why 99.9%, not 99.99%
- 99.9% allows about 43 minutes of budget a month: enough to deploy weekly and
  run experiments safely.
- 99.99% allows about 4 minutes: that needs multi-region and much more cost.
  We will raise it only if customers tell us 99.9% is not good enough.

## Alerting (see `slo/rules.yaml`)
| Alert | Burn rate | Windows | Budget spent | Action |
|---|---|---|---|---|
| FastBurn | 14.4x | 1h and 5m | 2% | Page now |
| SlowBurn | 6x | 6h and 30m | 5% | Page |
| Ticket | 3x | 1d and 2h | 10% | Ticket, working hours |
| TargetDown | n/a | 5m | unknown | Page: we are blind |

## Error budget policy
- **Budget left:** ship normally.
- **Budget spent:** feature releases pause. Only reliability fixes ship,
  behind a canary, until the 30-day SLI is back above target.
- **Disagreements** go to the engineering manager and product lead, with data.
- **One incident spends more than 20% of the budget:** a postmortem is required.

## What is not covered (yet)
- Latency (planned SLO: 99% of requests under 300 ms)
- Data freshness (are the departure times correct, not just returned?)
- Dependencies: the real service would call vehicle-location feeds.

## Review
Monthly, with product. An SLO nobody reviews is just a number.