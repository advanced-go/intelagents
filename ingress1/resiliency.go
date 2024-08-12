package ingress1

import (
	"fmt"
	"github.com/advanced-go/guidance/resiliency1"
	core2 "github.com/advanced-go/intelagents/core"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	Class                    = "ingress-resiliency1"
	resiliencyTickInterval   = time.Minute * 2
	resiliencyReviseInterval = time.Hour * 1
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
func newResiliencyAgent(origin core.Origin, handler messaging.OpsAgent) messaging.Agent {
	return newResiliency(origin, handler, resiliencyTickInterval, time.Minute*2, resiliencyReviseInterval)
}

func newResiliency(origin core.Origin, handler messaging.OpsAgent, tickerDur, pollerDur, reviseDur time.Duration) *resiliency {
	c := new(resiliency)
	c.origin = origin
	c.agentId = resiliencyAgentUri(origin)
	c.state = newResiliencyState()
	c.ticker = messaging.NewTicker(tickerDur)
	c.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.handler = handler
	return c
}

// String - identity
func (c *resiliency) String() string { return c.agentId }

// Uri - agent identifier
func (c *resiliency) Uri() string { return c.agentId }

// Message - message the agent
func (c *resiliency) Message(m *messaging.Message) { messaging.Mux(m, c.ctrlC, nil, nil) }

// Add - add a shutdown function
func (c *resiliency) Add(f func()) { c.shutdownFunc = messaging.AddShutdown(c.shutdownFunc, f) }

// Shutdown - shutdown the agent
func (c *resiliency) Shutdown() {
	if !c.running {
		return
	}
	c.running = false
	if c.shutdownFunc != nil {
		c.shutdownFunc()
	}
	msg := messaging.NewControlMessage(c.agentId, c.agentId, messaging.ShutdownEvent)
	if c.ctrlC != nil {
		c.ctrlC <- msg
	}
}

// Run - run the agent
func (c *resiliency) Run() {
	if c.running {
		return
	}
	go resiliencyRun(c, resilience, observe, exp, guide)
}

// startup - start tickers
func (c *resiliency) startup() {
	c.ticker.Start(-1)
}

// shutdown - close resources
func (c *resiliency) shutdown() {
	close(c.ctrlC)
	c.ticker.Stop()
}

func (c *resiliency) updateTicker(newDuration time.Duration) {
	c.ticker.Start(newDuration)
}

// run - ingress resiliency
func resiliencyRun(r *resiliency, fn *resiliencyFunc, observe *observation, exp *experience, guide *guidance) {
	if r == nil {
		return
	}
	// initialize percentile and rate limiting state
	percentile, _ := guide.percentile(r.handler, r.origin, defaultPercentile)
	fn.init(r, exp)
	r.startup()
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
				if msg.IsContentType(core2.ContentTypeProfile) {
					r.handler.AddActivity(r.agentId, "onDataChange() - profile")
					// Process revising the ticker based on the profile.
				} else {
					if msg.IsContentType(core2.ContentTypePercentile) {
						r.handler.AddActivity(r.agentId, "onDataChange() - percentile")
						percentile, _ = guide.percentile(r.handler, r.origin, percentile)
					}
				}
			default:
			}
		default:
		}
	}
}
