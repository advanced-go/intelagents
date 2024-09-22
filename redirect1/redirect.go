package redirect1

import (
	"github.com/advanced-go/guidance/redirect1"
	"github.com/advanced-go/intelagents/common1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	Class                  = "ingress-redirect1"
	redirectDuration       = time.Second * 60
	RedirectCompletedEvent = "event:redirect-completed"
)

// Responsibilities:
//  1. Startup + Restart Events

type redirect struct {
	running bool
	agentId string
	origin  core.Origin
	state   redirect1.IngressEntry

	// Channels
	lhc *messaging.Channel
	rhc *messaging.Channel

	handler      messaging.OpsAgent
	shutdownFunc func()
}

func redirectAgentUri(origin core.Origin) string {
	return origin.Uri(Class)
}

// NewRedirectAgent - create a new redirect agent
func NewRedirectAgent(origin core.Origin, handler messaging.OpsAgent) messaging.Agent {
	return newRedirect(origin, handler)
}

func newRedirect(origin core.Origin, handler messaging.OpsAgent) *redirect {
	r := new(redirect)
	r.agentId = redirectAgentUri(origin)
	r.origin = origin
	r.rhc = messaging.NewEnabledChannel()
	r.lhc = messaging.NewEnabledChannel()
	r.handler = handler
	return r
}

// String - identity
func (r *redirect) String() string { return r.Uri() }

// Uri - agent identifier
func (r *redirect) Uri() string { return r.agentId }

// Message - message the agent
func (r *redirect) Message(m *messaging.Message) {
	// How to determine which channel??
	if m.Channel() == messaging.ChannelLeft {
		r.lhc.Send(m)
	} else {
		r.rhc.Send(m)
	}
}

// Add - add a shutdown function
func (r *redirect) Add(f func()) { r.shutdownFunc = messaging.AddShutdown(r.shutdownFunc, f) }

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
	r.rhc.Enable()
	r.rhc.Send(msg)
	r.lhc.Enable()
	r.lhc.Send(msg)

}

// Run - run the agent
// TODO: notification/error message if no redirect is found?
func (r *redirect) Run() {
	if r.running {
		return
	}
	var status *core.Status
	r.state, status = common1.RedirectGuidance.Ingress(r.handler, r.origin)
	if !status.OK() {
		// Remove agent from exchange if registered
		if r.shutdownFunc != nil {
			r.shutdownFunc()
		}
		return
	}
	// TODO : If the redirect has a start time configured, then process that with a timer and go routing
	go runRedirectRHC(r, common1.IngressExperience, common1.RedirectGuidance)
	go runRedirectLHC(r, common1.TimeseriesEvents)
	r.running = true
}
