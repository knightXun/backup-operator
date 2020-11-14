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
