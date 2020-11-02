package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	crdv1 "github.com/murali-bashyam/crdexample/pkg/apis/crd.example.com/v1"
	crdclientset "github.com/murali-bashyam/crdexample/pkg/client/clientset/versioned"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
)

var incluster bool = false

func PrintPod(pod *apiv1.Pod) {
	fmt.Println("Pod Name ", pod.Name, " Host IP ", pod.Status.HostIP, " IP ", pod.Status.PodIP)
}

func createOrUpdatePool(c *crdclientset.Clientset, poolname string, quota int, domain string) error {

	pool, err := c.CrdV1().StoragePools("default").Get(context.TODO(), poolname, metav1.GetOptions{})

	if err != nil {
		pool := &crdv1.StoragePool{
			ObjectMeta: metav1.ObjectMeta{
				Name: poolname,
			},
			PoolSpec: crdv1.StoragePoolSpec{
				FailureDomain: domain,
				Quota:         quota,
			},
		}
		_, err = c.CrdV1().StoragePools("default").Create(context.TODO(), pool, metav1.CreateOptions{})
		if err == nil {
			fmt.Printf("Storage pool created %s \n", poolname)
		}
	} else {
		err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
			pool, err = c.CrdV1().StoragePools("default").Get(context.TODO(), poolname, metav1.GetOptions{})
			if err == nil {
				pool.PoolSpec.Quota = quota
				pool.PoolSpec.FailureDomain = domain
				_, err = c.CrdV1().StoragePools("default").Update(context.TODO(), pool, metav1.UpdateOptions{})
				if err == nil {
					fmt.Printf("Storage pool updated %s \n", poolname)
				}
			}
			return err
		})
	}

	return err
}

func deletePool(c *crdclientset.Clientset, poolname string) error {
	err := c.CrdV1().StoragePools("default").Delete(context.TODO(), poolname, metav1.DeleteOptions{})
	if err == nil {
		fmt.Printf("Storage pool deleted %s \n", poolname)
	}
	return err
}

func main() {
	var kubeconfig *string
	var config *restclient.Config

	if len(os.Args) > 1 {
		if os.Args[1] == "incluster" {
			incluster = true
		}
	}

	if incluster == false {
		kubeconfig = flag.String("kubeconfig", "/home/mbcoder/kubeconfig", "kubeconfig file")
		flag.Parse()

		kconfig, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			log.Fatal("Failed to build config", err)
		}
		config = kconfig
	} else {
		kconfig, err := restclient.InClusterConfig()
		if err != nil {
			log.Fatal("Failed to build config", err)
		} else {
			fmt.Println("Successfully built incluster config")
		}
		config = kconfig
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("Failed to initialize clientset", err)
	}

	crdclient, err := crdclientset.NewForConfig(config)
	if err != nil {
		log.Fatal("Failed to initialize CRD clientset", err)
	}

	for {
		fmt.Println("")
		fmt.Println("-----------------------------------------------------------")
		err = createOrUpdatePool(crdclient, "blockpool", 1024, "rack")
		if err != nil {
			log.Fatal("Failed to create storage pool ", err)
		}

		err = createOrUpdatePool(crdclient, "bpool", 512, "host")
		if err != nil {
			log.Fatal("Failed to create storage pool ", err)
		}

		pool_list, err := crdclient.CrdV1().StoragePools("default").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatal("Error listing all storage pools", err)
		}

		for _, pool := range pool_list.Items {
			fmt.Printf("Pool %s with quota %d and failureDomain %s \n", pool.Name,
				pool.PoolSpec.Quota, pool.PoolSpec.FailureDomain)
		}

		err = deletePool(crdclient, "bpool")
		if err != nil {
			log.Fatal("Failed to delete storage pool ", err)
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

		time.Sleep(10 * time.Second)
	}
}
