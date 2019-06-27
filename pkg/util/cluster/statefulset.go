package cluster

import (
	"fmt"

	api "github.com/infinivision/citus-operator/pkg/apis/infinivision/v1alpha1"
	"github.com/infinivision/citus-operator/pkg/util"
	"github.com/infinivision/citus-operator/pkg/util/label"
	apps "k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultKeeperStorageSize = "5Gi"
)

// GetOwnerRef returns CitusCluster's OwnerReference
func GetOwnerRef(cc *api.CitusCluster) metav1.OwnerReference {
	controller := true
	blockOwnerDeletion := true
	return metav1.OwnerReference{
		APIVersion:         cc.APIVersion,
		Kind:               cc.Kind,
		Name:               cc.GetName(),
		UID:                cc.GetUID(),
		Controller:         &controller,
		BlockOwnerDeletion: &blockOwnerDeletion,
	}
}

func NewKeeperStatefulset(cc *api.CitusCluster) *apps.StatefulSet {
	instanceName := cc.GetLabels()[label.InstanceLabelKey]
	keeperLabel := label.New().Instance(instanceName).Keeper()
	stolonDataVol := "data"
	volMounts := []corev1.VolumeMount{
		{Name: stolonDataVol, MountPath: "/stolon-data"},
		{Name: "stolon", ReadOnly: false, MountPath: "/etc/secrets/stolon"},
	}
	vols := []corev1.Volume{
		{Name: "stolon", VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: "stolon",
			}},
		},
	}
	annos := map[string]string{
		"pod.alpha.kubernetes.io/initialized": "true",
		"prometheus.io/scrape":                "true",
		"prometheus.io/port/keeper":           "8080",
		"prometheus.io/port/sentinel":         "8081",
	}
	q, _ := resource.ParseQuantity(defaultKeeperStorageSize)
	if cc.Spec.Keeper.Requests != nil {
		size := cc.Spec.Keeper.Requests.Storage
		var err error
		q, err = resource.ParseQuantity(size)
		if err != nil {
			panic(fmt.Errorf("cant' get storage size: %s for CitusCluster: %s/%s, %v", size, cc.Namespace, cc.Name, err))
		}
	}
	cmd := []string{
		"/bin/bash",
		"-ec",
		`IFS='-' read -ra ADDR <<< "$(hostname)"
		export STKEEPER_UID=keeper"${ADDR[-1]}"
		export POD_IP=$(hostname -i)
		export STKEEPER_PG_LISTEN_ADDRESS=$POD_IP
		export STOLON_DATA=/stolon-data
		chown stolon:stolon $STOLON_DATA
		exec gosu stolon stolon-sentinel &;
		exec gosu stolon stolon-keeper --data-dir $STOLON_DATA&`,
	}
	ss := &apps.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:            cc.Name,
			Namespace:       cc.Namespace,
			Labels:          keeperLabel.Labels(),
			OwnerReferences: []metav1.OwnerReference{GetOwnerRef(cc)},
		},
		Spec: apps.StatefulSetSpec{
			ServiceName: cc.Name,
			Replicas:    func() *int32 { r := cc.Spec.Keeper.Size; return &r }(),
			Selector:    keeperLabel.LabelSelector(),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      keeperLabel.Labels(),
					Annotations: annos,
				},
				Spec: corev1.PodSpec{

					// SchedulerName: cc.Spec.SchedulerName,
					Containers: []corev1.Container{
						{
							Name:            "stolon-keeper",
							Image:           cc.Spec.Keeper.Image,
							Command:         cmd,
							ImagePullPolicy: cc.Spec.Keeper.ImagePullPolicy,
							Ports: []corev1.ContainerPort{
								{
									Name:          "keeperPort",
									ContainerPort: int32(5432),
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          "keeperMetricsPort",
									ContainerPort: int32(8080),
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          "sentinelMetricsPort",
									ContainerPort: int32(8081),
									Protocol:      corev1.ProtocolTCP,
								},
							},
							VolumeMounts: volMounts,
							Resources:    util.ResourceRequirement(cc.Spec.Keeper.ContainerSpec),
							Env: []corev1.EnvVar{
								{
									Name: "NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
								{
									Name: "POD_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
								{
									Name:  "STKEEPER_CLUSTER_NAME",
									Value: cc.Name,
								},
								{
									Name:  "STKEEPER_STORE_BACKEND",
									Value: "kubernetes",
								},
								{
									Name:  "STKEEPER_KUBE_RESOURCE_KIND",
									Value: "configmap",
								},
								{
									Name:  "STKEEPER_PG_REPL_USERNAME",
									Value: "repluser",
								},
								{
									Name:  "STKEEPER_PG_REPL_PASSWORD",
									Value: "replpassword",
								},
								{
									Name:  "STKEEPER_PG_SU_USERNAME",
									Value: "stolon",
								},
								{
									Name:  "STKEEPER_PG_SU_PASSWORDFILE",
									Value: "/etc/secrets/stolon/password",
								},
								{
									Name:  "STKEEPER_METRICS_LISTEN_ADDRESS",
									Value: "0.0.0.0:8080",
								},
								{
									Name:  "STSENTINEL_CLUSTER_NAME",
									Value: cc.Name,
								},
								{
									Name:  "STSENTINEL_STORE_BACKEND",
									Value: "kubernetes",
								},
								{
									Name:  "STSENTINEL_KUBE_RESOURCE_KIND",
									Value: "configmap",
								},
								{
									Name:  "STSENTINEL_METRICS_LISTEN_ADDRESS",
									Value: "0.0.0.0:8081",
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyAlways,
					// Tolerations:   cc.Spec.Store.Tolerations,
					Volumes: vols,
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				VolumeClaimTemplates(q, stolonDataVol, func() *string { s := cc.Spec.Keeper.StorageClassName; return &s }()),
			},
			PodManagementPolicy: apps.ParallelPodManagement,
			UpdateStrategy: apps.StatefulSetUpdateStrategy{
				Type: apps.RollingUpdateStatefulSetStrategyType,
				RollingUpdate: &apps.RollingUpdateStatefulSetStrategy{
					Partition: func() *int32 { r := cc.Spec.Keeper.Size; return &r }(),
				},
			},
		},
	}

	return ss
}

func VolumeClaimTemplates(q resource.Quantity, metaName string, storageClassName *string) corev1.PersistentVolumeClaim {
	return corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{Name: metaName},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			StorageClassName: storageClassName,
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: q,
				},
			},
		},
	}
}
