package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CitusClusterSpec defines the desired state of CitusCluster
// +k8s:openapi-gen=true
type CitusClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	// Size int32 `json:"size"`
	Keeper KeeperSpec `json:"keeper"`
	Proxy  ProxySpec  `json:"proxy"`
}

// CitusClusterStatus defines the observed state of CitusCluster
// +k8s:openapi-gen=true
type CitusClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	// Nodes are the names of the citus pods
	Nodes []string `json:"nodes"`
}

// ResourceRequirement is resource requirements for a pod
type ResourceRequirement struct {
	// CPU is how many cores a pod requires
	CPU string `json:"cpu,omitempty"`
	// Memory is how much memory a pod requires
	Memory string `json:"memory,omitempty"`
	// Storage is storage size a pod requires
	Storage string `json:"storage,omitempty"`
}

// ContainerSpec is the container spec of a pod
type ContainerSpec struct {
	Image           string               `json:"image"`
	ImagePullPolicy corev1.PullPolicy    `json:"imagePullPolicy,omitempty"`
	Requests        *ResourceRequirement `json:"requests,omitempty"`
	Limits          *ResourceRequirement `json:"limits,omitempty"`
}

// KeeperSpec the keeper specification
type KeeperSpec struct {
	ContainerSpec
	Size                 int32               `json:"size"`
	NodeSelector         map[string]string   `json:"nodeSelector,omitempty"`
	NodeSelectorRequired bool                `json:"nodeSelectorRequired,omitempty"`
	StorageClassName     string              `json:"storageClassName,omitempty"`
	Tolerations          []corev1.Toleration `json:"tolerations,omitempty"`
}

// ProxySpec the proxy specification
type ProxySpec struct {
	MasterPort  int32  `json:"masterPort"`
	StandbyPort int32  `json:"standbyPort"`
	Type        string `json:"type,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CitusCluster is the Schema for the citusclusters API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type CitusCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CitusClusterSpec   `json:"spec,omitempty"`
	Status CitusClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CitusClusterList contains a list of CitusCluster
type CitusClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CitusCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CitusCluster{}, &CitusClusterList{})
}
