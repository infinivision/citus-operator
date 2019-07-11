package cituscluster

import (
	"context"

	infinivisionv1alpha1 "github.com/infinivision/citus-operator/pkg/apis/infinivision/v1alpha1"
	"github.com/infinivision/citus-operator/pkg/util"
	clusterutil "github.com/infinivision/citus-operator/pkg/util/cluster"
	apps "k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_cituscluster")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new CitusCluster Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCitusCluster{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("cituscluster-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource CitusCluster
	err = c.Watch(&source.Kind{Type: &infinivisionv1alpha1.CitusCluster{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner CitusCluster
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &infinivisionv1alpha1.CitusCluster{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileCitusCluster implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileCitusCluster{}

// ReconcileCitusCluster reconciles a CitusCluster object
type ReconcileCitusCluster struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a CitusCluster object and makes changes based on the state read
// and what is in the CitusCluster.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCitusCluster) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling CitusCluster")

	// Fetch the CitusCluster instance
	instance := &infinivisionv1alpha1.CitusCluster{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Check the stolon keeper statefulset
	ss := &apps.StatefulSet{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, ss)
	if err != nil && errors.IsNotFound(err) {
		// TODO: should initialize the stolon cluster
		if err := r.initCluster(instance); err != nil {
			return reconcile.Result{}, err
		}
		reqLogger.Info("Creating a new Statefulset for the stolon keeper", "Namespace", instance.Namespace, "Name", instance.Name)
		ss = clusterutil.NewKeeperStatefulset(instance)
		// Set CitusCluster instance as the owner and controller
		if err = controllerutil.SetControllerReference(instance, ss, r.scheme); err != nil {
			return reconcile.Result{}, err
		}
		err = r.client.Create(context.TODO(), ss)
		if err != nil {
			return reconcile.Result{}, err
		}

		// The stolon keeper statefulset created successfully - don't requeue
		return reconcile.Result{}, nil
	}

	keeperSpec := instance.Spec.Keeper
	needUpdate := false

	// Scale up
	if keeperSpec.Size > *ss.Spec.Replicas {
		ss.Spec.Replicas = &keeperSpec.Size
		reqLogger.Info("Scale up the stolon keeper to ", keeperSpec.Size)
		needUpdate = true
	} else if keeperSpec.Size < *ss.Spec.Replicas {
		// Scale down
		ss.Spec.Replicas = &keeperSpec.Size
		reqLogger.Info("Scale down the stolon keeper to ", keeperSpec.Size)
		needUpdate = true
	}

	ctn := &ss.Spec.Template.Spec.Containers[0]
	// Upgrade
	if len(keeperSpec.Image) > 0 && keeperSpec.Image != ctn.Image {
		reqLogger.Info("Upgrade the stolon keeper image to ", keeperSpec.Image)
		ctn.Image = keeperSpec.Image
		needUpdate = true
	}

	// Update request resources
	// cpuSpec := resource.MustParse(keeperSpec.Requests.CPU)
	// memSpec := resource.MustParse(keeperSpec.Requests.Memory)
	// cpu := ctn.Resources.Limits.Cpu()
	// mem := ctn.Resources.Limits.Memory()
	// if cpuSpec.Cmp(*cpu) != 0 || memSpec.Cmp(*mem) != 0 {
	// 	reqLogger.Info("Update the stolon keeper resource: cpu=", keeperSpec.Requests.CPU, "memory=", keeperSpec.Requests.Memory)
	// 	// ctn.Resources.Limits.
	// 	needUpdate = true
	// }

	if needUpdate {
		err = r.client.Update(context.TODO(), ss)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

// initCluster init the stolon cluster
func (r *ReconcileCitusCluster) initCluster(clus *infinivisionv1alpha1.CitusCluster) error {
	// stolonctl --cluster-name=kube-stolon --store-backend=kubernetes --kube-resource-kind=configmap init
	_, result, err := util.ExecCommand("/bin/stolon-ctl", []string{"--cluster-name=" + clus.ClusterName,
		"--store-backend=kubernetes", "--kube-resource-kind=configmap", "init"})
	log.V(3).Info("init cluster result:", result)
	if err != nil {
		log.Error(err, "init cluster failed.")
		return err
	}

	return nil
}
