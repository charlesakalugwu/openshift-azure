//+build e2e

package e2e

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/ghodss/yaml"
	project "github.com/openshift/api/project/v1"
	projectclient "github.com/openshift/client-go/project/clientset/versioned/typed/project/v1"
	routev1client "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	templatev1client "github.com/openshift/client-go/template/clientset/versioned/typed/template/v1"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var c *testClient

type testClient struct {
	kc        *kubernetes.Clientset
	pc        *projectclient.ProjectV1Client
	rc        *routev1client.RouteV1Client
	tc        *templatev1client.TemplateV1Client
	namespace string

	artifactDir string
}

func newTestClient(kubeconfig, artifactDir string) *testClient {
	var err error
	var config *rest.Config

	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err)
		}
	} else {
		// use in-cluster config if no kubeconfig has been specified
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}

	// create the clientset
	kc, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// create a project client for creating and tearing down namespaces
	pc, err := projectclient.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// create a template client
	tc, err := templatev1client.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// create a route client
	rc, err := routev1client.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return &testClient{
		kc: kc,
		pc: pc,
		rc: rc,
		tc: tc,

		artifactDir: artifactDir,
	}
}

func (t *testClient) createProject(namespace string) error {
	if _, err := t.pc.ProjectRequests().Create(&project.ProjectRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		}}); err != nil {
		return err
	}
	t.namespace = namespace

	if err := wait.PollImmediate(2*time.Second, time.Minute, t.selfSarSuccess); err != nil {
		return fmt.Errorf("failed to wait for self-sar success: %v", err)
	}
	if err := wait.PollImmediate(2*time.Second, time.Minute, t.defaultServiceAccountIsReady); err != nil {
		return fmt.Errorf("failed to wait for the default service account provision: %v", err)
	}
	return nil
}

func (t *testClient) cleanupProject(timeout time.Duration) error {
	if t.namespace == "" {
		return nil
	}
	if err := t.pc.Projects().Delete(t.namespace, &metav1.DeleteOptions{}); err != nil {
		return err
	}
	if err := wait.PollImmediate(2*time.Second, timeout, t.projectIsCleanedUp); err != nil {
		return fmt.Errorf("failed to wait for project cleanup: %v", err)
	}
	return nil
}

func (t *testClient) dumpInfo() error {
	// gather events
	eventList, err := t.kc.CoreV1().Events(t.namespace).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	eventBuf := bytes.NewBuffer(nil)
	for _, event := range eventList.Items {
		b, err := yaml.Marshal(event)
		if err != nil {
			return err
		}
		if _, err := eventBuf.Write(b); err != nil {
			return err
		}
		if _, err := eventBuf.Write([]byte("\n")); err != nil {
			return err
		}
	}

	// gather pods
	podList, err := t.kc.CoreV1().Pods(t.namespace).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	podBuf := bytes.NewBuffer(nil)
	for _, pod := range podList.Items {
		b, err := yaml.Marshal(pod)
		if err != nil {
			return err
		}
		if _, err := podBuf.Write(b); err != nil {
			return err
		}
		if _, err := podBuf.Write([]byte("\n")); err != nil {
			return err
		}
	}

	if t.artifactDir != "" {
		if err := ioutil.WriteFile(filepath.Join(t.artifactDir, fmt.Sprintf("events-%s.yaml", t.namespace)), eventBuf.Bytes(), 0777); err != nil {
			logrus.Warn(err)
			fmt.Println(eventBuf.String())
		}
		if err := ioutil.WriteFile(filepath.Join(t.artifactDir, fmt.Sprintf("pods-%s.yaml", t.namespace)), podBuf.Bytes(), 0777); err != nil {
			logrus.Warn(err)
			fmt.Println(podBuf.String())
		}
	} else {
		fmt.Println(eventBuf.String())
		fmt.Println(podBuf.String())
	}
	return nil
}
