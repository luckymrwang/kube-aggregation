# Copyright 2018 The KubeSphere Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.


# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

GV="cluster:v1alpha1"
MANIFESTS="cluster/*"

# App Version
APP_VERSION = v1.0.0

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

OUTPUT_DIR=bin
ifeq (${GOFLAGS},)
	# go build with vendor by default.
	export GOFLAGS=-mod=vendor
endif
define ALL_HELP_INFO
# Build code.
#
# Args:
#   WHAT: Directory names to build.  If any of these directories has a 'main'
#     package, the build will produce executable files under $(OUT_DIR).
#     If not specified, "everything" will be built.
#   GOFLAGS: Extra flags to pass to 'go' when building.
#   GOLDFLAGS: Extra linking flags passed to 'go' when building.
#   GOGCFLAGS: Additional go compile flags passed to 'go' when building.
#
# Example:
#   make
#   make all
#   make all WHAT=cmd/ks-apiserver
#     Note: Use the -N -l options to disable compiler optimizations an inlining.
#           Using these build options allows you to subsequently use source
#           debugging tools like delve.
endef
.PHONY: all
all: test ks-apiserver;$(info $(M)...Begin to test and build all of binary.) @ ## Test and build all of binary.

help:
	@grep -hE '^[ a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-17s\033[0m %s\n", $$1, $$2}'

.PHONY: binary
# Build all of binary
binary: | ks-apiserver; $(info $(M)...Build all of binary.) @ ## Build all of binary.

# Build ks-apiserver binary
ks-apiserver: ; $(info $(M)...Begin to build ks-apiserver binary.)  @ ## Build ks-apiserver.
	 hack/gobuild.sh cmd/ks-apiserver;

# Run go vet against code
vet: ;$(info $(M)...Begin to run go vet against code.)  @ ## Run go vet against code.
	go vet ./pkg/... ./cmd/...

# Generate manifests e.g. CRD, RBAC etc.
manifests: ;$(info $(M)...Begin to generate manifests e.g. CRD, RBAC etc..)  @ ## Generate manifests e.g. CRD, RBAC etc.
	hack/generate_manifests.sh ${CRD_OPTIONS} ${MANIFESTS}

deploy: manifests ;$(info $(M)...Begin to deploy.)  @ ## Deploy.
	kubectl apply -f config/crds
	kustomize build config/default | kubectl apply -f -

deepcopy: ;$(info $(M)...Begin to deepcopy.)  @ ## Deepcopy.
	hack/generate_group.sh "deepcopy" kubesphere.io/api kubesphere.io/api ${GV} --output-base=staging/src/  -h "hack/boilerplate.go.txt"

container: ;$(info $(M)...Begin to build the docker image.)  @ ## Build the docker image.
	DRY_RUN=true hack/docker_build.sh

container-push: ;$(info $(M)...Begin to build and push.)  @ ## Build and Push.
	hack/docker_build.sh

container-cross: ; $(info $(M)...Begin to build container images for multiple platforms.)  @ ## Build container images for multiple platforms. Currently, only linux/amd64,linux/arm64 are supported.
	DRY_RUN=true hack/docker_build_multiarch.sh

container-cross-push: ; $(info $(M)...Begin to build and push.)  @ ## Build and Push.
	hack/docker_build_multiarch.sh

.PHONY: clean
clean: ;$(info $(M)...Begin to clean.)  @ ## Clean.
	-make -C ./pkg/version clean
	@echo "ok"

clientset:  ;$(info $(M)...Begin to find or download controller-gen.)  @ ## Find or download controller-gen,download controller-gen if necessary.
	./hack/generate_client.sh ${GV}

# Fix invalid file's license.
update-licenses: ;$(info $(M)...Begin to update licenses.)
	@hack/update-licenses.sh
