package ingress1

import (
	"fmt"
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

const (
	FieldOperativeClass    = "ingress-field-operative1"
	RedirectCompletedEvent = "event:redirect-completed"
)

// Responsibilities:
//  1. Startup + Restart Events
//

type fieldOperative struct {
	running      bool
	agentId      string
	origin       core.Origin
	profile      *common.Profile
	state        *resiliency1.IngressRedirectState
	resiliency   messaging.Agent
	redirect     messaging.Agent
	ctrlC        chan *messaging.Message
	handler      messaging.OpsAgent
	shutdownFunc func()
}

func FieldOperativeUri(origin core.Origin) string {
	return origin.Uri(FieldOperativeClass)
}

// NewFieldOperative - create a new field operative
func NewFieldOperative(origin core.Origin, profile *common.Profile, handler messaging.OpsAgent) messaging.OpsAgent {
	return newFieldOperative(origin, profile, nil, handler)
}

func newFieldOperative(origin core.Origin, profile *common.Profile, resilience messaging.Agent, handler messaging.OpsAgent) *fieldOperative {
	f := new(fieldOperative)
	f.agentId = FieldOperativeUri(origin)
	f.origin = origin
	f.state = resiliency1.NewIngressRedirectState()
	f.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	f.handler = handler
	if resilience == nil {
		f.resiliency = newResiliencyAgent(origin, profile, f)
	} else {
		f.resiliency = resilience
	}
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
	return f.handler.Handle(status, requestId)
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
	// Removes agent from its exchange if registered
	if f.shutdownFunc != nil {
		f.shutdownFunc()
	}
	f.resiliency.Shutdown()
	if f.redirect != nil {
		f.redirect.Shutdown()
	}
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
	f.resiliency.Run()
	go runFieldOperative(f, operative, localGuidance)
}

func (f *fieldOperative) shutdown() {
	close(f.ctrlC)
}

func runFieldOperative(f *fieldOperative, fn *operativeFunc, guide *guidance) {
	if f == nil {
		return
	}
	fn.processRedirect(f, fn, guide)

	for {
		select {
		case msg := <-f.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				f.shutdown()
				f.handler.AddActivity(f.agentId, messaging.ShutdownEvent)
				return
			case RedirectCompletedEvent:
				f.redirect = nil
			case messaging.DataChangeEvent:
				f.handler.AddActivity(f.agentId, fmt.Sprintf("%v - %v", msg.Event(), msg.ContentType()))
				if msg.ContentType() == common.ContentTypeRedirectPlan {
					fn.processRedirect(f, fn, guide)
				} else {
					forwardDataChangeEvent(f, msg)
				}
			default:
				f.handler.Handle(common.MessageEventErrorStatus(f.agentId, msg), "")
			}
		default:
		}
	}
}

func forwardDataChangeEvent(f *fieldOperative, msg *messaging.Message) {
	switch msg.Header.Get(messaging.ContentType) {
	case common.ContentTypeProfile:
		if p := common.GetProfile(f.handler, f.agentId, msg); p != nil && p.Next().IsScaleUp() {
			// Need to send data change event if the next window of the profile is scaling up. This means that
			// a periodic routine done to update SLOs has been completed in the Off-Peak window.
			m := messaging.NewControlMessage(f.resiliency.Uri(), f.agentId, messaging.DataChangeEvent)
			m.SetContentType(common.ContentTypePercentileSLO)
			f.resiliency.Message(m)
		}
		f.resiliency.Message(msg)
	default:
		f.handler.Handle(common.MessageContentTypeErrorStatus(f.agentId, msg), "")
	}
}
