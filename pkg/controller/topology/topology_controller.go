package topology

import (
	"context"

	striotv1alpha1 "github.com/adam-cattermole/striot-operator/pkg/apis/striot/v1alpha1"
	apps "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_topology")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Topology Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileTopology{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("topology-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Topology
	err = c.Watch(&source.Kind{Type: &striotv1alpha1.Topology{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Topology
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &striotv1alpha1.Topology{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileTopology implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileTopology{}

// ReconcileTopology reconciles a Topology object
type ReconcileTopology struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Topology object and makes changes based on the state read
// and what is in the Topology.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileTopology) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Topology")

	// Fetch the Topology instance
	instance := &striotv1alpha1.Topology{}
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

	// // Define a new Pod object
	// pod := newPodForCR(instance)

	// // Set Topology instance as the owner and controller
	// if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
	// 	return reconcile.Result{}, err
	// }

	// // Check if this Pod already exists
	// found := &corev1.Pod{}
	// err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	// if err != nil && errors.IsNotFound(err) {
	// 	reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
	// 	err = r.client.Create(context.TODO(), pod)
	// 	if err != nil {
	// 		return reconcile.Result{}, err
	// 	}

	// 	// Pod created successfully - don't requeue
	// 	return reconcile.Result{}, nil
	// } else if err != nil {
	// 	return reconcile.Result{}, err
	// }

	// // Pod already exists - don't requeue
	// reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	// return reconcile.Result{}, nil

	// Define a new Deployment Object
	deployments := createDeployments(instance)
	for _, d := range *deployments {

		// Set Topology instance as the owner and controller
		if err = controllerutil.SetControllerReference(instance, &d, r.scheme); err != nil {
			return reconcile.Result{}, err
		}

		// Check if this deployment already exists
		found := &apps.Deployment{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: d.Name, Namespace: d.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", d.Namespace, "Deployment.Name", d.Name)
			err = r.client.Create(context.TODO(), &d)
			if err != nil {
				return reconcile.Result{}, err
			}
			// Deployment created successfully - don't requeue
			return reconcile.Result{}, nil
		} else if err != nil {
			return reconcile.Result{}, err
		}

		// Deployment already exists - don't requeue
		reqLogger.Info("Skip reconcile: Deployment already exists", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
		return reconcile.Result{}, nil
	}
	return reconcile.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *striotv1alpha1.Topology) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}

// DeployInfo provides additional deployment info
type DeployInfo struct {
	Name      string
	Namespace string
	Replicas  int32
}

func createDeployments(cr *striotv1alpha1.Topology) *[]apps.Deployment {
	// name := uuid.New().String()
	replicas := int32(1)
	deploy := []apps.Deployment{}

	for _, p := range cr.Spec.Partitions {
		di := DeployInfo{
			Name:      "striot-node-" + string(p.ID),
			Namespace: cr.Namespace,
			Replicas:  replicas,
		}
		deploy = append(deploy, *createStriotDeployment(&di, &p))
	}

	return &deploy
}

func createStriotDeployment(di *DeployInfo, part *striotv1alpha1.Partition) *apps.Deployment {

	labels := map[string]string{
		"app":    di.Name,
		"all":    "striot",
		"striot": "monitor",
	}
	env := *generateEnv(part)

	return &apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      di.Name,
			Namespace: di.Namespace,
			Labels:    labels,
		},
		Spec: apps.DeploymentSpec{
			Replicas: &di.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": di.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  di.Name,
							Image: part.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "striot",
									ContainerPort: 9001,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          "prom",
									ContainerPort: 8080,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							Env: env,
						},
					},
				},
			},
		},
	}
}

func generateEnv(cr *striotv1alpha1.Partition) *[]corev1.EnvVar {

	test := map[string]string{
		"STRIOT_NODE_NAME":              "striot-src",
		"STRIOT_INGRESS_TYPE":           striotv1alpha1.ProtocolTCP.String(),
		"STRIOT_INGRESS_HOST":           "",
		"STRIOT_INGRESS_PORT":           "",
		"STRIOT_EGRESS_TYPE":            striotv1alpha1.ProtocolKafka.String(),
		"STRIOT_EGRESS_HOST":            "my-cluster-kafka-bootstrap",
		"STRIOT_EGRESS_PORT":            "9092",
		"STRIOT_EGRESS_KAFKA_TOPIC":     "striot-queue",
		"STRIOT_EGRESS_KAFKA_CON_GROUP": "none",
		"STRIOT_CHAN_SIZE":              "10",
	}
	env := []corev1.EnvVar{}

	for k, v := range test {
		env = append(env, corev1.EnvVar{Name: k, Value: v})
	}

	return &env
}

// TODO:
//		- Create a function that iterates over Partitions and creates all necessary containers/deploys
// 		- Find a way to give default values to the containers depending on the ingress/egress types
//		- Find out how to interface with the strimzi operator to create the required topics
//		- Check how to create unique hashes in Go so that we can name the topics appropriately - DONE
