package scheme

import (
	"github.com/inspur/pkg/api/aggregation/install"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var scheme = runtime.NewScheme()

var Codecs = serializer.NewCodecFactory(scheme)

var ParameterCodec = runtime.NewParameterCodec(scheme)

func init() {
	install.Install(scheme)
}
