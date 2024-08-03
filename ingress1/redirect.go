package ingress1

import (
	"fmt"
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

type redirectState struct {
	host     string
	location string
	percent  int
}

type redirect struct {
	running bool
	agentId string

	// Assignment
	origin       core.Origin
	state        *redirectState
	ticker       *messaging.Ticker
	interval     time.Duration
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

// newRedirectAgent - create a new lead agent
func newRedirectAgent(origin core.Origin, handler messaging.OpsAgent) messaging.Agent {
	return newRedirect(origin, handler, redirectDuration)
}

func newRedirect(origin core.Origin, handler messaging.OpsAgent, tickerDur time.Duration) *redirect {
	c := new(redirect)
	c.agentId = redirectAgentUri(origin)
	c.origin = origin
	c.state = new(redirectState)
	c.ticker = messaging.NewTicker(tickerDur)
	c.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.handler = handler

	return c
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
	go runRedirect(r, observe, guide)
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

func runRedirect(r *redirect, observe *observation, guide *guidance) {
	if r == nil {
		return
	}
	r.startup()

	for {
		select {
		case <-r.ticker.C():

		// control channel
		case msg := <-r.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				r.shutdown()
				r.handler.AddActivity(r.agentId, messaging.ShutdownEvent)
				return
			case messaging.RestartEvent:
				// How to handle duplicate restart events as multiple pods will be sending restarts?
				// TODO : read/update from guidance
				//m := messaging.NewControlMessage(a.dependencyAgent.Uri(),a.uri,messaging.RestartEvent)
				//a.dependencyAgent.Message(m)
				//a.controllers.Broadcast(m)

			// Duplicates?? Should not happen.
			// Changesets only contain the changes
			case messaging.ChangesetApplyEvent:
			case messaging.ChangesetRollbackEvent:
				// TODO : apply and rollback changeset
			default:
			}
		default:
		}
	}
}
