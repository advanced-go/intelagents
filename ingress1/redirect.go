package ingress1

import (
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	RedirectClass = "ingress-redirect1"
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

type redirectState struct {
	Location string
	Percent  string
}

type redirect struct {
	running bool
	agentId string

	// Assignment
	origin core.Origin
	state  *redirectState

	interval     time.Duration
	ctrlC        chan *messaging.Message
	handler      messaging.OpsAgent
	shutdownFunc func()
}

func redirectAgentUri(origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v.%v", RedirectClass, origin.Region, origin.Zone, origin.Host)
	}
	return fmt.Sprintf("%v:%v.%v.%v.%v", RedirectClass, origin.Region, origin.Zone, origin.SubZone, origin.Host)
}

// newRedirectAgent - create a new lead agent
func newRedirectAgent(origin core.Origin, handler messaging.OpsAgent) messaging.OpsAgent {
	return newRedirect(origin, handler)
}

func newRedirect(origin core.Origin, handler messaging.OpsAgent) *redirect {
	c := new(redirect)
	c.agentId = redirectAgentUri(origin)
	c.origin = origin
	c.state = new(redirectState)
	c.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.handler = handler

	return c
}

// String - identity
func (r *redirect) String() string {
	return r.agentId
}

// Uri - agent identifier
func (r *redirect) Uri() string {
	return r.agentId
}

// Message - message the agent
func (r *redirect) Message(m *messaging.Message) {
	messaging.Mux(m, r.ctrlC, nil, nil)
}

// Handle - error handler
func (r *redirect) Handle(status *core.Status, requestId string) *core.Status {
	// TODO : Any operations specific processing ??  If not then forward to handler
	return r.handler.Handle(status, requestId)
}

// AddActivity - add activity
func (r *redirect) AddActivity(agentId string, content any) {
	// TODO : Any operations specific processing ??  If not then forward to handler
	//return a.handler.Handle(status, requestId)
}

// Add - add a shutdown function
//func (a *redirect) Add(f func()) {
//	a.shutdownFunc = messaging.AddShutdown(a.shutdownFunc, f)
//
//}

// Shutdown - shutdown the agent
func (r *redirect) Shutdown() {
	if !r.running {
		return
	}
	r.running = false
	if r.shutdownFunc != nil {
		r.shutdownFunc()
	}
	msg := messaging.NewControlMessage(r.agentId, r.agentId, messaging.ShutdownEvent)
	if r.ctrlC != nil {
		r.ctrlC <- msg
	}
}

// Run - run the agent
func (r *redirect) Run() {
	if r.running {
		return
	}
	go runRedirect(r, newObservation(r.handler), newGuidance(r.handler), newOperations(r.handler))
}

// shutdown - close resources
func (r *redirect) shutdown() {
	close(r.ctrlC)
	r.stopTickers()
}

func (r *redirect) stopTickers() {

}
