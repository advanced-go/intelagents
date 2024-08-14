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
	Class = "ingress-resiliency1"
)

var (
	defaultPercentile = resiliency1.Percentile{Percent: 99, Latency: 2000}
)

type resiliencyState struct {
	rateLimit float64
	rateBurst int
}

func newResiliencyState() *resiliencyState {
	l := new(resiliencyState)
	l.rateLimit = -1
	l.rateBurst = -1
	return l
}

type resiliency struct {
	running      bool
	agentId      string
	origin       core.Origin
	state        *resiliencyState
	profile      *common.Profile
	percentile   *resiliency1.Percentile
	ticker       *messaging.Ticker
	ctrlC        chan *messaging.Message
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
func newResiliencyAgent(origin core.Origin, profile *common.Profile, percentile *resiliency1.Percentile, handler messaging.OpsAgent) messaging.Agent {
	return newResiliency(origin, profile, percentile, handler)
}

func newResiliency(origin core.Origin, profile *common.Profile, percentile *resiliency1.Percentile, handler messaging.OpsAgent) *resiliency {
	r := new(resiliency)
	r.origin = origin
	r.agentId = resiliencyAgentUri(origin)
	r.state = newResiliencyState()
	r.profile = profile
	r.percentile = percentile
	r.ticker = messaging.NewTicker(profile.ResiliencyDuration(-1))
	r.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	r.handler = handler
	return r
}

// String - identity
func (r *resiliency) String() string { return r.Uri() }

// Uri - agent identifier
func (r *resiliency) Uri() string { return r.agentId }

// Message - message the agent
func (r *resiliency) Message(m *messaging.Message) { messaging.Mux(m, r.ctrlC, nil, nil) }

// Add - add a shutdown function
func (r *resiliency) Add(f func()) { r.shutdownFunc = messaging.AddShutdown(r.shutdownFunc, f) }

// Run - run the agent
func (r *resiliency) Run() {
	if r.running {
		return
	}
	go resiliencyRun(r, resilience, observe, exp)
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
	if r.ctrlC != nil {
		r.ctrlC <- msg
	}
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

// run - ingress resiliency
func resiliencyRun(r *resiliency, fn *resiliencyFunc, observe *observation, exp *experience) {
	percentile, _ := fn.startup(r, exp)

	for {
		// main agent processing : on tick -> observe access -> process inference with percentile -> create action
		select {
		case <-r.ticker.C():
			r.handler.AddActivity(r.agentId, "onTick")
			fn.process(r, percentile, observe, exp)
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
				if msg.IsContentType(common.ContentTypeProfile) {
					r.handler.AddActivity(r.agentId, "onDataChange() - profile")
					if p, ok := msg.Body.(*common.Profile); ok {
						r.reviseTicker(p.ResiliencyDuration(-1))
					} else {
						r.handler.Handle(common.ProfileTypeErrorStatus(msg.Body), "")
					}
				} else {
					if msg.IsContentType(common.ContentTypePercentile) {
						r.handler.AddActivity(r.agentId, "onDataChange() - percentile")
						if p, ok := msg.Body.(*resiliency1.Percentile); ok {
							percentile = p
						} else {
							r.handler.Handle(common.PercentileTypeErrorStatus(msg.Body), "")
						}
					}
				}
			default:
			}
		default:
		}
	}
}
