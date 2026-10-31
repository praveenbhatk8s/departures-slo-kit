CLUSTER := slo
IMAGE   := departures-api:0.1.0

.PHONY: up down cluster monitoring image

## Create everything from scratch
up: cluster monitoring image

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