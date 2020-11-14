package generate

import (
	"errors"
	crd2 "github.com/backup-operator/pkg/utils/crd"

	"github.com/spf13/cobra"
	"github.com/backup-operator/pkg/apis/mydumper/v1alpha1"
	crdutils "github.com/yisaer/crd-validation/pkg"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

const (
	usage = "usage: to-crdgen generate backup|restore"
)

func AddGenerateCommand(config *crdutils.Config) *cobra.Command {
	generatedCommand := &cobra.Command{
		Use:   "generate",
		Short: "Generate CRD",
		Long:  "Generate CRD by crd-util according to types",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(generate(config, args))
		},
	}
	return generatedCommand
}

func initConfig(kind v1alpha1.CrdKind, config *crdutils.Config) {
	config.Kind = kind.Kind
	config.Plural = kind.Plural
	config.ShortNames = kind.ShortNames
	config.SpecDefinitionName = kind.SpecName
}

func generate(config *crdutils.Config, args []string) error {
	if len(args) < 1 || len(args) > 1 {
		return errors.New(usage)
	}
	crdKind, err := crd2.GetCrdKindFromKindName(args[0])
	if err != nil {
		return errors.New(usage)
	}
	initConfig(crdKind, config)
	crd := crd2.NewCustomResourceDefinition(
		crdKind,
		config.Group, config.Labels.LabelsMap, config.EnableValidation)
	return crdutils.MarshallCrd(crd, config.OutputFormat)
}
