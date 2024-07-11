package guidance1

import (
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	Class           = "ingress-controller1"
	defaultInterval = time.Second * 3
)

type schedule struct {
	running  bool
	uri      string
	origin   core.Origin
	interval time.Duration
	ticker   *time.Ticker
	ctrlC    chan *messaging.Message
	opsAgent messaging.OpsAgent
	shutdown func()
}

func scheduleAgentUri(origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v.%v", Class, origin.Region, origin.Zone, origin.Host)
	}
	return fmt.Sprintf("%v:%v.%v.%v.%v", Class, origin.Region, origin.Zone, origin.SubZone, origin.Host)
}

// NewscheduleAgent - create a new schedule agent
func NewscheduleAgent(origin core.Origin, opsAgent messaging.OpsAgent) messaging.Agent {
	return newscheduleAgent(origin, opsAgent)
}

func newscheduleAgent(origin core.Origin, opsAgent messaging.OpsAgent) *schedule {
	c := new(schedule)
	c.origin = origin
	c.uri = scheduleAgentUri(origin)
	c.interval = defaultInterval
	c.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.opsAgent = opsAgent
	return c
}

// String - identity
func (c *schedule) String() string {
	return c.uri
}

// Uri - agent identifier
func (c *schedule) Uri() string {
	return c.uri
}

// Message - message the agent
func (c *schedule) Message(m *messaging.Message) {
	messaging.Mux(m, c.ctrlC, nil, nil)
}

// Shutdown - shutdown the agent
func (c *schedule) Shutdown() {
	if !c.running {
		return
	}
	c.running = false
	if c.shutdown != nil {
		c.shutdown()
	}
	msg := messaging.NewControlMessage(c.uri, c.uri, messaging.ShutdownEvent)
	if c.ctrlC != nil {
		c.ctrlC <- msg
	}
}

// Run - run the agent
func (c *schedule) Run() {
	if c.running {
		return
	}

}

func (c *schedule) startTicker(interval time.Duration) {
	if interval <= 0 {
		interval = c.interval
	} else {
		c.interval = interval
	}
	if c.ticker != nil {
		c.ticker.Stop()
	}
	c.ticker = time.NewTicker(interval)
}

func (c *schedule) stopTicker() {
	c.ticker.Stop()
}
