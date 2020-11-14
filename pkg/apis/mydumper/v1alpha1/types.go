package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +k8s:openapi-gen=true
// Backup is the control script's spec
type Backup struct {
	metav1.TypeMeta `json:",inline"`
	// +k8s:openapi-gen=false
	metav1.ObjectMeta `json:"metadata"`

	// Spec defines the behavior of a Backup
	Spec BackupSpec `json:"spec"`

	// +k8s:openapi-gen=false
	// Most recently observed status of the Backup
	Status BackupStatus `json:"status"`
}


// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +k8s:openapi-gen=true
// BackupList is Backup list
type BackupList struct {
	metav1.TypeMeta `json:",inline"`
	// +k8s:openapi-gen=false
	metav1.ListMeta `json:"metadata"`

	Items []Backup `json:"items"`
}


// +k8s:openapi-gen=true
// BackupSpec describes the attributes that a user creates on a mydumper
type BackupSpec struct {

}


// BackupStatus represents the current status of a mydumper.
type BackupStatus struct {
	Conditions []BackupCondition `json:"conditions,omitempty"`
}


type BackupCondition struct {
	// Type of the condition.
	Type BackupConditionType `json:"type"`

	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// The last time this condition was updated.

	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// +optional
	Reason string `json:"reason,omitempty"`
	// +optional
	Message string `json:"message,omitempty"`
}

type BackupConditionType string

const (
	BackupDone BackupConditionType = "Done"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +k8s:openapi-gen=true
// Restore is the control script's spec
type Restore struct {
	metav1.TypeMeta `json:",inline"`
	// +k8s:openapi-gen=false
	metav1.ObjectMeta `json:"metadata"`

	// Spec defines the behavior of a Restore
	Spec RestoreSpec `json:"spec"`

	// +k8s:openapi-gen=false
	// Most recently observed status of the Restore
	Status RestoreStatus `json:"status"`
}


// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +k8s:openapi-gen=true
// RestoreList is Restore list
type RestoreList struct {
	metav1.TypeMeta `json:",inline"`
	// +k8s:openapi-gen=false
	metav1.ListMeta `json:"metadata"`

	Items []Backup `json:"items"`
}


// +k8s:openapi-gen=true
// RestoreSpec describes the attributes that a user creates on a mydumper
type RestoreSpec struct {

}


// BackupStatus represents the current status of a mydumper.
type RestoreStatus struct {
	Conditions []RestoreCondition `json:"conditions,omitempty"`
}


type RestoreCondition struct {
	// Type of the condition.
	Type RestoreConditionType `json:"type"`

	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// The last time this condition was updated.

	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// +optional
	Reason string `json:"reason,omitempty"`
	// +optional
	Message string `json:"message,omitempty"`
}

type RestoreConditionType string

const (
	RestoreDone RestoreConditionType = "Done"
)