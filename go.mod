module github.com/backup-operator

go 1.15

require (
	github.com/aws/aws-sdk-go v1.35.25
	github.com/dlintw/goconf v0.0.0-20120228082610-dcc070983490
	github.com/pingcap/br v4.0.0-beta.2.0.20201110065050-a3517d674652+incompatible
	github.com/pingcap/errors v0.11.4
	github.com/pingcap/kvproto v0.0.0-20201104042953-62eb316d5182
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1
	github.com/xelabs/go-mysqlstack v0.0.0-20200603045106-7ffcfc8ed3c2
	go.uber.org/zap v1.16.0
	k8s.io/api v0.19.3
	k8s.io/apimachinery v0.19.3
	k8s.io/klog v1.0.0
)
