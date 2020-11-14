package scheme

import (
	nbscheme "github.com/backup-operator/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	kubescheme "k8s.io/client-go/kubernetes/scheme"
)

// Scheme gathers the schemes of native resources and custom resources used by nebula-operator
// in favor of the generic controller-runtime/client
var Scheme = runtime.NewScheme()

func init() {
	v1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})
	utilruntime.Must(nbscheme.AddToScheme(Scheme))
	utilruntime.Must(kubescheme.AddToScheme(Scheme))
}