.PHONY: default

default:

# Go Build

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

GO_FILES = $(shell find . -name *.go)

dist/echo-server-$(GOOS)-$(GOARCH): $(GO_FILES) go.mod go.sum
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $@ ./cmd/echo-server

dist/echo-client-$(GOOS)-$(GOARCH): $(GO_FILES) go.mod go.sum
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $@ ./cmd/echo-client

.PHONY: build

build: dist/echo-server-$(GOOS)-$(GOARCH) dist/echo-client-$(GOOS)-$(GOARCH)

DOCKER_IMAGE ?= ghcr.io/htdvisser/echo-server
DOCKER_IMAGE_VERSION ?= latest

# Docker Build

.PHONY: docker

docker:
	GOOS=linux GOARCH=amd64 $(MAKE) build
	docker build --build-arg TARGETOS=linux --build-arg TARGETARCH=amd64 -t $(DOCKER_IMAGE):$(DOCKER_IMAGE_VERSION) -f build/Dockerfile .
	docker push $(DOCKER_IMAGE):$(DOCKER_IMAGE_VERSION)

.PHONY: docker-multiarch

docker-multiarch:
	GOOS=linux GOARCH=arm $(MAKE) build
	GOOS=linux GOARCH=arm64 $(MAKE) build
	GOOS=linux GOARCH=386 $(MAKE) build
	GOOS=linux GOARCH=amd64 $(MAKE) build
	DOCKER_CLI_EXPERIMENTAL=enabled docker buildx build --push --platform linux/arm/v7,linux/arm64/v8,linux/386,linux/amd64 -t $(DOCKER_IMAGE):$(DOCKER_IMAGE_VERSION) -f build/Dockerfile .

# Istio

ISTIO_VERSION=1.4.2

istio-$(ISTIO_VERSION):
	curl -sSL https://github.com/istio/istio/releases/download/$(ISTIO_VERSION)/istio-$(ISTIO_VERSION)-osx.tar.gz | tar xz

deploy/istio-manifest.yaml: istio-$(ISTIO_VERSION)
	istio-$(ISTIO_VERSION)/bin/istioctl manifest generate -f deploy/istio-overlay.yaml > deploy/istio-manifest.yaml

.PHONY: istio-install

istio-install: istio-$(ISTIO_VERSION) deploy/istio-manifest.yaml
	istio-$(ISTIO_VERSION)/bin/istioctl manifest apply -f deploy/istio-overlay.yaml -f deploy/istio-manifest.yaml

.PHONY: istio-verify

istio-verify: istio-$(ISTIO_VERSION) deploy/istio-manifest.yaml
	istio-$(ISTIO_VERSION)/bin/istioctl verify-install -f deploy/istio-manifest.yaml

.PHONY: istio-uninstall

istio-uninstall: istio-$(ISTIO_VERSION) deploy/istio-manifest.yaml
	kubectl delete -f deploy/istio-manifest.yaml

.PHONY: istio-clean

istio-clean:
	rm -f istio-$(ISTIO_VERSION) deploy/istio-manifest.yaml

# Deploy

.PHONY: deploy

deploy:
	kubectl apply -f deploy/1-deployment-v1alpha1.yaml
	kubectl apply -f deploy/2-service.yaml
	kubectl apply -f deploy/3-gateway.yaml
	kubectl apply -f deploy/4-virtual-service.yaml

.PHONY: undeploy

undeploy:
	kubectl delete -f deploy/4-virtual-service.yaml
	kubectl delete -f deploy/3-gateway.yaml
	kubectl delete -f deploy/2-service.yaml
	kubectl delete -f deploy/1-deployment-v1alpha1.yaml
