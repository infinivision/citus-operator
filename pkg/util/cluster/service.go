package cluster

import (
	"fmt"

	api "github.com/infinivision/citus-operator/pkg/apis/infinivision/v1alpha1"
	"github.com/infinivision/citus-operator/pkg/util/label"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	// KeeperPortName stolon proxy service name
	KeeperPortName = "keeper"
	// ProxyServiceTypeMaster proxy master service
	ProxyServiceTypeMaster ProxyServiceType = "master"
	// ProxyServiceTypeStandby proxy standby service
	ProxyServiceTypeStandby ProxyServiceType = "standby"
)

// ProxyServiceType proxy service type, master or standby
type ProxyServiceType string

// ProxyServiceName generate the cluster proxy name
func ProxyServiceName(clusterName string, proxyType ProxyServiceType) string {
	return fmt.Sprintf("%s-proxy-%s", clusterName, proxyType)
}

// PortToName generate port name
func PortToName(port int32) string {
	return fmt.Sprintf("port%d", port)
}

// NewHeadlessService new a citus cluster headless service
func NewHeadlessService(cc *api.CitusCluster) *corev1.Service {
	instanceName := cc.GetLabels()[label.InstanceLabelKey]
	keeperLabel := label.New().Instance(instanceName).Keeper()
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            cc.Name,
			Namespace:       cc.Namespace,
			Labels:          keeperLabel.Labels(),
			OwnerReferences: []metav1.OwnerReference{GetOwnerRef(cc)},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       KeeperPortName,
					Protocol:   corev1.ProtocolTCP,
					Port:       int32(KeeperPort),
					TargetPort: intstr.FromInt(KeeperPort),
				},
			},
			ClusterIP: "None",
			Type:      corev1.ServiceTypeClusterIP,
			Selector:  keeperLabel,
		},
	}
}

// NewProxyService new a citus cluster proxy service(master or standby)
func NewProxyService(proxyType ProxyServiceType, cc *api.CitusCluster) *corev1.Service {
	instanceName := cc.GetLabels()[label.InstanceLabelKey]
	keeperLabel := label.New().Instance(instanceName).Keeper()
	serviceType := corev1.ServiceTypeClusterIP
	if cc.Spec.Proxy.Type == string(corev1.ServiceTypeNodePort) {
		serviceType = corev1.ServiceTypeNodePort
	}
	servicePort := cc.Spec.Proxy.StandbyPort
	if proxyType == ProxyServiceTypeMaster {
		servicePort = cc.Spec.Proxy.MasterPort
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            ProxyServiceName(cc.Name, proxyType),
			Namespace:       cc.Namespace,
			Labels:          keeperLabel.Labels(),
			OwnerReferences: []metav1.OwnerReference{GetOwnerRef(cc)},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       PortToName(servicePort),
					Protocol:   corev1.ProtocolTCP,
					Port:       int32(servicePort),
					TargetPort: intstr.FromInt(KeeperPort),
				},
			},
			Type: serviceType,
		},
	}
}

// NewProxyServiceEndpoints new proxy service endpoints related to the proxy service
func NewProxyServiceEndpoints(proxyType ProxyServiceType, cc *api.CitusCluster, addresses []string, ports []int32) *corev1.Endpoints {
	addrs := make([]corev1.EndpointAddress, 0)
	for _, v := range addresses {
		addrs = append(addrs, corev1.EndpointAddress{IP: v})
	}
	pts := make([]corev1.EndpointPort, 0)
	for _, v := range ports {
		pts = append(pts, corev1.EndpointPort{Name: PortToName(v), Port: v})
	}

	return &corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:            ProxyServiceName(cc.Name, proxyType),
			Namespace:       cc.Namespace,
			OwnerReferences: []metav1.OwnerReference{GetOwnerRef(cc)},
		},
		Subsets: []corev1.EndpointSubset{
			corev1.EndpointSubset{
				Addresses: addrs,
				Ports:     pts,
			},
		},
	}
}
