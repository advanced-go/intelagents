package caseofficer1

import (
	"fmt"
	"github.com/advanced-go/stdlib/access"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

const (
	EntryCDCClass = "entry-cdc1"
)

type entryCDC struct {
	running      bool
	agentId      string
	traffic      string
	origin       core.Origin
	ticker       *messaging.Ticker
	ctrlC        chan *messaging.Message
	handler      messaging.OpsAgent
	exchange     *messaging.Exchange
	shutdownFunc func()
}

func entryCDCUri(origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v", EntryCDCClass, origin.Region, origin.Zone)
	}
	return fmt.Sprintf("%v:%v.%v.%v", EntryCDCClass, origin.Region, origin.Zone, origin.SubZone)
}

// entryCDCAgent - create a new entry CDC agent
func entryCDCAgent(origin core.Origin, traffic string, exchange *messaging.Exchange, handler messaging.OpsAgent) messaging.Agent {
	return newEntryCDC(origin, traffic, exchange, handler)
}

// newEntryCDC - create a new entryCDC struct
func newEntryCDC(origin core.Origin, traffic string, exchange *messaging.Exchange, handler messaging.OpsAgent) *entryCDC {
	e := new(entryCDC)
	e.agentId = entryCDCUri(origin)
	e.origin = origin
	e.traffic = traffic
	//e.ticker = messaging.NewTicker(profile.CaseOfficerDuration(-1))
	e.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	e.handler = handler
	e.exchange = exchange
	return e
}

// String - identity
func (e *entryCDC) String() string { return e.Uri() }

// Uri - agent identifier
func (e *entryCDC) Uri() string { return e.agentId }

// Message - message the agent
func (e *entryCDC) Message(m *messaging.Message) { messaging.Mux(m, e.ctrlC, nil, nil) }

// Add - add a shutdown function
func (e *entryCDC) Add(f func()) { e.shutdownFunc = messaging.AddShutdown(e.shutdownFunc, f) }

// Run - run the agent
func (e *entryCDC) Run() {
	if e.running {
		return
	}
	e.running = true
	go runEntryCDC(e, officer, guide)
}

// Shutdown - shutdown the agent
func (e *entryCDC) Shutdown() {
	if !e.running {
		return
	}
	e.running = false
	// Is this needed or called in the right place??
	if e.shutdownFunc != nil {
		e.shutdownFunc()
	}
	msg := messaging.NewControlMessage(e.agentId, e.agentId, messaging.ShutdownEvent)
	if e.ctrlC != nil {
		e.ctrlC <- msg
	}

}

func (e *entryCDC) startup() {
	e.ticker.Start(-1)
}

func (e *entryCDC) shutdown() {
	close(e.ctrlC)
	e.ticker.Stop()
}

func (e *entryCDC) isIngress() bool {
	return e.traffic == access.IngressTraffic
}

func runEntryCDC(e *entryCDC, fn *caseOfficerFunc, guide *guidance) {
	for {
		select {
		case msg := <-e.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				e.shutdown()
				e.handler.AddActivity(e.agentId, messaging.ShutdownEvent)
				return
			case messaging.StartupEvent:
				e.handler.AddActivity(e.agentId, messaging.StartupEvent)
			default:
			}
		default:
		}
	}
}

/*
select {
case <-e.ticker.C():
default:
}

*/
