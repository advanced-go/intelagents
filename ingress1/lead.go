package ingress1

import (
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	LeadClass = "ingress-lead1"
)

// Responsibilities:
//  1. Startup + Restart Events
//     a. Read all egress route configurations
//     b. If authority routing is configured, read all host names
//     c. Create all egress controller agents and a dependency agent if configured
//  2. Changeset Apply Event
//     a. Read new egress and dependency configurations, update controllers and dependency agent
//  3. Changeset Rollback Event
//     b. Read previous egress and dependency configurations, update controllers and dependency agent
//
// 4. Polling - What if an event is missed?? Need some way to save events in database.
type lead struct {
	running bool
	uri     string

	// Assignment
	origin core.Origin

	// Agents
	controller messaging.Agent
	redirect   messaging.Agent

	interval     time.Duration
	ctrlC        chan *messaging.Message
	handler      messaging.OpsAgent
	shutdownFunc func()
}

func LeadAgentUri(origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v.%v", LeadClass, origin.Region, origin.Zone, origin.Host)
	}
	return fmt.Sprintf("%v:%v.%v.%v.%v", LeadClass, origin.Region, origin.Zone, origin.SubZone, origin.Host)
}

// NewLeadAgent - create a new lead agent
func NewLeadAgent(origin core.Origin, handler messaging.OpsAgent) messaging.OpsAgent {
	c := new(lead)
	c.uri = LeadAgentUri(origin)
	c.origin = origin

	c.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.handler = handler

	c.controller = newControllerAgent(origin, c)
	c.redirect = newRedirectAgent(origin, c)
	return c
}

// String - identity
func (l *lead) String() string {
	return l.uri
}

// Uri - agent identifier
func (l *lead) Uri() string {
	return l.uri
}

// Message - message the agent
func (l *lead) Message(m *messaging.Message) {
	messaging.Mux(m, l.ctrlC, nil, nil)
}

// Handle - error handler
func (l *lead) Handle(status *core.Status, requestId string) *core.Status {
	// TODO : Any operations specific processing ??  If not then forward to handler
	return l.handler.Handle(status, requestId)
}

// AddActivity - add activity
func (l *lead) AddActivity(agentId string, content any) {
	// TODO : Any operations specific processing ??  If not then forward to handler
	//return a.handler.Handle(status, requestId)
}

// Add - add a shutdown function
func (l *lead) Add(f func()) {
	l.shutdownFunc = messaging.AddShutdown(l.shutdownFunc, f)
}

// Shutdown - shutdown the agent
func (l *lead) Shutdown() {
	if !l.running {
		return
	}
	l.running = false
	if l.shutdownFunc != nil {
		l.shutdownFunc()
	}
	l.controller.Shutdown()
	l.redirect.Shutdown()
	msg := messaging.NewControlMessage(l.uri, l.uri, messaging.ShutdownEvent)
	if l.ctrlC != nil {
		l.ctrlC <- msg
	}
}

// Run - run the agent
func (l *lead) Run() {
	if l.running {
		return
	}
	go leadRun(l, newObservation(l.handler), newGuidance(l.handler), newOperations(l.handler))
}

// shutdown - close resources
func (l *lead) shutdown() {
	close(l.ctrlC)

}
