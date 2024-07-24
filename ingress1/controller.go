package ingress1

import (
	"fmt"
	"github.com/advanced-go/guidance/percentile1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	Class           = "ingress-controller1"
	defaultInterval = time.Minute * 2
)

type controller struct {
	running      bool
	uri          string
	origin       core.Origin
	ticker       *messaging.Ticker
	poller       *messaging.Ticker
	ctrlC        chan *messaging.Message
	handler      messaging.OpsAgent
	shutdownFunc func()
}

func ControllerAgentUri(origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v.%v", Class, origin.Region, origin.Zone, origin.Host)
	}
	return fmt.Sprintf("%v:%v.%v.%v.%v", Class, origin.Region, origin.Zone, origin.SubZone, origin.Host)
}

// NewControllerAgent - create a new controller agent
func NewControllerAgent(origin core.Origin, handler messaging.OpsAgent) messaging.Agent {
	return newControllerAgent(origin, handler)
}

func newControllerAgent(origin core.Origin, handler messaging.OpsAgent) *controller {
	c := new(controller)
	c.origin = origin
	c.uri = ControllerAgentUri(origin)
	c.ticker = messaging.NewTicker(defaultInterval)
	c.poller = messaging.NewTicker(percentile1.PercentilePollingDuration)
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
	if c.shutdownFunc != nil {
		c.shutdownFunc()
	}
	msg := messaging.NewControlMessage(c.uri, c.uri, messaging.ShutdownEvent)
	if c.ctrlC != nil {
		c.ctrlC <- msg
	}
}

// shutdown - close resources
func (c *controller) shutdown() {
	close(c.ctrlC)
	c.stopTickers()
}

// Run - run the agent
func (c *controller) Run() {
	if c.running {
		return
	}
	go runControl(c, newObservation(c.handler), newGuidance(c.handler), newInference(c.handler), newOperations(c.handler))
}

func (c *controller) stopTickers() {
	c.ticker.Stop()
	c.poller.Stop()
}

func (c *controller) updateTicker(observe *observation) {

}

/*
func (c *controller) startTicker(duration time.Duration) {
	c.ticker.Start(duration)
	if c.ticker != nil {
		c.ticker.Stop()
	}
	c.ticker = time.NewTicker(interval)
}

func (c *controller) stopTicker() {
	c.ticker.Stop()
}

func (c *controller) startPoller(interval time.Duration) {
	if interval <= 0 {
		interval = c.tickInterval
	} else {
		c.tickInterval = interval
	}
	if c.ticker != nil {
		c.ticker.Stop()
	}
	c.ticker = time.NewTicker(interval)
}

func (c *controller) stopPoller() {
	c.ticker.Stop()
}


*/
