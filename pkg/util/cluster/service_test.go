package cluster

import (
	"testing"

	"github.com/infinivision/citus-operator/pkg/apis/infinivision/v1alpha1"
	api "github.com/infinivision/citus-operator/pkg/apis/infinivision/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestNewHeadlessService(t *testing.T) {
	clus := &api.CitusCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CitusCluster",
			APIVersion: "infinivision.com/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "citus-demo",
			Namespace: "default",
			UID:       types.UID("test"),
		},
		Spec: v1alpha1.CitusClusterSpec{
			Keeper: v1alpha1.KeeperSpec{
				Size: 3,
			},
			Proxy: v1alpha1.ProxySpec{
				MasterPort:  5432,
				StandbyPort: 5433,
				Type:        "ClusterIP",
			},
		},
	}

	svc := NewHeadlessService(clus)
	t.Logf("%#v", svc)
}

func TestNewProxyService(t *testing.T) {
	clus := &api.CitusCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CitusCluster",
			APIVersion: "infinivision.com/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "citus-demo",
			Namespace: "default",
			UID:       types.UID("test"),
		},
		Spec: v1alpha1.CitusClusterSpec{
			Keeper: v1alpha1.KeeperSpec{
				Size: 3,
			},
			Proxy: v1alpha1.ProxySpec{
				MasterPort:  5432,
				StandbyPort: 5433,
				Type:        "ClusterIP",
			},
		},
	}

	svc := NewProxyService(ProxyServiceTypeMaster, clus)
	t.Logf("%#v", svc)
}

func TestNewProxyServiceEndpoints(t *testing.T) {
	clus := &api.CitusCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CitusCluster",
			APIVersion: "infinivision.com/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "citus-demo",
			Namespace: "default",
			UID:       types.UID("test"),
		},
		Spec: v1alpha1.CitusClusterSpec{
			Keeper: v1alpha1.KeeperSpec{
				Size: 3,
			},
			Proxy: v1alpha1.ProxySpec{
				MasterPort:  5432,
				StandbyPort: 5433,
				Type:        "ClusterIP",
			},
		},
	}

	ep := NewProxyServiceEndpoints(ProxyServiceTypeMaster, clus, []string{"192.168.101.174"}, []int32{5432})
	t.Logf("%#v", ep)
}
