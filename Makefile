CLUSTER := slo
IMAGE   := departures-api:0.1.0

.PHONY: up down cluster monitoring image deploy break heal

## Create everything from scratch
up: cluster monitoring image deploy

cluster:
	k3d cluster list $(CLUSTER) >/dev/null 2>&1 || k3d cluster create $(CLUSTER) --wait

monitoring:
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts --force-update
	helm upgrade --install monitoring prometheus-community/kube-prometheus-stack \
		--namespace monitoring --create-namespace \
		-f deploy/monitoring-values.yaml --wait --timeout 10m

## Build the app image and load it into the cluster (no registry needed)
image:
	docker build -f app/Dockerfile -t $(IMAGE) .
	k3d image import $(IMAGE) -c $(CLUSTER)

## Delete the whole lab
down:
	k3d cluster delete $(CLUSTER)

## Install or upgrade the app
deploy:
	helm upgrade --install departures-api deploy/chart \
		--namespace departures --create-namespace --wait

## Start an incident: 5% of requests fail
break:
	helm upgrade departures-api deploy/chart -n departures --reuse-values --set errorRate=0.05 --wait

## End the incident
heal:
	helm upgrade departures-api deploy/chart -n departures --reuse-values --set errorRate=0 --wait