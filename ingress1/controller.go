package ingress1

import (
	"fmt"
	"github.com/advanced-go/guidance/percentile1"
	core2 "github.com/advanced-go/intelagents/core"
	"github.com/advanced-go/observation/timeseries1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	Class                    = "ingress-controller1"
	controllerTickInterval   = time.Minute * 2
	controllerReviseInterval = time.Hour * 1
)

var (
	defaultPercentile = percentile1.Entry{Percent: 99, Latency: 2000}
)

type controllerState struct {
	rateLimit float64
	rateBurst int
}

func newControllerState() *controllerState {
	l := new(controllerState)
	l.rateLimit = -1
	l.rateBurst = -1
	return l
}

type controller struct {
	running bool
	agentId string
	origin  core.Origin
	state   *controllerState
	ticker  *messaging.Ticker
	poller  *messaging.Ticker
	//revise       *messaging.Ticker
	ctrlC        chan *messaging.Message
	dataC        chan *messaging.Message
	handler      messaging.OpsAgent
	entries      []timeseries1.Entry
	shutdownFunc func()
}

func controllerAgentUri(origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v.%v", Class, origin.Region, origin.Zone, origin.Host)
	}
	return fmt.Sprintf("%v:%v.%v.%v.%v", Class, origin.Region, origin.Zone, origin.SubZone, origin.Host)
}

// NewControllerAgent - create a new controller agent
func newControllerAgent(origin core.Origin, handler messaging.OpsAgent) messaging.Agent {
	return newController(origin, handler, controllerTickInterval, percentile1.PercentilePollingDuration, controllerReviseInterval)
}

func newController(origin core.Origin, handler messaging.OpsAgent, tickerDur, pollerDur, reviseDur time.Duration) *controller {
	c := new(controller)
	c.origin = origin
	c.agentId = controllerAgentUri(origin)
	c.state = newControllerState()
	c.ticker = messaging.NewTicker(tickerDur)
	c.poller = messaging.NewTicker(pollerDur)
	//c.revise = messaging.NewTicker(reviseDur)

	c.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.dataC = make(chan *messaging.Message, messaging.ChannelSize)
	c.handler = handler
	return c
}

// String - identity
func (c *controller) String() string { return c.agentId }

// Uri - agent identifier
func (c *controller) Uri() string { return c.agentId }

// Message - message the agent
func (c *controller) Message(m *messaging.Message) { messaging.Mux(m, c.ctrlC, nil, nil) }

// Add - add a shutdown function
func (c *controller) Add(f func()) { c.shutdownFunc = messaging.AddShutdown(c.shutdownFunc, f) }

// Shutdown - shutdown the agent
func (c *controller) Shutdown() {
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
func (c *controller) Run() {
	if c.running {
		return
	}
	go controllerRun(c, controllerFunc, controllerInitFunc, observe, exp, guide)
}

// startup - start tickers
func (c *controller) startup() {
	c.ticker.Start(-1)
	c.poller.Start(-1)
	//c.revise.Start(-1)
}

// shutdown - close resources
func (c *controller) shutdown() {
	close(c.ctrlC)
	c.ticker.Stop()
	c.poller.Stop()
	//c.revise.Stop()
}

func (c *controller) updateTicker(newDuration time.Duration) {
	c.ticker.Start(newDuration)
}

func (c *controller) addEntry(entries []timeseries1.Entry) {
	c.entries = append(c.entries, entries...)
}

type controllerFn func(c *controller, percentile percentile1.Entry, observe *observation, exp *experience) ([]timeseries1.Entry, *core.Status)
type controllerInitFn func(c *controller, observe *observation) *core.Status

// run - ingress controller
func controllerRun(c *controller, ctrlFn controllerFn, initFn controllerInitFn, observe *observation, exp *experience, guide *guidance) {
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
			default:
			}
		case msg := <-c.dataC:
			switch msg.Event() {
			case messaging.DataChangeEvent:
				if msg.IsContentType(core2.ContentTypeProfile) {
					c.handler.AddActivity(c.agentId, "onDataChange() - profile")
					// Process revising the ticker based on the profile.
				} else {
					if msg.IsContentType(core2.ContentTypePercentile) {
						c.handler.AddActivity(c.agentId, "onDataChange() - percentile")
						// Process revising the ticker based on the profile.
					}
				}

			default:
			}
		default:
		}
	}
}
