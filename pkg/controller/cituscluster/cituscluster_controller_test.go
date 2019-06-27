package cituscluster

import (
	"testing"

	"github.com/infinivision/citus-operator/pkg/apis/infinivision/v1alpha1"
	api "github.com/infinivision/citus-operator/pkg/apis/infinivision/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestCitusCluster(t *testing.T) {
	cluster := &api.CitusCluster{
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
				Size: 3,
			},
		},
	}

	t.Logf("%+v", cluster)
}
