package port_forward

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
