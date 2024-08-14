package ingress1

import (
	"fmt"
	"github.com/advanced-go/guidance/resiliency1"
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
type lead struct {
	running      bool
	agentId      string
	origin       core.Origin
	profile      *common.Profile
	percentile   *resiliency1.Percentile
	resiliency   messaging.Agent
	redirect     messaging.Agent
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
func NewLeadAgent(origin core.Origin, profile *common.Profile, percentile *resiliency1.Percentile, handler messaging.OpsAgent) messaging.OpsAgent {
	l := new(lead)
	l.agentId = LeadAgentUri(origin)
	l.origin = origin
	l.profile = profile
	l.percentile = percentile

	l.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	l.handler = handler

	l.resiliency = newResiliencyAgent(origin, profile, percentile, l)
	l.redirect = newRedirectAgent(origin, l)
	return l
}

// String - identity
func (l *lead) String() string {
	return l.agentId
}

// Uri - agent identifier
func (l *lead) Uri() string {
	return l.agentId
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
	l.resiliency.Shutdown()
	l.redirect.Shutdown()
	msg := messaging.NewControlMessage(l.agentId, l.agentId, messaging.ShutdownEvent)
	if l.ctrlC != nil {
		l.ctrlC <- msg
	}
}

// Run - run the agent
func (l *lead) Run() {
	if l.running {
		return
	}
	go leadRun(l, guide)
}

// shutdown - close resources
func (l *lead) shutdown() {
	close(l.ctrlC)

}

func leadRun(l *lead, guide *guidance) {
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
