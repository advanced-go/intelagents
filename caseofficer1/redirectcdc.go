package caseofficer1

import (
	"fmt"
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

const (
	RedirectCDCClass = "redirect-cdc1"
)

type redirectCDC struct {
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

func redirectCDCUri(origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v", RedirectCDCClass, origin.Region, origin.Zone)
	}
	return fmt.Sprintf("%v:%v.%v.%v", RedirectCDCClass, origin.Region, origin.Zone, origin.SubZone)
}

// redirectCDCAgent - create a new redirect CDC agent
func redirectCDCAgent(origin core.Origin, traffic string, exchange *messaging.Exchange, handler messaging.OpsAgent) messaging.Agent {
	return newRedirectCDC(origin, traffic, exchange, handler)
}

// newRedirectCDC - create a new redirectCDC struct
func newRedirectCDC(origin core.Origin, traffic string, exchange *messaging.Exchange, handler messaging.OpsAgent) *redirectCDC {
	r := new(redirectCDC)
	r.agentId = redirectCDCUri(origin)
	r.origin = origin
	r.traffic = traffic
	//e.ticker = messaging.NewTicker(profile.CaseOfficerDuration(-1))
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

func (r *redirectCDC) startup() {
	r.ticker.Start(-1)
}

func (r *redirectCDC) shutdown() {
	close(r.ctrlC)
	r.ticker.Stop()
}

func runRedirectCDC(r *redirectCDC, guide *guidance) {
	for {
		select {
		case msg := <-r.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				r.shutdown()
				r.handler.AddActivity(r.agentId, messaging.ShutdownEvent)
				return
			case messaging.StartupEvent:
				r.handler.AddActivity(r.agentId, messaging.StartupEvent)
				cdc, status := guide.redirectCDC(r.handler, r.origin)
				if status.OK() {
					for _, e := range cdc {
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

func newRedirectMessage(e resiliency1.CDCRedirect) *messaging.Message {
	// TODO: create valid TO
	origin := core.Origin{Route: e.RouteName}
	to := origin.Route
	msg := messaging.NewControlMessage(to, "", messaging.DataChangeEvent)
	//msg.SetContent(ContentTypeRedirectCDC, nil)
	return msg
}

/*
select {
case <-e.ticker.C():
default:
}

*/
