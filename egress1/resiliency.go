package egress1

import (
	"fmt"
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	ResiliencyClass = "egress-resiliency1"
)

// TODO : need to determine a way to increase/decrease the rate of observations if the traffic does not
//         match the profile.

type resiliency struct {
	running      bool
	agentId      string
	origin       core.Origin
	state        resiliency1.EgressState
	profile      *common.Profile
	ticker       *messaging.Ticker
	ctrlC        chan *messaging.Message
	handler      messaging.OpsAgent
	shutdownFunc func()
}

func ResiliencyAgentUri(origin core.Origin) string {
	return origin.Uri(ResiliencyClass)
}

// newResiliencyAgent - create a new Resiliency agent
func newResiliencyAgent(origin core.Origin, profile *common.Profile, state resiliency1.EgressState, handler messaging.OpsAgent) messaging.OpsAgent {
	return newResiliency(origin, profile, state, handler)
}

func newResiliency(origin core.Origin, profile *common.Profile, state resiliency1.EgressState, handler messaging.OpsAgent) *resiliency {
	c := new(resiliency)
	c.agentId = ResiliencyAgentUri(origin)
	c.origin = origin
	c.profile = profile
	c.state = state
	c.ticker = messaging.NewTicker(profile.ResiliencyDuration(-1))
	c.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.handler = handler
	return c
}

// String - identity
func (r *resiliency) String() string { return r.Uri() }

// Uri - agent identifier
func (r *resiliency) Uri() string { return r.agentId }

// Message - message the agent
func (r *resiliency) Message(m *messaging.Message) { r.ctrlC <- m }

// Handle - error handler
func (r *resiliency) Handle(status *core.Status) *core.Status {
	return r.handler.Handle(status)
}

// AddActivity - add activity
func (r *resiliency) AddActivity(agentId string, content any) {
	r.handler.AddActivity(agentId, content)
}

// Add - add a shutdown function
func (r *resiliency) Add(f func()) { r.shutdownFunc = messaging.AddShutdown(r.shutdown, f) }

// Shutdown - shutdown the agent
func (r *resiliency) Shutdown() {
	if !r.running {
		return
	}
	r.running = false
	if r.shutdownFunc != nil {
		r.shutdown()
	}
	msg := messaging.NewControlMessage(r.agentId, r.agentId, messaging.ShutdownEvent)
	if r.ctrlC != nil {
		r.ctrlC <- msg
	}
}

// Run - run the agent
func (r *resiliency) Run() {
	if r.running {
		return
	}
	go runResiliency(r, resilience, common.Observe, common.Exp, common.Guide)
}

func (r *resiliency) startup() {
	r.ticker.Start(-1)
}

func (r *resiliency) shutdown() {
	close(r.ctrlC)
	r.ticker.Stop()
}

func (r *resiliency) reviseTicker(newDuration time.Duration) {
	r.ticker.Start(newDuration)
}

// run - egress resiliency
func runResiliency(r *resiliency, fn *resiliencyFunc, observe *common.Observation, exp *common.Experience, guide *common.Guidance) {
	r.startup()

	for {
		// main agent processing : on tick -> observe access -> process inference with percentile -> create action
		select {
		case <-r.ticker.C():
			r.handler.AddActivity(r.agentId, "onTick")
			fn.process(r, observe, exp, guide)
		default:
		}
		// control channel processing
		select {
		case msg := <-r.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				r.shutdown()
				r.handler.AddActivity(r.agentId, messaging.ShutdownEvent)
				return
			case messaging.DataChangeEvent:
				r.handler.AddActivity(r.agentId, fmt.Sprintf("%v - %v", msg.Event(), msg.ContentType()))
				processDataChangeEvent(r, msg, exp)
			default:
				r.handler.Handle(common.MessageEventErrorStatus(r.agentId, msg))
			}
		default:
		}
	}
}

func processDataChangeEvent(r *resiliency, msg *messaging.Message, exp *common.Experience) {
	switch msg.ContentType() {
	case common.ContentTypeProfile:
		if p := common.GetProfile(r.handler, r.agentId, msg); p != nil {
			r.reviseTicker(p.ResiliencyDuration(-1))
		}
		return
	case common.ContentTypeEgressConfig:
		// Only a configuration update, insert and delete are handled by field operative
		if p, ok := msg.Body.(resiliency1.EgressConfig); ok {
			// Threshold is safe to change as there is no shared external state via actions
			r.state.FailoverThreshold = p.FailoverThreshold
			// If the scope changes then initialize all the following state:
			// 1. Current action
			// 2. Current location and percentage.
			if r.state.FailoverScope != p.FailoverScope {
				exp.ResetRoutingAction(r.handler, p.Origin(), r.agentId)
				r.state.FailoverScope = p.FailoverScope
				r.state.Location = ""
				r.state.Percentage = -1
			}
			return
		}
	default:
	}
	r.handler.Handle(common.MessageContentTypeErrorStatus(r.agentId, msg))
}
