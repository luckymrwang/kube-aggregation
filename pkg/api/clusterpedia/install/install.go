package install

import (
	internal "github.com/inspur/pkg/api/clusterpedia"
	"github.com/inspur/pkg/api/clusterpedia/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

func Install(scheme *runtime.Scheme) {
	utilruntime.Must(internal.Install(scheme))
	utilruntime.Must(v1beta1.Install(scheme))
	utilruntime.Must(scheme.SetVersionPriority(v1beta1.SchemeGroupVersion))
}
