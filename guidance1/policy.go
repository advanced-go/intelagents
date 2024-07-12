package guidance1

import (
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	PolicyMonitorClass = "policy-monitor1"
)

type policy struct {
	running  bool
	uri      string
	origin   core.Origin
	ticker   *messaging.Ticker
	ctrlC    chan *messaging.Message
	handler  messaging.OpsAgent
	shutdown func()
}

func PolicyMonitorAgentUri(origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v.%v", PolicyMonitorClass, origin.Region, origin.Zone, origin.Host)
	}
	return fmt.Sprintf("%v:%v.%v.%v.%v", PolicyMonitorClass, origin.Region, origin.Zone, origin.SubZone, origin.Host)
}

// NewPolicyMonitorAgent - create a new schedule agent
func NewPolicyMonitorAgent(interval time.Duration, origin core.Origin, handler messaging.OpsAgent) messaging.Agent {
	return newPolicyMonitorAgent(interval, origin, handler)
}

func newPolicyMonitorAgent(interval time.Duration, origin core.Origin, handler messaging.OpsAgent) *policy {
	c := new(policy)
	c.uri = PolicyMonitorAgentUri(origin)
	c.origin = origin
	c.ticker = messaging.NewTicker(interval)
	c.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.handler = handler
	return c
}

// String - identity
func (p *policy) String() string {
	return p.uri
}

// Uri - agent identifier
func (p *policy) Uri() string {
	return p.uri
}

// Message - message the agent
func (p *policy) Message(m *messaging.Message) {
	messaging.Mux(m, p.ctrlC, nil, nil)
}

// Shutdown - shutdown the agent
func (p *policy) Shutdown() {
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
func (p *policy) Run() {
	if p.running {
		return
	}
	go runPolicy(p)
}
