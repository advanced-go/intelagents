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
	RedirectClass    = "ingress-redirect1"
	redirectDuration = time.Second * 60
)

// Responsibilities:
//  1. Startup + Restart Events

type redirect struct {
	running      bool
	agentId      string
	origin       core.Origin
	state        resiliency1.IngressRedirectState
	ticker       *messaging.Ticker
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

// newRedirectAgent - create a new redirect agent
func newRedirectAgent(origin core.Origin, state resiliency1.IngressRedirectState, handler messaging.OpsAgent) messaging.Agent {
	return newRedirect(origin, state, handler, redirectDuration)
}

func newRedirect(origin core.Origin, state resiliency1.IngressRedirectState, handler messaging.OpsAgent, tickerDur time.Duration) *redirect {
	r := new(redirect)
	r.agentId = redirectAgentUri(origin)
	r.origin = origin
	r.state = state
	r.ticker = messaging.NewTicker(tickerDur)
	r.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	r.handler = handler
	return r
}

// String - identity
func (r *redirect) String() string { return r.agentId }

// Uri - agent identifier
func (r *redirect) Uri() string { return r.agentId }

// Message - message the agent
func (r *redirect) Message(m *messaging.Message) { messaging.Mux(m, r.ctrlC, nil, nil) }

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
	if r.ctrlC != nil {
		r.ctrlC <- msg
	}
}

// Run - run the agent
func (r *redirect) Run() {
	if r.running {
		return
	}
	go runRedirect(r, redirection, common.Observe, common.Guide)
}

// startup - start tickers
func (r *redirect) startup() {
	r.ticker.Start(-1)
}

// shutdown - close resources
func (r *redirect) shutdown() {
	close(r.ctrlC)
	r.ticker.Stop()
}

func (r *redirect) updatePercentage() {
	switch r.state.Percentage {
	case 0:
		r.state.Percentage = 10
	case 10:
		r.state.Percentage = 20
	case 20:
		r.state.Percentage = 40
	case 40:
		r.state.Percentage = 70
	case 70:
		r.state.Percentage = 100
	default:
	}
}

func runRedirect(r *redirect, fn *redirectFunc, observe *common.Observation, guide *common.Guidance) {
	r.startup()
	for {
		select {
		case <-r.ticker.C():
			r.handler.AddActivity(r.agentId, "onTick()")
			completed, status := fn.process(r, observe, guide)
			if completed {
				fn.update(r, guide, status.OK())
				r.handler.Message(messaging.NewControlMessage("", r.agentId, RedirectCompletedEvent))
				r.shutdown()
				return
			}
		case msg := <-r.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				r.shutdown()
				r.handler.AddActivity(r.agentId, messaging.ShutdownEvent)
				return
			default:
			}
		default:
		}
	}
}
