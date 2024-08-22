package caseofficer1

import (
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

const (
	FailoverCDCClass = "failover-cdc1"
)

type failoverCDC struct {
	running      bool
	agentId      string
	lastId       int
	origin       core.Origin
	ctrlC        chan *messaging.Message
	handler      messaging.OpsAgent
	exchange     *messaging.Exchange
	shutdownFunc func()
}

func failoverCDCUri(origin core.Origin) string {
	return origin.Uri(FailoverCDCClass)
}

// failoverCDCAgent - create a new redirect CDC agent
func failoverCDCAgent(origin core.Origin, lastId int, exchange *messaging.Exchange, handler messaging.OpsAgent) messaging.Agent {
	return newFailoverCDC(origin, lastId, exchange, handler)
}

// newFailoverCDC - create a new failoverCDCC struct
func newFailoverCDC(origin core.Origin, lastId int, exchange *messaging.Exchange, handler messaging.OpsAgent) *failoverCDC {
	r := new(failoverCDC)
	r.agentId = failoverCDCUri(origin)
	r.lastId = lastId
	r.origin = origin
	r.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	r.handler = handler
	r.exchange = exchange
	return r
}

// String - identity
func (f *failoverCDC) String() string { return f.Uri() }

// Uri - agent identifier
func (f *failoverCDC) Uri() string { return f.agentId }

// Message - message the agent
func (f *failoverCDC) Message(m *messaging.Message) { messaging.Mux(m, f.ctrlC, nil, nil) }

// Add - add a shutdown function
func (f *failoverCDC) Add(fn func()) { f.shutdownFunc = messaging.AddShutdown(f.shutdownFunc, fn) }

// Run - run the agent
func (f *failoverCDC) Run() {
	if f.running {
		return
	}
	f.running = true
	go runFailoverCDC(f, guide)
}

// Shutdown - shutdown the agent
func (f *failoverCDC) Shutdown() {
	if !f.running {
		return
	}
	f.running = false
	// Is this needed or called in the right place??
	if f.shutdownFunc != nil {
		f.shutdownFunc()
	}
	msg := messaging.NewControlMessage(f.agentId, f.agentId, messaging.ShutdownEvent)
	if f.ctrlC != nil {
		f.ctrlC <- msg
	}
}

func (f *failoverCDC) shutdown() { close(f.ctrlC) }

func runFailoverCDC(f *failoverCDC, guide *guidance) {
	for {
		select {
		case msg := <-f.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				f.shutdown()
				f.handler.AddActivity(f.agentId, messaging.ShutdownEvent)
				return
			case messaging.ProcessEvent:
				f.handler.AddActivity(f.agentId, messaging.ProcessEvent)
				entry, status := guide.updatedFailoverPlans(f.handler, f.origin, f.lastId)
				if status.OK() {
					f.lastId = entry[len(entry)-1].EntryId
					for _, e := range entry {
						err := f.exchange.Send(newFailoverMessage(e))
						if err != nil {
							f.handler.Handle(core.NewStatusError(core.StatusInvalidArgument, err), "")
						}
					}
				}
			default:
			}
		default:
		}
	}
}

func newFailoverMessage(e resiliency1.FailoverPlan) *messaging.Message {
	// TODO: create valid TO
	origin := core.Origin{Route: e.RouteName}
	to := failoverCDCUri(origin)
	msg := messaging.NewControlMessage(to, "", messaging.DataChangeEvent)
	msg.SetContentType(common.ContentTypeFailoverPlan)
	return msg
}
