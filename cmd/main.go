package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	crdclientset "github.com/murali-bashyam/crdexample/pkg/client/clientset/versioned"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	crdclient, err := crdclientset.NewForConfig(config)
	if err != nil {
		log.Fatal("Failed to initialize CRD clientset", err)
	}

	pool_list, err := crdclient.CrdV1().StoragePools("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal("Error listing all storage pool: %v", err)
	}

	for _, pool := range pool_list.Items {
		fmt.Printf("Pool %s with quota %d and failureDomain %s \n", pool.Name,
			pool.PoolSpec.Quota, pool.PoolSpec.FailureDomain)
	}

	pod, err := clientset.CoreV1().Pods("rook-ceph").Get(context.TODO(), "rook-discover-w5vz5", metav1.GetOptions{})
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
