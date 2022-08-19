// This is a generated file. Do not edit directly.
// Run hack/pin-dependency.sh to change pinned dependency versions.
// Run hack/update-vendor.sh to update go.mod files and the vendor directory.

module kubesphere.io/api

go 1.13

require (
	github.com/go-openapi/spec v0.19.7
	k8s.io/api v0.23.2
	k8s.io/apimachinery v0.23.2
	k8s.io/kube-openapi v0.0.0-20210421082810-95288971da7e
)

replace (
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20210305001622-591a79e4bda7
)
