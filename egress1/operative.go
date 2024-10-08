package egress1

import (
	"fmt"
	"github.com/advanced-go/intelagents/common"
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
	running      bool
	agentId      string
	origin       core.Origin
	profile      *common.Profile
	agents       *messaging.Exchange
	interval     time.Duration
	ctrlC        chan *messaging.Message
	handler      messaging.OpsAgent
	shutdownFunc func()
}

func FieldOperativeUri(origin core.Origin) string {
	return origin.Uri(FieldOperativeClass)
}

// NewFieldOperative - create a new field operative
func NewFieldOperative(origin core.Origin, profile *common.Profile, handler messaging.OpsAgent) messaging.OpsAgent {
	return newFieldOperative(origin, profile, handler)
}

func newFieldOperative(origin core.Origin, profile *common.Profile, handler messaging.OpsAgent) *fieldOperative {
	f := new(fieldOperative)
	f.agentId = FieldOperativeUri(origin)
	f.origin = origin
	f.profile = profile
	f.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	f.handler = handler
	f.agents = messaging.NewExchange()
	return f
}

// String - identity
func (f *fieldOperative) String() string { return f.Uri() }

// Uri - agent identifier
func (f *fieldOperative) Uri() string { return f.agentId }

// Message - message the agent
func (f *fieldOperative) Message(m *messaging.Message) { f.ctrlC <- m }

// Handle - error handler
func (f *fieldOperative) Handle(status *core.Status) *core.Status {
	return f.handler.Handle(status)
}

// AddActivity - add activity
func (f *fieldOperative) AddActivity(agentId string, content any) {
	f.handler.AddActivity(agentId, content)
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
	f.agents.Shutdown()
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
	go runFieldOperative(f, operative, common.Guide)
}

func (f *fieldOperative) shutdown() {
	close(f.ctrlC)
}

func runFieldOperative(f *fieldOperative, fn *operativeFunc, guide *common.Guidance) {
	if f == nil {
		return
	}
	fn.startup(f, guide)

	for {
		select {
		case msg := <-f.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				f.shutdown()
				f.handler.AddActivity(f.agentId, messaging.ShutdownEvent)
				return
			case messaging.DataChangeEvent:
				f.handler.AddActivity(f.agentId, fmt.Sprintf("%v - %v", msg.Event(), msg.ContentType()))
				if msg.ContentType() == common.ContentTypeEgressConfig {
					fn.onDataChange(f, guide, msg)
				} else {
					forwardDataChangeEvent(f, msg)
				}
			default:
				f.handler.Handle(common.MessageEventErrorStatus(f.agentId, msg))
			}
		default:
		}
	}
}

func forwardDataChangeEvent(f *fieldOperative, msg *messaging.Message) {
	switch msg.ContentType() {
	case common.ContentTypeProfile:
		f.agents.Broadcast(msg)
	default:
		f.handler.Handle(common.MessageContentTypeErrorStatus(f.agentId, msg))
	}
}
