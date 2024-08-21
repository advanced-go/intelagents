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
	running      bool
	agentId      string
	origin       core.Origin
	state        *resiliency1.IngressRedirectState
	ticker       *messaging.Ticker
	interval     time.Duration
	ctrlC        chan *messaging.Message
	handler      messaging.OpsAgent
	exchange     *messaging.Exchange
	shutdownFunc func()
}

func redirectAgentUri(origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v.%v", RedirectClass, origin.Region, origin.Zone, origin.Host)
	}
	return fmt.Sprintf("%v:%v.%v.%v.%v", RedirectClass, origin.Region, origin.Zone, origin.SubZone, origin.Host)
}

// newRedirectAgent - create a new lead agent
func newRedirectAgent(origin core.Origin, state *resiliency1.IngressRedirectState, handler messaging.OpsAgent) messaging.Agent {
	return newRedirect(origin, state, handler, redirectDuration)
}

func newRedirect(origin core.Origin, state *resiliency1.IngressRedirectState, handler messaging.OpsAgent, tickerDur time.Duration) *redirect {
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
	go runRedirect(r, redirection, common.Observe, common.Exp, localGuidance)
}

// startup - start tickers
func (r *redirect) startup() {
	r.ticker.Start(-1)
}

// shutdown - close resources
func (r *redirect) shutdown() {
	//msg := messaging.NewControlMessage(r.agentId, r.agentId, messaging.ShutdownEvent)
	close(r.ctrlC)
	r.ticker.Stop()
}

/*
func (r *redirect) updateRedirectPlan(guide *guidance) {
	p, status := guide.redirectPlan(r.handler, r.origin)
	if status.OK() {
		r.state.Location = p.Location
		r.state.Status = p.Status
	}
}

func (r *redirect) updatePercentileSLO(guide *guidance) {
	p, status := guide.percentileSLO(r.handler, r.origin)
	if status.OK() {
		r.state.Percent = p.Percent
		r.state.Latency = p.Latency
		r.state.Minimum = p.Minimum
	}
}


*/

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

func runRedirect(r *redirect, fn *redirectFunc, observe *common.Observation, exp *common.Experience, guide *guidance) {
	for {
		select {
		case <-r.ticker.C():
			r.handler.AddActivity(r.agentId, "onTick()")
			fn.process(r, observe)
			// TODO : based on process need to do the following:
			// 1. Update percentage and send action
			// 2. if status fail, then update redirect
			// 3. if status succeed, then update redirect and set redirect action
			// 4. IF done, then message parent and shutdown
		// control channel
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
