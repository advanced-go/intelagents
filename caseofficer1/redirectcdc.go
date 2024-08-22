package caseofficer1

import (
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

const (
	RedirectCDCClass = "redirect-cdc1"
)

type redirectCDC struct {
	running      bool
	agentId      string
	lastId       int
	origin       core.Origin
	ctrlC        chan *messaging.Message
	handler      messaging.OpsAgent
	exchange     *messaging.Exchange
	shutdownFunc func()
}

func redirectCDCUri(origin core.Origin) string {
	return origin.Uri(RedirectCDCClass)
}

// redirectCDCAgent - create a new redirect CDC agent
func redirectCDCAgent(origin core.Origin, lastId int, exchange *messaging.Exchange, handler messaging.OpsAgent) messaging.Agent {
	return newRedirectCDC(origin, lastId, exchange, handler)
}

// newRedirectCDC - create a new redirectCDC struct
func newRedirectCDC(origin core.Origin, lastId int, exchange *messaging.Exchange, handler messaging.OpsAgent) *redirectCDC {
	r := new(redirectCDC)
	r.agentId = redirectCDCUri(origin)
	r.lastId = lastId
	r.origin = origin
	r.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	r.handler = handler
	r.exchange = exchange
	return r
}

// String - identity
func (r *redirectCDC) String() string { return r.Uri() }

// Uri - agent identifier
func (r *redirectCDC) Uri() string { return r.agentId }

// Message - message the agent
func (r *redirectCDC) Message(m *messaging.Message) { messaging.Mux(m, r.ctrlC, nil, nil) }

// Add - add a shutdown function
func (r *redirectCDC) Add(f func()) { r.shutdownFunc = messaging.AddShutdown(r.shutdownFunc, f) }

// Run - run the agent
func (r *redirectCDC) Run() {
	if r.running {
		return
	}
	r.running = true
	go runRedirectCDC(r, guide)
}

// Shutdown - shutdown the agent
func (r *redirectCDC) Shutdown() {
	if !r.running {
		return
	}
	r.running = false
	// Is this needed or called in the right place??
	if r.shutdownFunc != nil {
		r.shutdownFunc()
	}
	msg := messaging.NewControlMessage(r.agentId, r.agentId, messaging.ShutdownEvent)
	if r.ctrlC != nil {
		r.ctrlC <- msg
	}
}

func (r *redirectCDC) shutdown() { close(r.ctrlC) }

func runRedirectCDC(r *redirectCDC, guide *guidance) {
	for {
		select {
		case msg := <-r.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				r.shutdown()
				r.handler.AddActivity(r.agentId, messaging.ShutdownEvent)
				return
			case messaging.ProcessEvent:
				r.handler.AddActivity(r.agentId, messaging.ProcessEvent)
				entry, status := guide.updatedRedirectPlans(r.handler, r.origin, r.lastId)
				if status.OK() {
					r.lastId = entry[len(entry)-1].EntryId
					for _, e := range entry {
						err := r.exchange.Send(newRedirectMessage(e))
						if err != nil {
							r.handler.Handle(core.NewStatusError(core.StatusInvalidArgument, err), "")
						}
					}
				}
			default:
			}
		default:
		}
	}
}

func newRedirectMessage(e resiliency1.RedirectPlan) *messaging.Message {
	// TODO: create valid TO
	origin := core.Origin{Route: e.RouteName}
	to := redirectCDCUri(origin)
	msg := messaging.NewControlMessage(to, "", messaging.DataChangeEvent)
	msg.SetContentType(common.ContentTypeRedirectPlan)
	return msg
}

/*
select {
case <-e.ticker.C():
default:
}

*/