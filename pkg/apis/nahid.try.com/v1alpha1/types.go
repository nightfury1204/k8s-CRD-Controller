package v1alpha1

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

//PodWatch is specification for PodWatch resource
type PodWatch struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodWatchSpec   `json:"spec"`
	Status PodWatchStatus `json:"status"`
}

//PodwatchSpec is spec for Podwatcher resource
type PodWatchSpec struct {
	metav1.LabelSelector `json:"selector"`
	Replicas             int32       `json:"replicas"`
	Template             PodTemplate `json:"template"`
}

//PodTemplate is template of pod for PodWatcher resource
type PodTemplate struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              apiv1.PodSpec `json:"spec"`
}

/*
//PodSpec is specification of Pod resource
type PodSpec struct {
	Name          string   `json:"name,omitempty"`
	Args          []string `json:"args,omitempty"`
	Image         string   `json:"image,omitempty"`
	ContainerPort int32    `json:"containerPort,,omitempty"`
}
*/

//PodWatchStatus is the status for PodWatcher resource
type PodWatchStatus struct {
	AvailabelReplicas int32 `json:"availabelReplicas"`
	CurrentlyProcessing int32 `json:"currently_processing"`
	PodNameList []string `json:"pod_name_list"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

//PodWatcherList is a list of PodWatcher resource
type PodWatchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []PodWatch `json:"items"`
}
