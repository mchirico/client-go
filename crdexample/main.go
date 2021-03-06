/*
Copyright 2017 The Kubernetes Authors.

Ref:
https://raw.githubusercontent.com/kubernetes/client-go/master/examples/create-update-delete-deployment/main.go
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cachev1alpha1 "github.com/mchirico/memcached-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	const (
		kind      = "Memcached"
		namespace = "default"
		name      = "memcached-sample"

		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	ctx := context.Background()

	key := types.NamespacedName{
		Name:      name + "-abcd",
		Namespace: namespace,
	}

	crd := &cachev1alpha1.Memcached{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Memcached",
			APIVersion: "cache.example.com/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      key.Name,
			Namespace: key.Namespace,
		},
		Status: cachev1alpha1.MemcachedStatus{},
	}

	resultCRD := &cachev1alpha1.Memcached{
		ObjectMeta: metav1.ObjectMeta{
			Name:      key.Name,
			Namespace: key.Namespace,
		},
		Status: cachev1alpha1.MemcachedStatus{},
	}

	crd.Kind = kind
	crd.Namespace = namespace
	crd.Name = name
	crd.Spec = cachev1alpha1.MemcachedSpec{}
	crd.Spec.Size = 7

	var k8sClient client.Client

	err = cachev1alpha1.AddToScheme(scheme.Scheme)
	//Scheme: scheme.Scheme
	k8sClient, err = client.New(config, client.Options{Scheme: scheme.Scheme})
	fmt.Printf("Load CRD\n")
	prompt()
	err = k8sClient.Create(ctx, crd)
	if err != nil {
		fmt.Errorf("error: %w\n", err)
	}
	for i := 0; i < 4; i++ {
		err = k8sClient.Get(ctx, key, resultCRD)
		if err != nil {
			fmt.Errorf("error: %w\n", err)
		}
		fmt.Printf("crd.Status.Nodes: %d\n", len(crd.Status.Nodes))
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("Now Update CRD to 2\n")
	prompt()

	crd.Spec.Size = 2
	err = k8sClient.Update(ctx, crd)
	if err != nil {
		fmt.Errorf("error: %w\n", err)
	}
	for i := 0; i < 7; i++ {
		err = k8sClient.Get(ctx, key, resultCRD)
		if err != nil {
			fmt.Errorf("error: %w\n", err)
		}
		fmt.Printf("crd.Status.Nodes: %d\n", len(crd.Status.Nodes))
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("Now change to 12")
	prompt()
	crd.Spec.Size = 12
	err = k8sClient.Update(ctx, crd)
	if err != nil {
		fmt.Errorf("error: %w\n", err)
	}
	// We won't get change immediately?
	for i := 0; i < 7; i++ {
		err = k8sClient.Get(ctx, key, resultCRD)
		if err != nil {
			fmt.Errorf("error: %w\n", err)
		}
		fmt.Printf("crd.Status.Nodes: %d\n", len(crd.Status.Nodes))
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("Now Delete CRD\n")
	prompt()
	err = k8sClient.Delete(ctx, crd)
	if err != nil {
		fmt.Errorf("error: %w\n", err)
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	// Create Deployment
	//fmt.Println("Creating crd deployment...")
	//result, err := deploymentsClient.Create(context.TODO(), memcached, metav1.CreateOptions{})

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	// Update Deployment
	prompt()
	fmt.Println("Updating deployment...")
	//    You have two options to Update() this Deployment:
	//
	//    1. Modify the "deployment" variable and call: Update(deployment).
	//       This works like the "kubectl replace" command and it overwrites/loses changes
	//       made by other clients between you Create() and Update() the object.
	//    2. Modify the "result" returned by Get() and retry Update(result) until
	//       you no longer get a conflict error. This way, you can preserve changes made
	//       by other clients between Create() and Update(). This is implemented below
	//			 using the retry utility package included with client-go. (RECOMMENDED)
	//
	// More Info:
	// https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := deploymentsClient.Get(context.TODO(), "demo-deployment", metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
		}

		result.Spec.Replicas = int32Ptr(1)                           // reduce replica count
		result.Spec.Template.Spec.Containers[0].Image = "nginx:1.13" // change nginx version
		_, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Println("Updated deployment...")

	// List Deployments
	prompt()
	fmt.Printf("Listing deployments in namespace %q:\n", apiv1.NamespaceDefault)
	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}

	// Delete Deployment
	prompt()
	fmt.Println("Deleting deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete(context.TODO(), "demo-deployment", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted deployment.")
}

func prompt() {
	fmt.Printf("-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}

func int32Ptr(i int32) *int32 { return &i }
