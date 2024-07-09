package ingress1

import (
	"fmt"
	"github.com/advanced-go/intelagents/guidance"
	"github.com/advanced-go/observation/access1"
	"github.com/advanced-go/observation/inference1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	Class = "ingress-controller1"
)

type controller struct {
	running  bool
	uri      string
	interval time.Duration
	ticker   *time.Ticker
	ctrlC    chan *messaging.Message
	handler  messaging.Agent
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
	return newControllerAgent(origin, handler)
}

func newControllerAgent(origin core.Origin, handler messaging.Agent) *controller {
	c := new(controller)
	c.uri = ControllerAgentUri(origin)
	c.interval = guidance.IngressInterval()
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
	go run(c, access1.IngressQuery, inference1.IngressQuery, nil, inference1.Insert)
}

func (c *controller) StartTicker(interval time.Duration) {
	if interval <= 0 {
		interval = c.interval
	} else {
		c.interval = interval
	}
	if c.ticker != nil {
		c.ticker.Stop()
	}
	c.ticker = time.NewTicker(interval)
}

func (c *controller) StopTicker() {
	c.ticker.Stop()
}
