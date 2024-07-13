package guidance1

import (
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	RoutingMonitorClass = "routing-monitor1"
)

type routing struct {
	running  bool
	uri      string
	origin   core.Origin
	ticker   *messaging.Ticker
	ctrlC    chan *messaging.Message
	handler  messaging.OpsAgent
	shutdown func()
}

func RoutingMonitorAgentUri(origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v.%v", RoutingMonitorClass, origin.Region, origin.Zone, origin.Host)
	}
	return fmt.Sprintf("%v:%v.%v.%v.%v", RoutingMonitorClass, origin.Region, origin.Zone, origin.SubZone, origin.Host)
}

// NewRoutingMonitorAgent - create a new routing monitor agent
func NewRoutingMonitorAgent(interval time.Duration, origin core.Origin, handler messaging.OpsAgent) messaging.Agent {
	return newRoutingMonitorAgent(interval, origin, handler)
}

func newRoutingMonitorAgent(interval time.Duration, origin core.Origin, handler messaging.OpsAgent) *routing {
	c := new(routing)
	c.uri = RoutingMonitorAgentUri(origin)
	c.origin = origin
	c.ticker = messaging.NewTicker(interval)
	c.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.handler = handler
	return c
}

// String - identity
func (p *routing) String() string {
	return p.uri
}

// Uri - agent identifier
func (p *routing) Uri() string {
	return p.uri
}

// Message - message the agent
func (p *routing) Message(m *messaging.Message) {
	messaging.Mux(m, p.ctrlC, nil, nil)
}

// Shutdown - shutdown the agent
func (p *routing) Shutdown() {
	if !p.running {
		return
	}
	p.running = false
	if p.shutdown != nil {
		p.shutdown()
	}
	msg := messaging.NewControlMessage(p.uri, p.uri, messaging.ShutdownEvent)
	if p.ctrlC != nil {
		p.ctrlC <- msg
	}
}

// Run - run the agent
func (p *routing) Run() {
	if p.running {
		return
	}
	go runRouting(p)
}
