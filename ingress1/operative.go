package ingress1

import (
	"fmt"
	"github.com/advanced-go/intelagents/common"
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
type fieldOperative struct {
	running      bool
	agentId      string
	origin       core.Origin
	resiliency   messaging.Agent
	redirect     messaging.Agent
	interval     time.Duration
	ctrlC        chan *messaging.Message
	handler      messaging.OpsAgent
	shutdownFunc func()
}

func FieldOperativeUri(origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v.%v", LeadClass, origin.Region, origin.Zone, origin.Host)
	}
	return fmt.Sprintf("%v:%v.%v.%v.%v", LeadClass, origin.Region, origin.Zone, origin.SubZone, origin.Host)
}

// NewFieldOperative - create a new lead agent
func NewFieldOperative(origin core.Origin, profile *common.Profile, handler messaging.OpsAgent) messaging.OpsAgent {
	f := new(fieldOperative)
	f.agentId = FieldOperativeUri(origin)
	f.origin = origin
	f.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	f.handler = handler
	f.resiliency = newResiliencyAgent(origin, profile, f)
	f.redirect = newRedirectAgent(origin, profile, f)
	return f
}

// String - identity
func (f *fieldOperative) String() string { return f.Uri() }

// Uri - agent identifier
func (f *fieldOperative) Uri() string { return f.agentId }

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
	//return a.handler.Handle(status, requestId)
}

// Add - add a shutdown function
func (f *fieldOperative) Add(fn func()) { f.shutdownFunc = messaging.AddShutdown(f.shutdownFunc, fn) }

// Shutdown - shutdown the agent
func (f *fieldOperative) Shutdown() {
	if !f.running {
		return
	}
	f.running = false
	if f.shutdownFunc != nil {
		f.shutdownFunc()
	}
	f.resiliency.Shutdown()
	f.redirect.Shutdown()
	msg := messaging.NewControlMessage(f.agentId, f.agentId, messaging.ShutdownEvent)
	if f.ctrlC != nil {
		f.ctrlC <- msg
	}
}

// Run - run the agent
func (f *fieldOperative) Run() {
	if f.running {
		return
	}
	go runFieldOperative(f, guide)
}

// shutdown - close resources
func (f *fieldOperative) shutdown() {
	close(f.ctrlC)

}

func runFieldOperative(l *fieldOperative, guide *guidance) {
	if l == nil {
		return
	}
	//entry, status := guide.controllers(l.handler, l.origin)
	//if entry.EntryId != 0 || status != nil {
	//}
	for {
		select {
		case msg := <-l.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				l.resiliency.Shutdown()
				l.redirect.Shutdown()
				l.shutdown()
				return
			case messaging.HostStartupEvent:
				//if
				//l.controller.Message(msg)
			case messaging.ChangesetApplyEvent:
				//l.controller.Message(msg)
			case messaging.ChangesetRollbackEvent:
				//l.controller.Message(msg)
			default:
			}
		default:
		}
	}
}
