package main

import (
    "os"

    "sigs.k8s.io/controller-runtime/pkg/controller"
    "sigs.k8s.io/controller-runtime/pkg/manager"
    "sigs.k8s.io/controller-runtime/pkg/reconcile"
    logf "sigs.k8s.io/controller-runtime/pkg/log"
    "context"
)

var (
    mgr manager.Manager

    log = logf.Log.WithName("controller-examples")
)




func main() {
    _, err := controller.New("pod-controller", mgr, controller.Options{
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
}
