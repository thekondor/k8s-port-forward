package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/thekondor/k8s-port-forward/app_skeleton"
	"log"
)

type output struct {
	remotePort int
	jsonOutput bool
}

func (self *output) OnReady(localPort int) {
	log.Println("PortForwarding is ready :)")
	pfInfo := self.getPfInfo(localPort)
	log.Println("[ports]", pfInfo)
}

func (self output) getPfInfo(localPort int) string {
	if self.jsonOutput {
		info, err := json.Marshal(struct {
			Remote int `json:"remote_port"`
			Local  int `json:"local_port"`
		}{
			Remote: self.remotePort,
			Local:  localPort,
		})
		if nil != err {
			log.Panicf("Failed to build JSON output: %v", err)
		}

		return string(info)
	}

	return fmt.Sprintf("remote_port=%d, local_port=%d", self.remotePort, localPort)
}

func (*output) OnExit() {
	log.Println("Exiting...")
}

func main() {
	kubeConfig := flag.String("kubeconfig", "", `Path to 'kubeconfig.yaml'`)
	podPort := flag.Int("pod-port", 27017, `Remote Pod's port`)
	podName := flag.String("pod-name", "", `Remote Pod's name`)
	ns := flag.String("ns", "", "Namespace")
	jsonOutput := flag.Bool("json", false, "Output active port forwarding info as JSON")

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

	app.Run(&output{remotePort: *podPort, jsonOutput: *jsonOutput})
}
