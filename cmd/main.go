package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// crdclientset "github.com/murali-bashyam/crdexample/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func PrintPod(pod *apiv1.Pod) {
	fmt.Println("Pod Name ", pod.Name, " Host IP ", pod.Status.HostIP, " IP ", pod.Status.PodIP)
}

func main() {
	var kubeconfig *string

	kubeconfig = flag.String("kubeconfig", "/home/mbcoder/kubeconfig", "kubeconfig file")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatal("Failed to build config", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("Failed to initialize clientset", err)
	}

/*
	crdclient, err := crdclientset.NewForConfig(config)
	if err != nil {
		log.Fatal("Failed to initialize CRD clientset", err)
	}
*/
	pod, err := clientset.CoreV1().Pods("rook-ceph").Get(context.TODO(), "rook-discover-qkbrg", metav1.GetOptions{})
	if err != nil {
		log.Fatal("Failed to get info about pod", err)
	}

	PrintPod(pod)
	podlist, err := clientset.CoreV1().Pods("rook-ceph").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal("Failed to get podlist", err)
	}

	for _, pod := range podlist.Items {
		PrintPod(&pod)
	}

}
