package ingress2

import (
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/intelagents/common1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	Class           = "ingress-resiliency2"
	defaultDuration = time.Minute * 5
)

// TODO : need to determine a way to increase/decrease the rate of observations if the traffic does not
//         match the profile.

type resiliency struct {
	running bool
	agentId string
	origin  core.Origin

	// Channels
	duration time.Duration
	lhc      *messaging.Channel
	rhc      *messaging.Channel

	handler      messaging.OpsAgent
	shutdownFunc func()
}

func resiliencyAgentUri(origin core.Origin) string {
	return origin.Uri(Class)
}

// NewResiliencyAgent - create a new resiliency agent
func NewResiliencyAgent(origin core.Origin, profile *common.Profile, handler messaging.OpsAgent) messaging.Agent {
	return newResiliency(origin, profile, handler)
}

func newResiliency(origin core.Origin, profile *common.Profile, handler messaging.OpsAgent) *resiliency {
	r := new(resiliency)
	r.origin = origin
	r.agentId = resiliencyAgentUri(origin)

	// Left channel
	r.duration = defaultDuration
	if profile != nil {
		r.duration = profile.ResiliencyDuration(-1)
	}
	r.lhc = messaging.NewEnabledChannel()

	// Right channel
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
	if m.IsContentType(common.ContentTypeProfile) || m.Channel() == messaging.ChannelLeft {
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
	go runResiliencyRHC(r, common1.IngressExperience)
	go runResiliencyLHC(r, common1.TimeseriesEvents)
	r.running = true
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
