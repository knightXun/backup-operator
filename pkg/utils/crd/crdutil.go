package crd

import (
	"errors"
	"github.com/backup-operator/pkg/apis/mydumper/v1alpha1"
	crdutils "github.com/yisaer/crd-validation/pkg"
	extensionsobj "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"strings"
)

var (
	restoreAdditionalPrinterColumns []extensionsobj.CustomResourceColumnDefinition
	restoreReadyColumn              = extensionsobj.CustomResourceColumnDefinition{
		Name:     "Ready",
		Type:     "string",
		JSONPath: `.status.conditions[?(@.type=="Ready")].status`,
	}
	restoreStatusMessageColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:     "Status",
		Type:     "string",
		JSONPath: `.status.conditions[?(@.type=="Ready")].message`,
		Priority: 1,
	}
	backupAdditionalPrinterColumns []extensionsobj.CustomResourceColumnDefinition
	backupPathColumn               = extensionsobj.CustomResourceColumnDefinition{
		Name:        "BackupPath",
		Type:        "string",
		Description: "The full path of mydumper data",
		JSONPath:    ".status.backupPath",
	}
)

func init() {
	restoreAdditionalPrinterColumns = append(restoreAdditionalPrinterColumns,
		restoreReadyColumn)
	backupAdditionalPrinterColumns = append(backupAdditionalPrinterColumns)
}


func GetCrdKindFromKindName(kindName string) (v1alpha1.CrdKind, error) {
	switch strings.ToLower(kindName) {
	case v1alpha1.RestoreKindKey:
		return v1alpha1.DefaultCrdKinds.Restore, nil
	case v1alpha1.BackupKindKey:
		return v1alpha1.DefaultCrdKinds.Backup, nil
	default:
		return v1alpha1.CrdKind{}, errors.New("unknown CrdKind Name")
	}
}

func NewCustomResourceDefinition(crdKind v1alpha1.CrdKind, group string, labels map[string]string, validation bool) *extensionsobj.CustomResourceDefinition {
	crd := crdutils.NewCustomResourceDefinition(crdutils.Config{
		SpecDefinitionName:    crdKind.SpecName,
		EnableValidation:      validation,
		Labels:                crdutils.Labels{LabelsMap: labels},
		ResourceScope:         string(extensionsobj.NamespaceScoped),
		Group:                 group,
		Kind:                  crdKind.Kind,
		Version:               v1alpha1.Version,
		Plural:                crdKind.Plural,
		ShortNames:            crdKind.ShortNames,
		GetOpenAPIDefinitions: v1alpha1.GetOpenAPIDefinitions,
	})
	addAdditionalPrinterColumnsForCRD(crd, crdKind)
	return crd
}

func addAdditionalPrinterColumnsForCRD(crd *extensionsobj.CustomResourceDefinition, crdKind v1alpha1.CrdKind) {
	switch crdKind.Kind {
	case v1alpha1.DefaultCrdKinds.Restore.Kind:
		crd.Spec.AdditionalPrinterColumns = restoreAdditionalPrinterColumns
		break
	case v1alpha1.DefaultCrdKinds.Backup.Kind:
		crd.Spec.AdditionalPrinterColumns = backupAdditionalPrinterColumns
		break
	default:
		break
	}
}