package port_forward

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
)

type PodSpec struct {
	Name      string
	Namespace string
}

type PortsSpec struct {
	Local  int
	Remote int
}

type PortForwardingSpec struct {
	KubeConfigPath string
	Pod            PodSpec
	Ports          PortsSpec
}

type PortForwarding struct {
	restConfig *rest.Config
	pod        v1.Pod
	localPort  int
	podPort    int
	stopCh     chan struct{}
	readyCh    chan struct{}
}
