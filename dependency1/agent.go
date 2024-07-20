package dependency1

import (
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	DependencyClass = "egress-dependency1"
)

type dependency struct {
	running      bool
	uri          string
	origin       core.Origin
	version      string        // Current version of origin configuration, helps to stop duplicate updates of egress routes
	interval     time.Duration // Needs to be configured dynamically during runtime
	ctrlC        chan *messaging.Message
	handler      messaging.OpsAgent
	shutdownFunc func()
}

func DependencyAgentUri(origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v.%v", DependencyClass, origin.Region, origin.Zone, origin.Host)
	}
	return fmt.Sprintf("%v:%v.%v.%v.%v", DependencyClass, origin.Region, origin.Zone, origin.SubZone, origin.Host)
}

// NewDependencyAgent - create a new dependency agent
func NewDependencyAgent(origin core.Origin, handler messaging.OpsAgent) messaging.Agent {
	c := new(dependency)
	c.uri = DependencyAgentUri(origin)
	c.origin = origin

	c.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.handler = handler
	return c
}

// String - identity
func (a *dependency) String() string {
	return a.uri
}

// Uri - agent identifier
func (a *dependency) Uri() string {
	return a.uri
}

// Message - message the agent
func (a *dependency) Message(m *messaging.Message) {
	messaging.Mux(m, a.ctrlC, nil, nil)
}

// Add - add a shutdown function
/*
func (a *dependency) Add(f func()) {
	a.shutdown = messaging.AddShutdown(a.shutdown, f)
}


*/

// Shutdown - shutdown the agent
func (a *dependency) Shutdown() {
	if !a.running {
		return
	}
	a.running = false
	if a.shutdownFunc != nil {
		a.shutdownFunc()
	}
	msg := messaging.NewControlMessage(a.uri, a.uri, messaging.ShutdownEvent)
	if a.ctrlC != nil {
		a.ctrlC <- msg
	}
}

// Run - run the agent
func (a *dependency) Run() {
	if a.running {
		return
	}
	go run(a)
}
