package redirect1

import (
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/intelagents/common2"
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

	// Channels
	lhc *messaging.Channel
	rhc *messaging.Channel

	handler      messaging.OpsAgent
	shutdownFunc func()
}

func redirectAgentUri(origin core.Origin) string {
	return origin.Uri(Class)
}

// newRedirectAgent - create a new redirect agent
func newRedirectAgent(origin core.Origin, state resiliency1.IngressRedirectState, handler messaging.OpsAgent) messaging.Agent {
	return newRedirect(origin, state, handler)
}

func newRedirect(origin core.Origin, state resiliency1.IngressRedirectState, handler messaging.OpsAgent) *redirect {
	r := new(redirect)
	r.agentId = redirectAgentUri(origin)
	r.origin = origin

	//r.state = state
	//r.ticker = messaging.NewTicker(tickerDur)
	//r.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)

	r.rhc = messaging.NewEnabledChannel()
	r.lhc = messaging.NewEnabledChannel()
	r.handler = handler
	return r
}

// String - identity
func (r *redirect) String() string { return r.agentId }

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
func (r *redirect) Run() {
	if r.running {
		return
	}
	go runRedirectRHC(r, redirection, common.Observe, common.Exp, common.Guide)
	go runRedirectLHC(r, common2.Event)
}
