package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

// const (
// 	POD_LABEL = "app"
// )

func main() {
	log.Print("Shared Informer app started")
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = "/Users/***/.kube/config"
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
	}

	factory := informers.NewSharedInformerFactory(clientset, 0)
	informer := factory.Core().V1().PersistentVolumes().Informer()
	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		UpdateFunc: onUpdate,
		DeleteFunc: nil,
	})
	go informer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(errors.New("timed out waiting for caches to sync"))
		return
	}
	<-stopper
}

// onAdd is the function executed when the kubernetes informer notified the
// presence of a new kubernetes node in the cluster
func onAdd(obj interface{}) {
	// Cast the obj as pv
	pv := obj.(*corev1.PersistentVolume)
	fmt.Printf("[onAdd] pv %s with claim %s in %s\n", pv.Name, fmt.Sprintf("%s/%s", pv.Spec.ClaimRef.Namespace, pv.Spec.ClaimRef.Name), pv.Status.Phase)
}

func onUpdate(_, obj interface{}) {
	// Cast the obj as pv
	pv := obj.(*corev1.PersistentVolume)
	fmt.Printf("[onUpdate] pv %s with claim %s in %s\n", pv.Name, fmt.Sprintf("%s/%s", pv.Spec.ClaimRef.Namespace, pv.Spec.ClaimRef.Name), pv.Status.Phase)
}
