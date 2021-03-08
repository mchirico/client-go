package main

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"time"
)

func WatchExperiment() {
	const shortDuration = 60 * time.Minute

	d := time.Now().Add(shortDuration)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	// Location of kubeconfig file
	kubeconfig := os.Getenv("HOME") + "/.kube/config"

	// Create a Config (k8s.io/client-go/rest)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// Create an API Clientset (k8s.io/client-go/kubernetes)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Create a CoreV1Client (k8s.io/client-go/kubernetes/typed/core/v1)
	coreV1Client := clientset.CoreV1()
	// Create an AppsV1Client (k8s.io/client-go/kubernetes/typed/apps/v1)




	namespace := "default"
	name := "name3"
	event := &v1.Event{
		TypeMeta:            metav1.TypeMeta{},
		ObjectMeta:          metav1.ObjectMeta{Namespace: namespace,Name: name},
		InvolvedObject:      v1.ObjectReference{},
		Reason:              "Some silly reason",
		Message:             "Message... some random message here... ",
		Source:              v1.EventSource{},
		FirstTimestamp:      metav1.Time{time.Now()},
		LastTimestamp:       metav1.Time{time.Now().Add(-78 * time.Minute)},
		Count:               10,
		Type:                "Special",
		EventTime:           metav1.MicroTime{},
		Series:              nil,
		Action:              "Dropped jaw",
		Related:             nil,
		ReportingController: "A Giant Panda",
		ReportingInstance:   "",
	}
	_ , err = coreV1Client.Events(namespace).Create(ctx, event , metav1.CreateOptions{})

	event.Action = "Smile"
	coreV1Client.Events(namespace).Update(ctx,event, metav1.UpdateOptions{})



	if err != nil {
		coreV1Client.Events(namespace).Delete(ctx, name , metav1.DeleteOptions{})
		log.Fatal(err.Error())
	}

	go func() {
		time.Sleep(7 * time.Second)
		err = coreV1Client.Events(namespace).Delete(ctx, name , metav1.DeleteOptions{})
	}()

/*
By hand

Need to get a list?

   k describe events

  k get events -A --field-selector involvedObject.kind=Environment



   watch, err := coreV1Client.Events("").Watch(ctx, metav1.ListOptions{
   		FieldSelector: "involvedObject.name=environment-sample",
   	})

Other options

   watch, err := coreV1Client.Events("").Watch(ctx, metav1.ListOptions{
      		FieldSelector: "involvedObject.name=environment-sample",
      	})

 */

	watch, err := coreV1Client.Events("").Watch(ctx, metav1.ListOptions{
		FieldSelector: "involvedObject.name=environment-sample",

	})

	//watch, err := coreV1Client.Pods("").Watch(ctx, metav1.ListOptions{})

	if err != nil {
		log.Fatal(err.Error())
	}
	go func() {
		for event := range watch.ResultChan() {
			fmt.Printf("Type: %v\n", event.Type)
			fmt.Printf("Event: %v\n", event)

		}
	}()
	time.Sleep(60 * time.Minute)

}

func main() {

	// Location of kubeconfig file
	kubeconfig := os.Getenv("HOME") + "/.kube/config"

	// Create a Config (k8s.io/client-go/rest)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// Create an API Clientset (k8s.io/client-go/kubernetes)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Create a CoreV1Client (k8s.io/client-go/kubernetes/typed/core/v1)
	coreV1Client := clientset.CoreV1()
	// Create an AppsV1Client (k8s.io/client-go/kubernetes/typed/apps/v1)
	appsV1Client := clientset.AppsV1()

	//-------------------------------------------------------------------------//
	// List pods (all namespaces)
	//-------------------------------------------------------------------------//

	// Get a *PodList (k8s.io/api/core/v1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pods, err := coreV1Client.Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// List each Pod (k8s.io/api/core/v1)
	for i, pod := range pods.Items {
		fmt.Printf("Pod %d: %s\n", i+1, pod.ObjectMeta.Name)
	}

	//-------------------------------------------------------------------------//
	// List nodes
	//-------------------------------------------------------------------------//

	// Get a *NodeList (k8s.io/api/core/v1)
	nodes, err := coreV1Client.Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// For each Node (k8s.io/api/core/v1)
	for i, node := range nodes.Items {
		fmt.Printf("Node %d: %s\n", i+1, node.ObjectMeta.Name)
	}

	//-------------------------------------------------------------------------//
	// List deployments (all namespaces)
	//-------------------------------------------------------------------------//

	// Get a *DeploymentList (k8s.io/api/apps/v1)
	deployments, err := appsV1Client.Deployments("").List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// For each Deployment (k8s.io/api/apps/v1)
	for i, deployment := range deployments.Items {
		fmt.Printf("Deployment %d: %s\n", i+1, deployment.ObjectMeta.Name)
	}


	//-------------------------------------------------------------------------//
	// Create Namespace
	//-------------------------------------------------------------------------//


	nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "ns-bubble"}}
	_, err = clientset.CoreV1().Namespaces().Create(context.TODO(),nsSpec, metav1.CreateOptions{})
	if err != nil {
		fmt.Errorf("err:%w",err)
	}


	//-------------------------------------------------------------------------//
	// CRD
	//-------------------------------------------------------------------------//


	//nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "ns-bubble"}}
	//_, err = clientset.CoreV1().Namespaces().Create(context.TODO(),nsSpec, metav1.CreateOptions{})
	//if err != nil {
	//	fmt.Errorf("err:%w",err)
	//}









	fmt.Printf("\n\nGetting Ready to run watch:\nCtl-c to kill\n\n")
	time.Sleep(2 * time.Second)
	fmt.Printf("")

	WatchExperiment()

}
