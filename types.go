package port_forward

import (
	"fmt"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type PortForwarding struct {
	restConfig *rest.Config
	pod        v1.Pod
	localPort  int
	podPort    int
	stopCh     chan struct{}
	readyCh    chan struct{}
}

func NewPortForwarding(spec PortForwardingSpec) (PortForwarding, error) {
	restConfig, err := clientcmd.BuildConfigFromFlags("", spec.KubeConfigPath)
	if nil != err {
		return PortForwarding{}, err
	}

	return PortForwarding{
		restConfig: restConfig,
		pod: v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      spec.Pod.Name,
				Namespace: spec.Pod.Namespace,
			},
		},
		localPort: spec.Ports.Local,
		podPort:   spec.Ports.Remote,
		stopCh:    make(chan struct{}, 1),
		readyCh:   make(chan struct{}),
	}, nil
}

func (self PortForwarding) Start() error {
	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", self.pod.Namespace, self.pod.Name)
	hostIP := strings.TrimLeft(self.restConfig.Host, "htps:/")

	transport, upgrader, err := spdy.RoundTripperFor(self.restConfig)
	if err != nil {
		return err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost,
		&url.URL{
			Scheme: "https",
			Path:   path,
			Host:   hostIP,
		},
	)

	ports := []string{fmt.Sprintf("%d:%d", self.localPort, self.podPort)}
	fw, err := portforward.New(dialer, ports, self.stopCh, self.readyCh, ioutil.Discard, os.Stderr)
	if err != nil {
		return err
	}
	return fw.ForwardPorts()
}

func (self PortForwarding) WaitForReady() {
	<-self.readyCh
}

func (self PortForwarding) Stop() {
	close(self.stopCh)
}
