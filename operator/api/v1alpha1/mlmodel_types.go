package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HealthThresholds struct {
	MaxDriftScore resource.Quantity `json:"maxDriftScore"`
	MinAccuracy   resource.Quantity `json:"minAccuracy"`
	MaxLatencyMs  int64             `json:"maxLatencyMs"`
}

type MLModelSpec struct {
	Image            string           `json:"image"`
	Port             int32            `json:"port"`
	Replicas         int32            `json:"replicas"`
	HealthThresholds HealthThresholds `json:"healthThresholds"`
}

type MLModelStatus struct {
	Phase      string            `json:"phase,omitempty"`
	DriftScore resource.Quantity `json:"driftScore,omitempty"`
	Message    string            `json:"message,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

type MLModel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MLModelSpec   `json:"spec,omitempty"`
	Status            MLModelStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

type MLModelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MLModel `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MLModel{}, &MLModelList{})
}
