package ocp4cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	l "github.com/synclib"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1client "k8s.io/client-go/kubernetes/typed/apps/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ListGeter interface {
	GetList(SessionT) []string
}

type SessionT struct {
	coreclient corev1client.CoreV1Client
	appsclient appsv1client.AppsV1Client
}

type IndexError struct{}
type RsNotFound struct{}

var (
	restconfig *rest.Config
)

func (err IndexError) Error() string {
	return "Index not found or is more pods than configs, please check the configmap or access rights"
}

func (err RsNotFound) Error() string {
	return "Replicastion set based on the revision is not found!"
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

func (session *SessionT) GetPods(ctx *context.Context, config *l.Config) (map[string]l.NginxInstance, error) {

	podList := l.New()
	var pod l.NginxInstance

	revision, err := session.getDeploymentRevision(ctx, &config.Deployment, &config.Namespace)
	if err != nil {
		panic(err)
	}

	replicaSet, _ := session.getRsBasedOnRevision(ctx, &revision, &config.Namespace, &config.Deployment)
	if err != nil {
		panic(err)
	}

	rse := strings.Split(replicaSet, "-")
	podsHash := rse[len(rse)-1]

	listOptions := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=ng-plus-apigw,pod-template-hash=%s", podsHash),
	}

	p, err := session.coreclient.Pods(config.Namespace).List(context.Background(), listOptions)

	if err != nil {
		return nil, err
	}

	for _, k := range p.Items {
		pod = l.NginxInstance{
			Address: k.Status.PodIP,
			Port:    config.HttpsPort,
		}
		podList.Pods[k.Name] = pod
	}

	return podList.Pods, nil
}

func (session *SessionT) getDeploymentRevision(ctx *context.Context, deployment *string, namespace *string) (string, error) {

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

func (session *SessionT) getRsBasedOnRevision(ctx *context.Context, revision *string, namespace *string, deployment *string) (string, error) {

	var err error
	var rs string

	listOptions := metav1.ListOptions{
		LabelSelector: "app=ng-plus-apigw",
	}

	rsl, err := session.appsclient.ReplicaSets(*namespace).List(*ctx, listOptions)

	if err != nil {
		return "", err
	}

	for _, k := range rsl.Items {
		for l, m := range k.Annotations {
			if l == "deployment.kubernetes.io/revision" {
				if m == *revision {
					rs = k.Name
					break
				}
			}
		}
		if rs != "" {
			break
		}
	}

	if rs == "" {
		err := RsNotFound{}
		return "", err
	}

	return rs, nil
}
