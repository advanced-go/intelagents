package ingress2

import (
	"fmt"
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	Class = "ingress-resiliency2"
)

// TODO : need to determine a way to increase/decrease the rate of observations if the traffic does not
//         match the profile.

type resiliency struct {
	running      bool
	agentId      string
	origin       core.Origin
	state        resiliency1.IngressResiliencyState
	profile      *common.Profile
	ticker       *messaging.Ticker
	lhc          *messaging.Channel
	rhc          *messaging.Channel
	handler      messaging.OpsAgent
	shutdownFunc func()
}

func resiliencyAgentUri(origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v.%v", Class, origin.Region, origin.Zone, origin.Host)
	}
	return fmt.Sprintf("%v:%v.%v.%v.%v", Class, origin.Region, origin.Zone, origin.SubZone, origin.Host)
}

// NewResiliencyAgent - create a new resiliency agent
func NewResiliencyAgent(origin core.Origin, profile *common.Profile, handler messaging.OpsAgent) messaging.Agent {
	return newResiliency(origin, profile, handler)
}

func newResiliency(origin core.Origin, profile *common.Profile, handler messaging.OpsAgent) *resiliency {
	r := new(resiliency)
	r.origin = origin
	r.agentId = resiliencyAgentUri(origin)
	resiliency1.NewIngressResiliencyState(&r.state)
	r.profile = profile
	r.ticker = messaging.NewTicker(profile.ResiliencyDuration(-1))
	r.lhc = messaging.NewEnabledChannel()
	r.rhc = messaging.NewEnabledChannel()
	r.handler = handler
	return r
}

// String - identity
func (r *resiliency) String() string { return r.Uri() }

// Uri - agent identifier
func (r *resiliency) Uri() string { return r.agentId }

// Message - message the agent
func (r *resiliency) Message(m *messaging.Message) {
	if m == nil {
		return
	}
	// Specifically for the lhc or profile content
	if m.IsContentType(common.ContentTypeProfile) || m.Channel() == messaging.ChannelLeftHemisphere {
		r.lhc.C <- m
	} else {
		r.rhc.C <- m
	}
}

// Add - add a shutdown function
func (r *resiliency) Add(f func()) { r.shutdownFunc = messaging.AddShutdown(r.shutdownFunc, f) }

// Run - run the agent
func (r *resiliency) Run() {
	if r.running {
		return
	}
	go runResiliencyRHC(r, nil, common.Observe, common.Exp, common.Guide)
}

// Shutdown - shutdown the agent
func (r *resiliency) Shutdown() {
	if !r.running {
		return
	}
	r.running = false
	if r.shutdownFunc != nil {
		r.shutdownFunc()
	}
	msg := messaging.NewControlMessage(r.agentId, r.agentId, messaging.ShutdownEvent)
	r.lhc.Enable()
	r.lhc.C <- msg
	r.rhc.Enable()
	r.rhc.C <- msg
}

func (r *resiliency) startup() {
	r.ticker.Start(-1)
}

func (r *resiliency) shutdown() {
	r.lhc.Close()
	r.rhc.Close()
	r.ticker.Stop()
}

func (r *resiliency) reviseTicker(newDuration time.Duration) {
	r.ticker.Start(newDuration)
}

func (r *resiliency) updatePercentileSLO(guide *common.Guidance) {
	p, status := guide.PercentileSLO(r.handler, r.origin)
	if status.OK() {
		r.state.Percent = p.Percent
		r.state.Latency = p.Latency
		r.state.Minimum = p.Minimum
	}
}

// run - ingress resiliency for the RHC
func runResiliencyRHC(r *resiliency, fn *resiliencyFunc, observe *common.Observation, exp *common.Experience, guide *common.Guidance) {
	fn.startup(r, guide)

	for {
		// main agent processing : on tick -> observe access -> process inference with percentile -> create action
		select {
		case <-r.ticker.C():
			r.handler.AddActivity(r.agentId, "onTick")
			fn.process(r, observe, exp)
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
				processDataChangeEvent(r, msg, guide)
			default:
				r.handler.Handle(common.MessageEventErrorStatus(r.agentId, msg))
			}
		default:
		}
	}
}

// run - ingress resiliency for the LHC
func runResiliencyLHC(r *resiliency, fn *resiliencyFunc, observe *common.Observation, exp *common.Experience, guide *common.Guidance) {
	fn.startup(r, guide)

	for {
		// main agent processing : on tick -> observe access -> process inference with percentile -> create action
		select {
		case <-r.ticker.C():
			r.handler.AddActivity(r.agentId, "onTick")
			fn.process(r, observe, exp)
		default:
		}
		// control channel processing
		select {
		case msg := <-r.lhc.C:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				r.shutdown()
				r.handler.AddActivity(r.agentId, messaging.ShutdownEvent)
				return
			case messaging.DataChangeEvent:
				r.handler.AddActivity(r.agentId, fmt.Sprintf("%v - %v", msg.Event(), msg.ContentType()))
				processDataChangeEvent(r, msg, guide)
			default:
				r.handler.Handle(common.MessageEventErrorStatus(r.agentId, msg))
			}
		default:
		}
	}
}

func processDataChangeEvent(r *resiliency, msg *messaging.Message, guide *common.Guidance) {
	switch msg.ContentType() {
	case common.ContentTypeProfile:
		// GetProfile errors on cast
		if p := common.GetProfile(r.handler, r.agentId, msg); p != nil {
			r.reviseTicker(p.ResiliencyDuration(-1))
		}
		return
	case common.ContentTypePercentileSLO:
		r.updatePercentileSLO(guide)
		return
	default:
		r.handler.Handle(common.MessageContentTypeErrorStatus(r.agentId, msg))
	}
}
