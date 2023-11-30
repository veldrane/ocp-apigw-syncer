package ocp4cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1client "k8s.io/client-go/kubernetes/typed/apps/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ListGeter interface {
	GetList(SessionT) []string
}

type PodsT struct{}

type SessionT struct {
	coreclient corev1client.CoreV1Client
	appsclient appsv1client.AppsV1Client
}

type IndexError struct{}

var (
	restconfig *rest.Config
)

func (err IndexError) Error() string {
	return fmt.Sprintf("Index not found or is more pods than configs, please check the configmap or access rights")
}

func Session() *SessionT {

	restconfig := getRestConfig()

	coreclient, err := corev1client.NewForConfig(restconfig)
	if err != nil {
		panic("Cannot login to the cluster")
	}

	appsclient, err := appsv1client.NewForConfig(restconfig)
	if err != nil {
		panic("Cannot login to the cluster")
	}
	res := SessionT{
		coreclient: *coreclient,
		appsclient: *appsclient,
	}

	return &res
}

func inPod() bool {

	if k := os.Getenv("KUBERNETES_PORT"); k == "" {
		return false
	}

	return true
}

func getRestConfig() *rest.Config {

	var err error

	if inPod() {
		restconfig, err = rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
	} else {
		kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			clientcmd.NewDefaultClientConfigLoadingRules(),
			&clientcmd.ConfigOverrides{},
		)
		restconfig, err = kubeconfig.ClientConfig()
		if err != nil {
			panic(err)
		}
	}

	return restconfig
}

func GetNamespace() string {

	namespace := os.Getenv("POD_NAMESPACE")
	return namespace
}

func (pods PodsT) GetList(session *SessionT, namespace *string, replicaSet *string) ([]string, error) {

	var podList []string

	listOptions := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=ng-plus-apigw,pod-template-hash=%s", *replicaSet),
	}

	p, err := session.coreclient.Pods(*namespace).List(context.Background(), listOptions)

	if err != nil {
		return nil, err
	}

	//sort.Slice(p.Items, func(i, j int) bool {
	//	return p.Items[i].CreationTimestamp.Before(&p.Items[j].CreationTimestamp)
	//} )

	for _, k := range p.Items {
		podList = append(podList, k.Name)
	}

	sort.Strings(podList)

	return podList, nil
}

func GetReplicationSet() string {

	hostname := os.Getenv("HOSTNAME")
	parts := strings.Split(hostname, "-")
	hash := ""
	for i := 0; i < (len(parts) - 1); i++ {
		hash = (parts[i])
	}

	//rs := ("ng-plus-apigw-" + hash)

	return hash
}

func GetIndex(podlist []string, hostname *string, path *string) (int, error) {

	files, _ := filepath.Glob(*path + "config-*.yaml")

	error := IndexError{}

	if len(files) < len(podlist) {
		return 0, error
	}

	idx := 0

	for i, k := range podlist {
		if k == *hostname {
			idx = i % len(podlist)
			break
		}
	}
	return idx, nil
}

func (session *SessionT) GetDeploymentRevision(ctx *context.Context, deployment *string, namespace *string) (string, error) {

	var revision string

	getOptions := metav1.GetOptions{}

	d, err := session.appsclient.Deployments(*namespace).Get(*ctx, *deployment, getOptions)

	if err != nil {
		return "", err
	}

	for k, v := range d.Annotations {

		if k == "deployment.kubernetes.io/revision" {
			revision = v
			break
		}
	}

	return revision, nil
}
