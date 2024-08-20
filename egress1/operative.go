package egress1

import (
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	FieldOperativeClass = "egress-field-operative1"
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

type fieldOperative struct {
	running bool
	uri     string

	// Assignment
	origin core.Origin

	// Guidance/configuration
	guideVersion         string // Version for authority and egress, helps to stop duplicate updates of egress routes
	processingScheduleId string

	// Routing controllers
	controllers *messaging.Exchange

	interval     time.Duration
	ctrlC        chan *messaging.Message
	handler      messaging.OpsAgent
	shutdownFunc func()
}

func FieldOperativeUri(origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v.%v", FieldOperativeClass, origin.Region, origin.Zone, origin.Host)
	}
	return fmt.Sprintf("%v:%v.%v.%v.%v", FieldOperativeClass, origin.Region, origin.Zone, origin.SubZone, origin.Host)
}

// NewFieldOperative - create a new field operative
func NewFieldOperative(origin core.Origin, handler messaging.OpsAgent) messaging.OpsAgent {
	f := new(fieldOperative)
	f.uri = FieldOperativeUri(origin)
	f.origin = origin

	f.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	f.handler = handler
	f.controllers = messaging.NewExchange()
	return f
}

// String - identity
func (f *fieldOperative) String() string { return f.Uri() }

// Uri - agent identifier
func (f *fieldOperative) Uri() string { return f.uri }

// Message - message the agent
func (f *fieldOperative) Message(m *messaging.Message) { messaging.Mux(m, f.ctrlC, nil, nil) }

// Handle - error handler
func (f *fieldOperative) Handle(status *core.Status, requestId string) *core.Status {
	// TODO : Any operations specific processing ??  If not then forward to handler
	return f.handler.Handle(status, requestId)
}

// AddActivity - add activity
func (f *fieldOperative) AddActivity(agentId string, content any) {
	// TODO : Any operations specific processing ??  If not then forward to handler
	f.handler.AddActivity(agentId, content)
}

// Add - add a shutdown function
func (f *fieldOperative) Add(fn func()) {
	f.shutdownFunc = messaging.AddShutdown(f.shutdownFunc, fn)
}

// Shutdown - shutdown the agent
func (f *fieldOperative) Shutdown() {
	if !f.running {
		return
	}
	f.running = false
	if f.shutdownFunc != nil {
		f.shutdownFunc()
	}
	msg := messaging.NewControlMessage(f.uri, f.uri, messaging.ShutdownEvent)
	if f.ctrlC != nil {
		f.ctrlC <- msg
	}
}

// Run - run the agent
func (f *fieldOperative) Run() {
	if f.running {
		return
	}
	//go runLead(a, newObservation(a.handler), newGuidance(a.handler), newOperations(a.handler))
}
