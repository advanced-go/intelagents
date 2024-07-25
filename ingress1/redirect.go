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
type redirect struct {
	running bool
	uri     string

	// Assignment
	origin core.Origin

	// Guidance/configuration
	guideVersion         string // Version for authority and egress, helps to stop duplicate updates of egress routes
	processingScheduleId string

	// Dependency processing
	dependencyUpdates    bool
	dependencyScheduleId string
	dependencyAgent      messaging.Agent

	// Routing controllers
	controllers *messaging.Exchange

	interval     time.Duration
	ctrlC        chan *messaging.Message
	handler      messaging.OpsAgent
	shutdownFunc func()
}

func RedirectAgentUri(origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v.%v", RedirectClass, origin.Region, origin.Zone, origin.Host)
	}
	return fmt.Sprintf("%v:%v.%v.%v.%v", RedirectClass, origin.Region, origin.Zone, origin.SubZone, origin.Host)
}

// NewRedirectAgent - create a new lead agent
func NewRedirectAgent(origin core.Origin, handler messaging.OpsAgent) messaging.OpsAgent {
	c := new(redirect)
	c.uri = LeadAgentUri(origin)
	c.origin = origin

	c.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.handler = handler
	//c.dependencyAgent = dependency1.NewDependencyAgent(origin, c)
	c.controllers = messaging.NewExchange()
	return c
}

// String - identity
func (a *redirect) String() string {
	return a.uri
}

// Uri - agent identifier
func (a *redirect) Uri() string {
	return a.uri
}

// Message - message the agent
func (a *redirect) Message(m *messaging.Message) {
	messaging.Mux(m, a.ctrlC, nil, nil)
}

// Handle - error handler
func (a *redirect) Handle(status *core.Status, requestId string) *core.Status {
	// TODO : Any operations specific processing ??  If not then forward to handler
	return a.handler.Handle(status, requestId)
}

// AddActivity - add activity
func (a *redirect) AddActivity(agentId string, content any) {
	// TODO : Any operations specific processing ??  If not then forward to handler
	//return a.handler.Handle(status, requestId)
}

// Add - add a shutdown function
func (a *redirect) Add(f func()) {
	a.shutdownFunc = messaging.AddShutdown(a.shutdownFunc, f)

}

// Shutdown - shutdown the agent
func (a *redirect) Shutdown() {
	if !a.running {
		return
	}
	a.running = false
	if a.shutdownFunc != nil {
		a.shutdownFunc()
	}
	msg := messaging.NewControlMessage(a.uri, a.uri, messaging.ShutdownEvent)
	if a.ctrlC != nil {
		a.ctrlC <- msg
	}
}

// Run - run the agent
func (a *redirect) Run() {
	if a.running {
		return
	}
	go runRedirect(a, newObservation(a.handler), newGuidance(a.handler), newOperations(a.handler))
}

// shutdown - close resources
func (a *redirect) shutdown() {
	close(a.ctrlC)
	a.stopTickers()
}

func (a *redirect) stopTickers() {

}
