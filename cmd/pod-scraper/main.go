package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rsevilla87/pod-scraper/pkg/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Config struct {
	nsLabel    *string
	podLabel   *string
	urlScheme  *string
	targetPort *int
}

func parseFlags() Config {
	config := Config{
		nsLabel:    flag.String("ns-label", "", "Target namespace label"),
		podLabel:   flag.String("pod-label", "", "Target pod label"),
		urlScheme:  flag.String("scheme", "", "URL scheme, http or https"),
		targetPort: flag.Int("port", 0, "Target port"),
	}
	flag.Parse()
	return config
}

func getClientSet() *kubernetes.Clientset {
	restConfig, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Loaded k8s configuration")
	return kubernetes.NewForConfigOrDie(restConfig)
}

func main() {
	config := parseFlags()
	clientSet := getClientSet()
	nsList, err := discovery.DiscoverNamespaces(clientSet, *config.nsLabel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(nsList.Items) > 0 {
		podList, err := discovery.DiscoverPods(clientSet, nsList, *config.podLabel)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, pod := range podList {
			fmt.Println(pod.Status.PodIP)
		}
	} else {
		fmt.Printf("No namespaces discovered with labels %v", *config.nsLabel)
	}
}
