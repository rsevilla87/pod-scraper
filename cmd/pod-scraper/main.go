package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/rsevilla87/pod-scraper/pkg/discovery"
	"github.com/rsevilla87/pod-scraper/pkg/scraper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Config struct {
	nsLabel    *string
	podLabel   *string
	urlScheme  *string
	endpoint   *string
	targetPort *int
	code       *int
}

func parseFlags() Config {
	config := Config{
		nsLabel:    flag.String("ns-label", "", "Target namespace label"),
		podLabel:   flag.String("pod-label", "", "Target pod label"),
		urlScheme:  flag.String("scheme", "", "URL scheme, http or https"),
		endpoint:   flag.String("endpoint", "/", "Target endpoint"),
		targetPort: flag.Int("port", 0, "Target port"),
		code:       flag.Int("code", 200, "Expected status code"),
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
	var wg sync.WaitGroup
	var failed int
	config := parseFlags()
	clientSet := getClientSet()
	nsList, err := discovery.DiscoverNamespaces(clientSet, *config.nsLabel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(nsList.Items) < 1 {
		fmt.Printf("No namespaces discovered with labels %v", *config.nsLabel)
		os.Exit(0)
	}
	podList, err := discovery.DiscoverPods(clientSet, nsList, *config.podLabel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(podList) < 1 {
		fmt.Printf("No pods discovered with labels %v", *config.podLabel)
		os.Exit(0)
	}
	for _, pod := range podList {
		target := fmt.Sprintf("%v://%v%v", config.urlScheme, pod.Status.PodIP, config.endpoint)
		go scraper.Scrape(target, *config.code, &wg, &failed)
	}
	wg.Wait()
	os.Exit(failed)
}
