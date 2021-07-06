package discovery

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func DiscoverNamespaces(clientSet *kubernetes.Clientset, nsLabels string) (*v1.NamespaceList, error) {
	return clientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
		LabelSelector: nsLabels,
	})
}

func DiscoverPods(clientSet *kubernetes.Clientset, nsList *v1.NamespaceList, podLabels string) ([]v1.Pod, error) {
	var podList []v1.Pod
	listOptions := metav1.ListOptions{
		LabelSelector: podLabels,
	}
	if nsList != nil {
		for _, ns := range nsList.Items {
			tempPodList, err := clientSet.CoreV1().Pods(ns.Name).List(context.TODO(), listOptions)
			if err != nil {
				return []v1.Pod{}, err
			}
			podList = append(podList, tempPodList.Items...)
		}
	} else {
		tempPodList, err := clientSet.CoreV1().Pods(v1.NamespaceAll).List(context.TODO(), listOptions)
		if err != nil {
			return []v1.Pod{}, err
		}
		podList = tempPodList.Items
	}
	return podList, nil
}
