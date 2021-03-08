package main

import (
    "os"

    "sigs.k8s.io/controller-runtime/pkg/client/config"
    "sigs.k8s.io/controller-runtime/pkg/manager"
    logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var

// NB: don't call SetLogger in init(), or else you'll mess up logging in the main suite.
log = logf.Log.WithName("manager-examples")

func main() {
    cfg, err := config.GetConfig()
    if err != nil {
        log.Error(err, "unable to get kubeconfig")
        os.Exit(1)
    }

    mgr, err := manager.New(cfg, manager.Options{})
    if err != nil {
        log.Error(err, "unable to set up manager")
        os.Exit(1)
    }
    log.Info("created manager", "manager", mgr)
}
