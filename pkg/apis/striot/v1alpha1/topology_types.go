package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TopologySpec defines the desired state of Topology
type TopologySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Partitions []Partition `json:"partitions"`
	Order      []int32     `json:"order"`
}

// TopologyStatus defines the observed state of Topology
type TopologyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Topology is the Schema for the topologies API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=topologies,scope=Namespaced
type Topology struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TopologySpec   `json:"spec,omitempty"`
	Status TopologyStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TopologyList contains a list of Topology
type TopologyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Topology `json:"items"`
}

// Partition contains an ID and definition for a partition
type Partition struct {
	ID          int32       `json:"id"`
	Image       string      `json:"image"`
	ConnectType ConnectType `json:"connectType"`
}

// ConnectType specifies ingress and egress types
type ConnectType struct {
	Ingress ConnectProtocol `json:"ingress"`
	Egress  ConnectProtocol `json:"egress"`
}

// ConnectProtocol defines the type of ingress / egress
type ConnectProtocol string

const (
	// ProtocolTCP is the TCP protocol
	ProtocolTCP ConnectProtocol = "TCP"
	// ProtocolKafka is the Kafka protocol
	ProtocolKafka ConnectProtocol = "KAFKA"
	// ProtocolMQTT is the MQTT protocol
	ProtocolMQTT ConnectProtocol = "MQTT"
)

func (cp ConnectProtocol) String() string {
	return string(cp)
}

func init() {
	SchemeBuilder.Register(&Topology{}, &TopologyList{})
}
