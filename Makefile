OS ?= $(shell go env GOOS)
ARCH ?= $(shell go env GOARCH)
ALL_PLATFORM = linux/amd64,linux/arm/v7,linux/arm64

# Image URL to use all building/pushing image targets
REGISTRY ?= arbreagile
IMAGE ?= cert-manager-webhook-bunny
FULL_IMAGE ?= $(REGISTRY)/$(IMAGE)

IMAGE_TAG ?= $(shell git rev-parse HEAD)

DOCKER_CLI_EXPERIMENTAL ?= enabled

KUBEBUILDER_VERSION=2.3.1

TEST_ZONE_NAME ?= example.com.

# Run tests
test: tests/kubebuilder
	TEST_ZONE_NAME=$(TEST_ZONE_NAME) go test -v ./... -coverprofile cover.out

cover:
	go tool cover -func=cover.out -o=cover.txt

test-with-cover: test cover

tests/kubebuilder:
	curl -fsSL https://github.com/kubernetes-sigs/kubebuilder/releases/download/v$(KUBEBUILDER_VERSION)/kubebuilder_$(KUBEBUILDER_VERSION)_$(OS)_$(ARCH).tar.gz -o kubebuilder-tools.tar.gz
	mkdir tests/kubebuilder
	tar -xvf kubebuilder-tools.tar.gz
	mv kubebuilder_$(KUBEBUILDER_VERSION)_$(OS)_$(ARCH)/bin tests/kubebuilder/
	rm kubebuilder-tools.tar.gz
	rm -R kubebuilder_$(KUBEBUILDER_VERSION)_$(OS)_$(ARCH)

clean-kubebuilder:
	rm -Rf tests/kubebuilder

compile:
	go build -v -o cert-manager-webhook-bunny main.go

docker-build:
	@echo "Building cert-manager-webhook-bunny for $(ARCH)"
	docker build . --platform=$(OS)/$(ARCH) -f Dockerfile -t $(FULL_IMAGE):$(IMAGE_TAG)-$(ARCH)

docker-buildx-all:
	@echo "Making release for tag $(IMAGE_TAG)"
	docker buildx build --platform=$(ALL_PLATFORM) --push -t $(FULL_IMAGE):$(IMAGE_TAG) .

release: docker-buildx-all
