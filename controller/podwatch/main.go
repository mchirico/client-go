package main

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	corev1 "k8s.io/api/core/v1"

	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"sigs.k8s.io/controller-runtime/pkg/source"

	"context"
)

var (
	mgr manager.Manager

	log = logf.Log.WithName("controller-examples")
)

func main() {

	ExampleController()
	//ExampleController_unstructured()
}


// This example starts a new Controller named "pod-controller" to Watch Pods and call a no-op Reconciler.
func ExampleController() {
	// mgr is a manager.Manager

	// Create a new Controller that will call the provided Reconciler function in response
	// to events.
	c, err := controller.New("pod-controller", mgr, controller.Options{
		Reconciler: reconcile.Func(func(context.Context, reconcile.Request) (reconcile.Result, error) {
			// Your business logic to implement the API by creating, updating, deleting objects goes here.
			return reconcile.Result{}, nil
		}),
		Log: log,
	})
	if err != nil {
		log.Error(err, "unable to create pod-controller")
		os.Exit(1)
	}


	// Watch for Pod create / update / delete events and call Reconcile
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		log.Error(err, "unable to watch pods")
		os.Exit(1)
	}

	// Start the Controller through the manager.
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "unable to continue running manager")
		os.Exit(1)
	}
}



// This example starts a new Controller named "pod-controller" to Watch Pods with the unstructured object and call a no-op Reconciler.
func ExampleController_unstructured() {
	// mgr is a manager.Manager

	// Create a new Controller that will call the provided Reconciler function in response
	// to events.
	c, err := controller.New("pod-controller", mgr, controller.Options{
		Reconciler: reconcile.Func(func(context.Context, reconcile.Request) (reconcile.Result, error) {
			// Your business logic to implement the API by creating, updating, deleting objects goes here.
			return reconcile.Result{}, nil
		}),
		Log: log,
	})
	if err != nil {
		log.Error(err, "unable to create pod-controller")
		os.Exit(1)
	}

	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Kind:    "Pod",
		Group:   "",
		Version: "v1",
	})
	// Watch for Pod create / update / delete events and call Reconcile
	err = c.Watch(&source.Kind{Type: u}, &handler.EnqueueRequestForObject{})
	if err != nil {
		log.Error(err, "unable to watch pods")
		os.Exit(1)
	}

	// Start the Controller through the manager.
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "unable to continue running manager")
		os.Exit(1)
	}
}

func ExampleNewUnmanaged() {
	// mgr is a manager.Manager

	// Configure creates a new controller but does not add it to the supplied
	// manager.
	c, err := controller.NewUnmanaged("pod-controller", mgr, controller.Options{
		Reconciler: reconcile.Func(func(context.Context, reconcile.Request) (reconcile.Result, error) {
			return reconcile.Result{}, nil
		}),
		Log: log,
	})
	if err != nil {
		log.Error(err, "unable to create pod-controller")
		os.Exit(1)
	}

	if err := c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForObject{}); err != nil {
		log.Error(err, "unable to watch pods")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Start our controller in a goroutine so that we do not block.
	go func() {
		// Block until our controller manager is elected leader. We presume our
		// entire process will terminate if we lose leadership, so we don't need
		// to handle that.
		<-mgr.Elected()

		// Start our controller. This will block until the context is
		// closed, or the controller returns an error.
		if err := c.Start(ctx); err != nil {
			log.Error(err, "cannot run experiment controller")
		}
	}()

	// Stop our controller.
	cancel()
}