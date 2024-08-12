package ingress1

import (
	"fmt"
	"github.com/advanced-go/guidance/resiliency1"
	core2 "github.com/advanced-go/intelagents/core"
	"github.com/advanced-go/observation/timeseries1"
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
	running bool
	agentId string
	origin  core.Origin
	state   *resiliencyState
	ticker  *messaging.Ticker
	poller  *messaging.Ticker
	//revise       *messaging.Ticker
	ctrlC        chan *messaging.Message
	handler      messaging.OpsAgent
	entries      []timeseries1.Entry
	shutdownFunc func()
}

func resiliencyAgentUri(origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v.%v", Class, origin.Region, origin.Zone, origin.Host)
	}
	return fmt.Sprintf("%v:%v.%v.%v.%v", Class, origin.Region, origin.Zone, origin.SubZone, origin.Host)
}

// NewResiliencyAgent - create a new resiliency agent
func newResiliency(origin core.Origin, handler messaging.OpsAgent) messaging.Agent {
	return newController(origin, handler, resiliencyTickInterval, time.Minute*2, resiliencyReviseInterval)
}

func newController(origin core.Origin, handler messaging.OpsAgent, tickerDur, pollerDur, reviseDur time.Duration) *resiliency {
	c := new(resiliency)
	c.origin = origin
	c.agentId = resiliencyAgentUri(origin)
	c.state = newResiliencyState()
	c.ticker = messaging.NewTicker(tickerDur)
	c.poller = messaging.NewTicker(pollerDur)
	//c.revise = messaging.NewTicker(reviseDur)

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
	go resiliencyRun(c, resiliencyFunc, resiliencyInitFunc, observe, exp, guide)
}

// startup - start tickers
func (c *resiliency) startup() {
	c.ticker.Start(-1)
	c.poller.Start(-1)
	//c.revise.Start(-1)
}

// shutdown - close resources
func (c *resiliency) shutdown() {
	close(c.ctrlC)
	c.ticker.Stop()
	c.poller.Stop()
	//c.revise.Stop()
}

func (c *resiliency) updateTicker(newDuration time.Duration) {
	c.ticker.Start(newDuration)
}

func (c *resiliency) addEntry(entries []timeseries1.Entry) {
	c.entries = append(c.entries, entries...)
}

type resiliencyFn func(c *resiliency, percentile resiliency1.Percentile, observe *observation, exp *experience) ([]timeseries1.Entry, *core.Status)
type resiliencyInitFn func(c *resiliency, observe *observation) *core.Status

// run - ingress resiliency
func resiliencyRun(c *resiliency, ctrlFn resiliencyFn, initFn resiliencyInitFn, observe *observation, exp *experience, guide *guidance) {
	if c == nil {
		return
	}
	// initialize percentile and rate limiting state
	percentile, _ := guide.percentile(c.handler, c.origin, defaultPercentile)
	initFn(c, observe)
	c.startup()
	for {
		// main agent processing
		select {
		case <-c.ticker.C():
			// main : on tick -> observe access -> process inference with percentile -> create action
			c.handler.AddActivity(c.agentId, "onTick")
			entry, status := ctrlFn(c, percentile, observe, exp)
			if status.OK() {
				c.addEntry(entry)
			}
		default:
		}
		// secondary processing
		select {
		case <-c.poller.C():
			c.handler.AddActivity(c.agentId, "onPoll")
			percentile, _ = guide.percentile(c.handler, c.origin, percentile)
		case msg := <-c.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				c.shutdown()
				c.handler.AddActivity(c.agentId, messaging.ShutdownEvent)
				return
			case messaging.DataChangeEvent:
				if msg.IsContentType(core2.ContentTypeProfile) {
					c.handler.AddActivity(c.agentId, "onDataChange() - profile")
					// Process revising the ticker based on the profile.
				} else {
					if msg.IsContentType(core2.ContentTypePercentile) {
						c.handler.AddActivity(c.agentId, "onDataChange() - percentile")
						percentile, _ = guide.percentile(c.handler, c.origin, percentile)
					}
				}
			default:
			}
		default:
		}
	}
}
