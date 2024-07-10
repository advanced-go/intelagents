package ingress1

import (
	"fmt"
	"github.com/advanced-go/intelagents/guidance"
	"github.com/advanced-go/intelagents/observation"
	"github.com/advanced-go/observation/inference1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	Class           = "ingress-controller1"
	defaultInterval = time.Second * 3
)

type controller struct {
	running  bool
	uri      string
	origin   core.Origin
	interval time.Duration
	ticker   *time.Ticker
	ctrlC    chan *messaging.Message
	opsAgent messaging.OpsAgent
	shutdown func()
}

func ControllerAgentUri(origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v.%v", Class, origin.Region, origin.Zone, origin.Host)
	}
	return fmt.Sprintf("%v:%v.%v.%v.%v", Class, origin.Region, origin.Zone, origin.SubZone, origin.Host)
}

// NewControllerAgent - create a new controller agent
func NewControllerAgent(origin core.Origin, opsAgent messaging.OpsAgent) messaging.Agent {
	return newControllerAgent(origin, opsAgent)
}

func newControllerAgent(origin core.Origin, opsAgent messaging.OpsAgent) *controller {
	c := new(controller)
	c.origin = origin
	c.uri = ControllerAgentUri(origin)
	c.interval = defaultInterval
	c.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.opsAgent = opsAgent
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
	go run(c, guidance.IngressProcessing, observation.AccessIngressQuery, observation.InferenceIngressQuery, inference1.Insert)
}

func (c *controller) startTicker(interval time.Duration) {
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

func (c *controller) stopTicker() {
	c.ticker.Stop()
}
