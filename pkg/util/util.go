package util

import (
	"github.com/golang/glog"
	"github.com/infinivision/citus-operator/pkg/apis/infinivision/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// ResourceRequirement creates ResourceRequirements for MemberSpec
// Optionally pass in a default value
func ResourceRequirement(spec v1alpha1.ContainerSpec, defaultRequests ...corev1.ResourceRequirements) corev1.ResourceRequirements {
	rr := corev1.ResourceRequirements{}
	if len(defaultRequests) > 0 {
		defaultRequest := defaultRequests[0]
		rr.Requests = make(map[corev1.ResourceName]resource.Quantity)
		rr.Requests[corev1.ResourceCPU] = defaultRequest.Requests[corev1.ResourceCPU]
		rr.Requests[corev1.ResourceMemory] = defaultRequest.Requests[corev1.ResourceMemory]
		rr.Limits = make(map[corev1.ResourceName]resource.Quantity)
		rr.Limits[corev1.ResourceCPU] = defaultRequest.Limits[corev1.ResourceCPU]
		rr.Limits[corev1.ResourceMemory] = defaultRequest.Limits[corev1.ResourceMemory]
	}
	if spec.Requests != nil {
		if rr.Requests == nil {
			rr.Requests = make(map[corev1.ResourceName]resource.Quantity)
		}
		if spec.Requests.CPU != "" {
			if q, err := resource.ParseQuantity(spec.Requests.CPU); err != nil {
				glog.Errorf("failed to parse CPU resource %s to quantity: %v", spec.Requests.CPU, err)
			} else {
				rr.Requests[corev1.ResourceCPU] = q
			}
		}
		if spec.Requests.Memory != "" {
			if q, err := resource.ParseQuantity(spec.Requests.Memory); err != nil {
				glog.Errorf("failed to parse memory resource %s to quantity: %v", spec.Requests.Memory, err)
			} else {
				rr.Requests[corev1.ResourceMemory] = q
			}
		}
	}
	if spec.Limits != nil {
		if rr.Limits == nil {
			rr.Limits = make(map[corev1.ResourceName]resource.Quantity)
		}
		if spec.Limits.CPU != "" {
			if q, err := resource.ParseQuantity(spec.Limits.CPU); err != nil {
				glog.Errorf("failed to parse CPU resource %s to quantity: %v", spec.Limits.CPU, err)
			} else {
				rr.Limits[corev1.ResourceCPU] = q
			}
		}
		if spec.Limits.Memory != "" {
			if q, err := resource.ParseQuantity(spec.Limits.Memory); err != nil {
				glog.Errorf("failed to parse memory resource %s to quantity: %v", spec.Limits.Memory, err)
			} else {
				rr.Limits[corev1.ResourceMemory] = q
			}
		}
	}
	return rr
}
