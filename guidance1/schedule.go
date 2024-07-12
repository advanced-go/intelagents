package guidance1

import (
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	ScheduleClass = "schedule1"
)

type schedule struct {
	running  bool
	uri      string
	ticker   *messaging.Ticker
	ctrlC    chan *messaging.Message
	handler  messaging.OpsAgent
	shutdown func()
}

func ScheduleAgentUri(origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v.%v", ScheduleClass, origin.Region, origin.Zone, origin.Host)
	}
	return fmt.Sprintf("%v:%v.%v.%v.%v", ScheduleClass, origin.Region, origin.Zone, origin.SubZone, origin.Host)
}

// NewScheduleAgent - create a new schedule agent
func NewScheduleAgent(interval time.Duration, handler messaging.OpsAgent) messaging.Agent {
	return newScheduleAgent(interval, handler)
}

func newScheduleAgent(interval time.Duration, handler messaging.OpsAgent) *schedule {
	c := new(schedule)
	c.uri = ScheduleClass
	c.ticker = messaging.NewTicker(interval)
	c.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.handler = handler
	return c
}

// String - identity
func (s *schedule) String() string {
	return s.uri
}

// Uri - agent identifier
func (s *schedule) Uri() string {
	return s.uri
}

// Message - message the agent
func (s *schedule) Message(m *messaging.Message) {
	messaging.Mux(m, s.ctrlC, nil, nil)
}

// Shutdown - shutdown the agent
func (s *schedule) Shutdown() {
	if !s.running {
		return
	}
	s.running = false
	if s.shutdown != nil {
		s.shutdown()
	}
	msg := messaging.NewControlMessage(s.uri, s.uri, messaging.ShutdownEvent)
	if s.ctrlC != nil {
		s.ctrlC <- msg
	}
}

// Run - run the agent
func (s *schedule) Run() {
	if s.running {
		return
	}
	go runSchedule(s)
}
