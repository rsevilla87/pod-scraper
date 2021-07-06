package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/rsevilla87/pod-scraper/pkg/discovery"
	"github.com/rsevilla87/pod-scraper/pkg/scraper"
	v1 "k8s.io/api/core/v1"
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
	timeout    *time.Duration
}

func parseFlags() Config {
	config := Config{
		nsLabel:    flag.String("ns-label", "", "Target namespace label"),
		podLabel:   flag.String("pod-label", "", "Target pod label"),
		urlScheme:  flag.String("scheme", "http", "URL scheme, http or https"),
		endpoint:   flag.String("endpoint", "/", "Target endpoint"),
		targetPort: flag.Int("port", 80, "Target port"),
		code:       flag.Int("code", 200, "Expected status code"),
		timeout:    flag.Duration("timeout", 10*time.Second, "Request timeout"),
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
	var err error
	var podList []v1.Pod
	var nsList *v1.NamespaceList
	var wg sync.WaitGroup
	config := parseFlags()
	if *config.nsLabel == "" && *config.podLabel == "" {
		fmt.Println("Either -ns-label or -pod-label flags must be set")
		os.Exit(1)
	}
	scraper := scraper.NewScraper(&wg, *config.code, *config.timeout)
	clientSet := getClientSet()
	// Discover pods within the namespaces
	if *config.nsLabel != "" {
		nsList, err = discovery.DiscoverNamespaces(clientSet, *config.nsLabel)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if len(nsList.Items) < 1 {
			fmt.Printf("No namespaces discovered with labels %v\n", *config.nsLabel)
			os.Exit(0)
		}
		fmt.Printf("Discovered %d namespaces with labels %s\n", len(nsList.Items), *config.nsLabel)
	}
	podList, err = discovery.DiscoverPods(clientSet, nsList, *config.podLabel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(podList) < 1 {
		fmt.Printf("0 pods discovered%v\n", *config.podLabel)
		os.Exit(0)
	}
	fmt.Printf("Discovered %d pods\n", len(podList))
	for _, pod := range podList {
		target := fmt.Sprintf("%v://%v:%d%v", *config.urlScheme, pod.Status.PodIP, *config.targetPort, *config.endpoint)
		wg.Add(1)
		go scraper.Scrape(target)
	}
	wg.Wait()
	os.Exit(scraper.Failed)
}
