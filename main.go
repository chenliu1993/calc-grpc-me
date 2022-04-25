package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

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

	informersSyncd := make([]cache.InformerSynced, 0)

	factory := informers.NewSharedInformerFactory(clientset, time.Hour*24)
	pvinformer := factory.Core().V1().PersistentVolumes().Informer()
	informersSyncd = append(informersSyncd, pvinformer.HasSynced)
	pvcinformer := factory.Core().V1().PersistentVolumeClaims().Informer()
	informersSyncd = append(informersSyncd, pvcinformer.HasSynced)

	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()
	pvinformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onPVAdd,
		UpdateFunc: onPVUpdate,
		DeleteFunc: onPVDelete,
	})
	pvcinformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onPVCAdd,
		UpdateFunc: onPVCUpdate,
		DeleteFunc: onPVCDelete,
	})
	go factory.Start(stopper)
	if !cache.WaitForCacheSync(stopper, informersSyncd...) {
		runtime.HandleError(errors.New("timed out waiting for caches to sync"))
		return
	}
	<-stopper
}

// onAdd is the function executed when the kubernetes informer notified the
// presence of a new kubernetes node in the cluster
func onPVAdd(obj interface{}) {
	// Cast the obj as pv
	pv := obj.(*corev1.PersistentVolume)
	fmt.Printf("[onPVAdd] pv %s with claim %s in %s\n", pv.Name, fmt.Sprintf("%s/%s", pv.Spec.ClaimRef.Namespace, pv.Spec.ClaimRef.Name), pv.Status.Phase)
}

func onPVUpdate(_, obj interface{}) {
	// Cast the obj as pv
	pv := obj.(*corev1.PersistentVolume)
	fmt.Printf("[onPVUpdate] pv %s with claim %s in %s\n", pv.Name, fmt.Sprintf("%s/%s", pv.Spec.ClaimRef.Namespace, pv.Spec.ClaimRef.Name), pv.Status.Phase)
}

func onPVDelete(obj interface{}) {
	pv := obj.(*corev1.PersistentVolume)
	fmt.Printf("[onPVDelete] pv %s with claim %s in %s deleted\n", pv.Name, fmt.Sprintf("%s/%s", pv.Spec.ClaimRef.Namespace, pv.Spec.ClaimRef.Name), pv.Status.Phase)
}

func onPVCAdd(obj interface{}) {
	// Cast the obj as pv
	pvc := obj.(*corev1.PersistentVolumeClaim)
	fmt.Printf("[onPVCAdd] pvc %s in %s\n", pvc.Name, pvc.Status.Phase)
}

func onPVCUpdate(old, obj interface{}) {
	// Cast the obj as pv
	pvc := obj.(*corev1.PersistentVolumeClaim)
	oldobj := old.(*corev1.PersistentVolumeClaim)
	fmt.Printf("[onPVCUpdate] old pvc %s in %s while new pvc %s in %s\n", oldobj.Name, oldobj.Status.Phase, pvc.Name, pvc.Status.Phase)
}

func onPVCDelete(obj interface{}) {
	// Cast the obj as pv
	pvc := obj.(*corev1.PersistentVolumeClaim)
	fmt.Printf("[onPVCDelete] pvc %s in %s deleted\n", pvc.Name, pvc.Status.Phase)
}
