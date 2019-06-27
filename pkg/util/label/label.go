package label

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	// ManagedByLabelKey is Kubernetes recommended label key, it represents the tool being used to manage the operation of an application
	// For resources managed by Citus Operator, its value is always citus-operator
	ManagedByLabelKey string = "app.kubernetes.io/managed-by"
	// ComponentLabelKey is Kubernetes recommended label key, it represents the component within the architecture
	ComponentLabelKey string = "app.kubernetes.io/component"
	// NameLabelKey is Kubernetes recommended label key, it represents the name of the application
	// It should always be cell-cluster in our case.
	NameLabelKey string = "app.kubernetes.io/name"
	// InstanceLabelKey is Kubernetes recommended label key, it represents a unique name identifying the instance of an application
	// It's set by helm when installing a release
	InstanceLabelKey string = "app.kubernetes.io/instance"
	// NamespaceLabelKey is label key used in PV for easy querying
	NamespaceLabelKey string = "app.kubernetes.io/namespace"
	// ProxyLabelVal is proxy label value
	ProxyLabelVal string = "proxy"
	// KeeperLabelVal is Keeper label value
	KeeperLabelVal string = "keeper"
)

// Label is the label field in metadata
type Label map[string]string

// New initialize a new Label
func New() Label {
	return Label{
		NameLabelKey:      "citus-cluster",
		ManagedByLabelKey: "citus-operator",
	}
}

// Instance adds instance kv pair to label
func (l Label) Instance(name string) Label {
	l[InstanceLabelKey] = name
	return l
}

// Component adds component kv pair to label
func (l Label) Component(name string) Label {
	l[ComponentLabelKey] = name
	return l
}

// Keeper assigns keeper to component key in label
func (l Label) Keeper() Label {
	l.Component(KeeperLabelVal)
	return l
}

// IsKeeper returns whether label is a Keeper
func (l Label) IsKeeper() bool {
	return l[ComponentLabelKey] == KeeperLabelVal
}

// Proxy assigns proxy to component key in label
func (l Label) Proxy() Label {
	l.Component(ProxyLabelVal)
	return l
}

// IsProxy returns whether label is a Proxy
func (l Label) IsProxy() bool {
	return l[ComponentLabelKey] == ProxyLabelVal
}

// Selector gets labels.Selector from label
func (l Label) Selector() (labels.Selector, error) {
	return metav1.LabelSelectorAsSelector(l.LabelSelector())
}

// LabelSelector gets LabelSelector from label
func (l Label) LabelSelector() *metav1.LabelSelector {
	return &metav1.LabelSelector{MatchLabels: l}
}

// Labels converts label to map[string]string
func (l Label) Labels() map[string]string {
	return l
}
