package guidance1

import (
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	RouteMonitorClass = "route-monitor1"
)

type route struct {
	running  bool
	uri      string
	origin   core.Origin
	ticker   *messaging.Ticker
	ctrlC    chan *messaging.Message
	handler  messaging.OpsAgent
	shutdown func()
}

func RouteMonitorAgentUri(origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v.%v", RouteMonitorClass, origin.Region, origin.Zone, origin.Host)
	}
	return fmt.Sprintf("%v:%v.%v.%v.%v", RouteMonitorClass, origin.Region, origin.Zone, origin.SubZone, origin.Host)
}

// NewRouteMonitorAgent - create a new routing monitor agent
func NewRouteMonitorAgent(interval time.Duration, origin core.Origin, handler messaging.OpsAgent) messaging.Agent {
	return newRouteMonitorAgent(interval, origin, handler)
}

func newRouteMonitorAgent(interval time.Duration, origin core.Origin, handler messaging.OpsAgent) *route {
	c := new(route)
	c.uri = RouteMonitorAgentUri(origin)
	c.origin = origin
	c.ticker = messaging.NewTicker(interval)
	c.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.handler = handler
	return c
}

// String - identity
func (p *route) String() string {
	return p.uri
}

// Uri - agent identifier
func (p *route) Uri() string {
	return p.uri
}

// Message - message the agent
func (p *route) Message(m *messaging.Message) {
	messaging.Mux(m, p.ctrlC, nil, nil)
}

// Shutdown - shutdown the agent
func (p *route) Shutdown() {
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
func (p *route) Run() {
	if p.running {
		return
	}
	go runRouting(p)
}
