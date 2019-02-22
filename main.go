package main

import (
	"flag"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"test.local/crd/api/types/v1alpha1"
	clientset "test.local/crd/clientset/v1alpha1"
	"time"
)

var (
	kubeconfig string
	masterurl string
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "kubeconfig path")
}

func main() {
	flag.Parse()

	var err error

	if kubeconfig == "" {
		panic("kubeconfig path must required")
	}

	var config *rest.Config
	if config, err = clientcmd.BuildConfigFromFlags(masterurl, kubeconfig); err != nil {
		panic(err)
	}

	v1alpha1.AddToScheme(runtime.NewScheme())

	var myClient *clientset.MyV1Aplpha1Client
	if myClient, err = clientset.NewClientForConfig(config); err != nil {
		panic(err)
	}

	var kubeClient *kubernetes.Clientset
	if kubeClient, err = kubernetes.NewForConfig(config); err != nil {
		panic(err)
	}

	for {
		func() {
			var err error

			defer time.Sleep(time.Second * 3)

			var redis v1alpha1.Redis
			err = myClient.
				RestClient.Get().
				Namespace("default").
				Resource("redis").
				Name("my-redis").
				Do().
				Into(&redis)
			if err != nil {
				log.Println("crd:", err)
				return
			}

			deployment, err := kubeClient.AppsV1().Deployments(redis.Namespace).Get(redis.Spec.DeploymentName, metav1.GetOptions{})
			if errors.IsNotFound(err) {
				deployment, err = kubeClient.AppsV1().Deployments(redis.Namespace).Create(newRedisDeployment(&redis))
			}
			if err != nil {
				log.Println(redis.Name, deployment.Name, err)
				return
			}

			if redis.Spec.Replicas != nil && *redis.Spec.Replicas != *deployment.Spec.Replicas {
				deployment, err = kubeClient.AppsV1().Deployments(redis.Namespace).Update(newRedisDeployment(&redis))
				if err != nil {
					log.Println(err)
					return
				}
			}
		}()
	}
}

func newRedisDeployment(myRedis *v1alpha1.Redis) *appsv1.Deployment{
	labels := map[string]string {
		"app": "redis",
		"controller": myRedis.Name,
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: myRedis.Spec.DeploymentName,
			Namespace: myRedis.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(myRedis, schema.GroupVersionKind{
					Group: v1alpha1.GroupName,
					Version: v1alpha1.GroupVersion,
					Kind: "Redis",
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: myRedis.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "redis",
							Image: "redis:latest",
						},
					},
				},
			},
		},
	}
}
