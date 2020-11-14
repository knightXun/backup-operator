package v1alpha1

import (
	extensionsobj "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

const (
	Version   = "v1alpha1"

	RestoreName    = "restores"
	RestoreKind    = "Restore"
	RestoreKindKey = "restore"

	BackupName    = "backups"
	BackupKind    = "Backup"
	BackupKindKey = "backup"
	SpecPath = "github.com/controller-operator/pkg/apis/mydumper/v1alpha1."
)

type CrdKind struct {
	Kind                    string
	Plural                  string
	SpecName                string
	ShortNames              []string
	AdditionalPrinterColums []extensionsobj.CustomResourceColumnDefinition
}

type CrdKinds struct {
	KindsString   	string
	Restore 		CrdKind
	Backup        	CrdKind
}

var DefaultCrdKinds = CrdKinds{
	KindsString:   	"",
	Restore: 		CrdKind{Plural: RestoreName, Kind: RestoreKind, ShortNames: []string{"rs"}, SpecName: SpecPath + RestoreKind},
	Backup:        	CrdKind{Plural: BackupName, Kind: BackupKind, ShortNames: []string{"bk"}, SpecName: SpecPath + BackupKind},
}