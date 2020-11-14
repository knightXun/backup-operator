module github.com/backup-operator

go 1.15

require (
	github.com/aws/aws-sdk-go v1.35.25
	github.com/dlintw/goconf v0.0.0-20120228082610-dcc070983490
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1
	github.com/xelabs/go-mysqlstack v0.0.0-20200603045106-7ffcfc8ed3c2
	github.com/yisaer/crd-validation v0.0.3
	go.uber.org/zap v1.16.0
	google.golang.org/grpc v1.27.0
	k8s.io/api v0.19.4
	k8s.io/apiextensions-apiserver v0.19.4
	k8s.io/apimachinery v0.19.4
	k8s.io/code-generator v0.19.4
	k8s.io/klog v1.0.0
	k8s.io/kubectl v0.19.4
)

replace (
	cloud.google.com/go => github.com/googleapis/google-cloud-go v0.53.0
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20200221231518-2aa609cf4a9d
	golang.org/x/exp => github.com/golang/exp v0.0.0-20200221183520-7c80518d1cc7
	golang.org/x/image => github.com/golang/image v0.0.0-20200119044424-58c23975cae1
	golang.org/x/lint => github.com/golang/lint v0.0.0-20200130185559-910be7a94367
	golang.org/x/mobile => github.com/golang/mobile v0.0.0-20200222142934-3c8601c510d0
	golang.org/x/mod => github.com/golang/mod v0.2.0
	golang.org/x/net => github.com/golang/net v0.0.0-20200222125558-5a598a2470a0
	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys => github.com/golang/sys v0.0.0-20200219091948-cb0a6d8edb6c
	golang.org/x/text => github.com/golang/text v0.3.2
	golang.org/x/time => github.com/golang/time v0.0.0-20191024005414-555d28b269f0
	golang.org/x/tools => github.com/golang/tools v0.0.0-20200221224223-e1da425f72fd
	golang.org/x/xerrors => github.com/golang/xerrors v0.0.0-20191204190536-9bdfabe68543
	google.golang.org/api => github.com/googleapis/google-api-go-client v0.17.0
	google.golang.org/appengine => github.com/golang/appengine v1.6.5
	google.golang.org/genproto => github.com/googleapis/go-genproto v0.0.0-20200218151345-dad8c97a84f5
	google.golang.org/grpc => github.com/grpc/grpc-go v1.26.0
	k8s.io/api => k8s.io/api v0.18.2
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.2
	k8s.io/apiserver => k8s.io/apiserver v0.18.2
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.18.6
	k8s.io/client-go => k8s.io/client-go v0.18.2
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.18.6
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.18.6
	k8s.io/code-generator => k8s.io/code-generator v0.18.2
	k8s.io/component-base => k8s.io/component-base v0.18.2
	k8s.io/cri-api => k8s.io/cri-api v0.18.6
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.18.6
	k8s.io/klog => k8s.io/klog v1.0.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.18.6
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.18.6
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.18.6
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.18.6
	k8s.io/kubectl => k8s.io/kubectl v0.18.6
	k8s.io/kubelet => k8s.io/kubelet v0.18.6
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.18.6
	k8s.io/metrics => k8s.io/metrics v0.18.2
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.18.6
)
