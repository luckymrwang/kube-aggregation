package clusterpedia

import (
	"net/url"

	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/inspur/pkg/api/clusterpedia/fields"
)

type OrderBy struct {
	Field string
	Desc  bool
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ListOptions struct {
	metainternal.ListOptions

	Names        []string
	ClusterNames []string
	Namespaces   []string
	OrderBy      []OrderBy

	OwnerName          string
	OwnerUID           string
	OwnerGroupResource schema.GroupResource
	OwnerSeniority     int

	Since  *metav1.Time
	Before *metav1.Time

	WithContinue       *bool
	WithRemainingCount *bool

	// +k8s:conversion-fn:drop
	EnhancedFieldSelector fields.Selector

	// +k8s:conversion-fn:drop
	ExtraLabelSelector labels.Selector

	// +k8s:conversion-fn:drop
	URLQuery url.Values

	// RelatedResources []schema.GroupVersionKind

	OnlyMetadata bool
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CollectionResource struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	ResourceTypes []CollectionResourceType
	Items         []runtime.Object

	Continue           string
	RemainingItemCount *int64
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CollectionResourceList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []CollectionResource
}

type CollectionResourceType struct {
	Group    string
	Version  string
	Kind     string
	Resource string
}

func (t CollectionResourceType) GroupResource() schema.GroupResource {
	return schema.GroupResource{
		Group:    t.Group,
		Resource: t.Resource,
	}
}
