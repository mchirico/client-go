package main

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func main() {

	r := reconcile.Func(func(_ context.Context, o reconcile.Request) (reconcile.Result, error) {
		// Create your business logic to create, update, delete objects here.
		fmt.Printf("Name: %s, Namespace: %s", o.Name, o.Namespace)

		return reconcile.Result{Requeue: true}, nil
	})

	res, err := r.Reconcile(context.Background(), reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "test"}})
	if err != nil || res.Requeue || res.RequeueAfter != time.Duration(0) {
		fmt.Printf("\ngot requeue request: %v, %v\n", err, res)
	}

}

