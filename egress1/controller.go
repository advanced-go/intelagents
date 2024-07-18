package egress1

import (
	"fmt"
	"github.com/advanced-go/observation/access1"
	"github.com/advanced-go/observation/inference1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	ControllerClass = "egress-controller1"
)

// What does this agent need to work
// 1. Observation data
//    - Access Log
//    - Experience
//    - Inference
//    - Action
// 2. Guidance data
//    - Global processing schedule
//    - Dependency processing schedule
//    - Controller Configuration
//

// TODO : add support for control messages or restart, apply-changes, rollback-changes

type controller struct {
	running  bool
	uri      string
	origin   core.Origin
	version  string        // Current version of origin configuration, helps to stop duplicate updates of egress routes
	interval time.Duration // Needs to be configured dynamically during runtime
	ctrlC    chan *messaging.Message
	handler  messaging.OpsAgent
	shutdown func()
}

func ControllerAgentUri(origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v.%v", ControllerClass, origin.Region, origin.Zone, origin.Host)
	}
	return fmt.Sprintf("%v:%v.%v.%v.%v", ControllerClass, origin.Region, origin.Zone, origin.SubZone, origin.Host)
}

// NewControllerAgent - create a new controller agent
func NewControllerAgent(origin core.Origin, handler messaging.OpsAgent) messaging.OpsAgent {
	c := new(controller)
	c.uri = ControllerAgentUri(origin)
	c.origin = origin

	c.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.handler = handler
	return c
}

// String - identity
func (a *controller) String() string {
	return a.uri
}

// Uri - agent identifier
func (a *controller) Uri() string {
	return a.uri
}

// Message - message the agent
func (a *controller) Message(m *messaging.Message) {
	messaging.Mux(m, a.ctrlC, nil, nil)
}

// Handle - error handler
func (a *controller) Handle(status *core.Status, requestId string) *core.Status {
	// TODO : Any controller specific processint ??  If not then forward to handler
	return a.handler.Handle(status, requestId)
}

// Add - add a shutdown function
func (a *controller) Add(f func()) {
	a.shutdown = messaging.AddShutdown(a.shutdown, f)

}

// Shutdown - shutdown the agent
func (a *controller) Shutdown() {
	if !a.running {
		return
	}
	a.running = false
	if a.shutdown != nil {
		a.shutdown()
	}
	msg := messaging.NewControlMessage(a.uri, a.uri, messaging.ShutdownEvent)
	if a.ctrlC != nil {
		a.ctrlC <- msg
	}
}

// Run - run the agent
func (a *controller) Run() {
	if a.running {
		return
	}
	go runController(a, access1.EgressQuery, inference1.EgressQuery, nil, nil)
}
