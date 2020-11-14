package cmd

import (
	"github.com/spf13/cobra"
	"github.com/backup-operator/pkg/apis/mydumper/v1alpha1"
	"github.com/backup-operator/pkg/to-crdgen/cmd/generate"
	crdutils "github.com/yisaer/crd-validation/pkg"
	extensionsobj "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

const (
	tkcLongDescription = `
		"to-crdgen (backup-operator crd generator) is a tool to help generate CRD automatically.
`
)

var (
	cfg crdutils.Config
)

func initFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&cfg.EnableValidation, "with-validation", true, "Add CRD validation field, default: true")
	cmd.Flags().StringVar(&cfg.Group, "apigroup", v1alpha1.GroupName, "CRD api group")
	cmd.Flags().StringVar(&cfg.OutputFormat, "output", "yaml", "output format: json|yaml")
	cmd.Flags().StringVar(&cfg.ResourceScope, "scope", string(extensionsobj.NamespaceScoped), "CRD scope: 'Namespaced' | 'Cluster'.  Default: Namespaced")
	cmd.Flags().StringVar(&cfg.Version, "version", v1alpha1.Version, "CRD version, default: 'v1alpha1'")
}

func NewToCrdGenRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "to-crdgen",
		Short: "to-crdgen is a small tool to generate crd",
		Long:  tkcLongDescription,
		Run:   runHelp,
	}
	initFlags(rootCmd)
	rootCmd.AddCommand(generate.AddGenerateCommand(&cfg))
	return rootCmd
}

func runHelp(cmd *cobra.Command, _ []string) {
	cmd.Help()
}
