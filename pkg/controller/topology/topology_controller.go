package topology

import (
	"context"
	"fmt"

	strimziv1beta1 "github.com/adam-cattermole/striot-operator/pkg/apis/strimzi/v1beta1"
	striotv1alpha1 "github.com/adam-cattermole/striot-operator/pkg/apis/striot/v1alpha1"
	"github.com/google/uuid"
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

	// first create all required kafka topics
	for _, topic := range deployments.Topics {
		reqLogger.Info("Creating a new KafkaTopic", "KafkaTopic.Namespace", topic.Namespace, "KafkaTopic.Name", topic.Name)
		err = r.client.Create(context.TODO(), &topic)
		if err != nil {
			reqLogger.Info("Failed to create topic", "err", err)
		}
	}

	for i, d := range deployments.Deployments {

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
			// Deployment created successfully - create relevant service
			service := deployments.Services[i]
			reqLogger.Info("Creating a new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
			err = r.client.Create(context.TODO(), &service)
			if err != nil {
				reqLogger.Info("Failed to create Service", "err", err)
			}
			continue
		} else if err != nil {
			return reconcile.Result{}, err
		}
		// Deployment already exists - don't requeue
		reqLogger.Info("Skip reconcile: Deployment already exists", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
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

// KafkaTopicInfo provides additional info for creating kafkatopics
type KafkaTopicInfo struct {
	Name        string
	ClusterName string
	Namespace   string
	Replicas    int32
	Partitions  int32
}

// Deploy contains all deployments and kafkatopics to deploy
type Deploy struct {
	Deployments []apps.Deployment
	Services    []corev1.Service
	Topics      []strimziv1beta1.KafkaTopic
}

func createDeployments(cr *striotv1alpha1.Topology) *Deploy {
	// name := uuid.New().String()
	replicas := int32(1)
	deploy := Deploy{
		Deployments: []apps.Deployment{},
		Topics:      []strimziv1beta1.KafkaTopic{},
	}

	partitions := []striotv1alpha1.Partition{}
	// build partitions map based on Order parameter
	for _, current := range cr.Spec.Order {
		// find partition next based on order
		for _, p := range cr.Spec.Partitions {
			if p.ID == current {
				partitions = append(partitions, p)
				break
			}
		}
	}

	for i, partition := range partitions {
		envTemp := map[string]string{}
		// setup deploy info
		di := DeployInfo{
			Name:      "striot-node-" + fmt.Sprint(partition.ID),
			Namespace: cr.Namespace,
			Replicas:  replicas,
		}
		// create connection specific info
		envTemp["STRIOT_NODE_NAME"] = di.Name

		// SETUP INGRESS
		switch partition.ConnectType.Ingress {
		case striotv1alpha1.ProtocolTCP:
			envTemp["STRIOT_INGRESS_TYPE"] = striotv1alpha1.ProtocolTCP.String()
			if i == 0 {
				// SOURCE
				envTemp["STRIOT_INGRESS_HOST"] = ""
				envTemp["STRIOT_INGRESS_PORT"] = ""
			} else {
				envTemp["STRIOT_INGRESS_HOST"] = "striot-node-" + fmt.Sprint(partitions[i-1].ID)
				envTemp["STRIOT_INGRESS_PORT"] = "9001"
			}

		case striotv1alpha1.ProtocolKafka:
			envTemp["STRIOT_INGRESS_TYPE"] = striotv1alpha1.ProtocolKafka.String()
			envTemp["STRIOT_INGRESS_HOST"] = "my-cluster-kafka-bootstrap"
			envTemp["STRIOT_INGRESS_PORT"] = "9092"
			envTemp["STRIOT_INGRESS_KAFKA_TOPIC"] = deploy.Topics[len(deploy.Topics)-1].Name
			envTemp["STRIOT_INGRESS_KAFKA_CON_GROUP"] = "striot-con-group"
		case striotv1alpha1.ProtocolMQTT:
		default:
		}

		// SETUP EGRESS
		switch partition.ConnectType.Egress {
		case striotv1alpha1.ProtocolTCP:
			envTemp["STRIOT_EGRESS_TYPE"] = striotv1alpha1.ProtocolTCP.String()
			if len(partitions)-1 == i {
				// SINK
				envTemp["STRIOT_EGRESS_HOST"] = ""
				envTemp["STRIOT_EGRESS_PORT"] = ""
			} else {
				envTemp["STRIOT_EGRESS_HOST"] = "striot-node-" + fmt.Sprint(partitions[i+1].ID)
				envTemp["STRIOT_EGRESS_PORT"] = "9001"
			}
		case striotv1alpha1.ProtocolKafka:
			envTemp["STRIOT_EGRESS_TYPE"] = striotv1alpha1.ProtocolKafka.String()
			kti := KafkaTopicInfo{
				Name:        "striot-topic-" + uuid.New().String(),
				ClusterName: "my-cluster",
				Namespace:   cr.Namespace,
				Replicas:    1,
				Partitions:  64,
			}
			envTemp["STRIOT_EGRESS_HOST"] = "my-cluster-kafka-bootstrap"
			envTemp["STRIOT_EGRESS_PORT"] = "9092"
			envTemp["STRIOT_EGRESS_KAFKA_TOPIC"] = kti.Name
			envTemp["STRIOT_EGRESS_KAFKA_CON_GROUP"] = "none"
			deploy.Topics = append(deploy.Topics, *createKafkaTopic(&kti))
		case striotv1alpha1.ProtocolMQTT:
		default:
		}
		// Set buffer size
		envTemp["STRIOT_CHAN_SIZE"] = "10"
		// Create deployment with info
		env := []corev1.EnvVar{}

		for k, v := range envTemp {
			env = append(env, corev1.EnvVar{Name: k, Value: v})
		}
		labels := map[string]string{
			"app":    di.Name,
			"all":    "striot",
			"striot": "monitor",
		}
		deploy.Deployments = append(deploy.Deployments, *createStriotDeployment(&di, &partition, &env, &labels))
		deploy.Services = append(deploy.Services, *createStriotService(&di, &labels))
	}
	return &deploy
}

func createStriotService(di *DeployInfo, labels *map[string]string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      di.Name,
			Namespace: di.Namespace,
			Labels:    *labels,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:     "striot",
					Port:     9001,
					Protocol: corev1.ProtocolTCP,
				},
				{
					Name:     "prom",
					Port:     8080,
					Protocol: corev1.ProtocolTCP,
				},
			},
			Selector: map[string]string{"app": di.Name},
		},
	}

}

func createStriotDeployment(di *DeployInfo, part *striotv1alpha1.Partition, env *[]corev1.EnvVar, labels *map[string]string) *apps.Deployment {
	return &apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      di.Name,
			Namespace: di.Namespace,
			Labels:    *labels,
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
					Labels: *labels,
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
							Env: *env,
						},
					},
				},
			},
		},
	}
}

func createKafkaTopic(kti *KafkaTopicInfo) *strimziv1beta1.KafkaTopic {
	return &strimziv1beta1.KafkaTopic{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kti.Name,
			Namespace: kti.Namespace,
			Labels: map[string]string{
				"strimzi.io/cluster": kti.ClusterName,
			}},
		Spec: strimziv1beta1.KafkaTopicSpec{
			Replicas:   kti.Replicas,
			Partitions: kti.Partitions,
		},
	}
}
