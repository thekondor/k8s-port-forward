package main

import (
	"flag"
	"github.com/thekondor/k8s-port-forward/app_skeleton"
	"log"
)

type output struct {
	remotePort int
}

func (self *output) OnReady(localPort int) {
	log.Println("PortForwarding is ready :)")
	log.Println("remote:", self.remotePort)
	log.Println("local:", localPort)
}

func (*output) OnExit() {
	log.Println("Exiting...")
}

func main() {
	kubeConfig := flag.String("kubeconfig", "", `Path to 'kubeconfig.yaml'`)
	podPort := flag.Int("pod-port", 27017, `Remote Pod's port`)
	podName := flag.String("pod-name", "", `Remote Pod's name`)
	ns := flag.String("ns", "", "Namespace")

	flag.Parse()

	if "" == *kubeConfig {
		log.Fatal(`No path to 'kubeconfig.yaml' specified`)
	}
	if "" == *podName {
		log.Fatal(`No pod name specified`)
	}

	app, err := app_skeleton.NewBasicAndNaive(app_skeleton.Spec{
		KubeConfigPath: *kubeConfig,
		PodName:        *podName,
		Namespace:      *ns,
		RemotePort:     *podPort,
	})
	if nil != err {
		log.Fatal(err)
	}

	app.Run(&output{remotePort: *podPort})
}
