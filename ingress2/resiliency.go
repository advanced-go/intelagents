package ingress2

import (
	"fmt"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/intelagents/common2"
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

	// Left channel
	duration time.Duration
	lhc      *messaging.Channel

	// Right channel

	rhc *messaging.Channel

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
	//go runResiliencyRHC(r, nil, common.Observe, common.Exp, common.Guide)
	go runResiliencyLHC(r, common2.Events)
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
