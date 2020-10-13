package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig *string
	// Get config file with absolute path
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	configNamespace := os.Getenv("cmNamespace")
	configName := os.Getenv("cmName")
	dataName := os.Getenv("dataName")
	flag.Parse()
	fmt.Println(configNamespace)
	fmt.Println(configName)
	fmt.Println(dataName)
	// Uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	// Get gang minMember configMap
	cm, err := clientset.CoreV1().ConfigMaps(configNamespace).Get(configName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	// Update configMap
	num, _ := strconv.Atoi(cm.Data[dataName])
	num -= 1
	cm.Data[dataName] = strconv.Itoa(num)
	_, _ = clientset.CoreV1().ConfigMaps(configNamespace).Update(cm)
	// Check if we satisfy gang minMember or not every 5 second
	for {
		cm, err := clientset.CoreV1().ConfigMaps(configNamespace).Get(configName, metav1.GetOptions{})
		if err != nil {
			panic(err.Error())
		}
		number, _ := strconv.Atoi(cm.Data[dataName])
		if number <= 0 {
			fmt.Println("satisfy gang minMember.")
			fmt.Println("start to run job.")
			time.Sleep(60 * time.Second) // means the application start running job
			break
		}
		time.Sleep(2 * time.Second)
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return ""
}
