package app_skeleton

import (
	"github.com/thekondor/anyport"
	. "github.com/thekondor/k8s-port-forward"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type app struct {
	portForwarding PortForwarding
	signals        chan os.Signal
	localPort      int
	wg             sync.WaitGroup
}

func NewBasicAndNaive(spec Spec) (app, error) {
	anyPort, err := anyport.ListenInsecure("localhost")
	if nil != err {
		return app{}, err
	}

	localPort := anyPort.PortNumber
	// TODO: not reliable, a point of a possible race
	anyPort.Listener.Close()

	portForwarding, err := NewPortForwarding(PortForwardingSpec{
		KubeConfigPath: spec.KubeConfigPath,
		Pod: PodSpec{
			Name: spec.PodName, Namespace: spec.Namespace,
		},
		Ports: PortsSpec{
			Local: localPort, Remote: spec.RemotePort,
		},
	})
	if nil != err {
		return app{}, err
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	return app{portForwarding: portForwarding, signals: signals, localPort: localPort}, nil
}

func (self *app) Run(handler Handler) {
	self.wg.Add(1)
	go func() {
		<-self.signals
		self.portForwarding.Stop()
		handler.OnExit()
		self.wg.Done()
	}()

	go func() {
		err := self.portForwarding.Start()
		if nil != err {
			log.Fatal("Failed to forward port", err)
		}
	}()

	self.portForwarding.WaitForReady()

	handler.OnReady(self.localPort)

	self.wg.Wait()
}
