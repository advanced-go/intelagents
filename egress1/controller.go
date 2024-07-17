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
	Class = "egress-controller1"
)

type controller struct {
	running  bool
	uri      string
	interval time.Duration // Needs to be configured dynamically during runtime
	ctrlC    chan *messaging.Message
	handler  messaging.Agent
	version  string // Current version of origin configuration, helps to stop duplicate updates of egress routes
	shutdown func()
}

func ControllerAgentUri(origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v.%v", Class, origin.Region, origin.Zone, origin.Host)
	}
	return fmt.Sprintf("%v:%v.%v.%v.%v", Class, origin.Region, origin.Zone, origin.SubZone, origin.Host)
}

// NewControllerAgent - create a new controller agent
func NewControllerAgent(origin core.Origin, handler messaging.Agent) messaging.Agent {
	c := new(controller)
	c.uri = ControllerAgentUri(origin)
	//c.interval = interval

	c.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.handler = handler
	return c
}

// String - identity
func (c *controller) String() string {
	return c.uri
}

// Uri - agent identifier
func (c *controller) Uri() string {
	return c.uri
}

// Message - message the agent
func (c *controller) Message(m *messaging.Message) {
	messaging.Mux(m, c.ctrlC, nil, nil)
}

// Add - add a shutdown function
func (c *controller) Add(f func()) {
	c.shutdown = messaging.AddShutdown(c.shutdown, f)

}

// Shutdown - shutdown the agent
func (c *controller) Shutdown() {
	if !c.running {
		return
	}
	c.running = false
	if c.shutdown != nil {
		c.shutdown()
	}
	msg := messaging.NewControlMessage(c.uri, c.uri, messaging.ShutdownEvent)
	if c.ctrlC != nil {
		c.ctrlC <- msg
	}
}

// Run - run the agent
func (c *controller) Run() {
	if c.running {
		return
	}
	go run(c, access1.EgressQuery, inference1.EgressQuery, nil, inference1.Insert)
}
