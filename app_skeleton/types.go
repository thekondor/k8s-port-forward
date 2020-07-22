package app_skeleton

type Spec struct {
	KubeConfigPath string
	PodName        string
	Namespace      string
	RemotePort     int
}

type Handler interface {
	OnReady(local_port int)
	OnExit()
}
