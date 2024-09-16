package caseofficer1

import (
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/intelagents/egress1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

const (
	EgressCDCClass = "egress-cdc1"
)

type egressCDC struct {
	running      bool
	agentId      string
	lastId       int
	origin       core.Origin
	ctrlC        chan *messaging.Message
	handler      messaging.OpsAgent
	exchange     *messaging.Exchange
	shutdownFunc func()
}

func egressCDCUri(origin core.Origin) string {
	return origin.Uri(EgressCDCClass)
}

// egressCDCAgent - create a new egress configuration CDC agent
func egressCDCAgent(origin core.Origin, lastId int, exchange *messaging.Exchange, handler messaging.OpsAgent) messaging.Agent {
	return newEgressCDC(origin, lastId, exchange, handler)
}

// newEgressCDC - create a new egressCDC struct
func newEgressCDC(origin core.Origin, lastId int, exchange *messaging.Exchange, handler messaging.OpsAgent) *egressCDC {
	e := new(egressCDC)
	e.agentId = egressCDCUri(origin)
	e.lastId = lastId
	e.origin = origin
	e.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	e.handler = handler
	e.exchange = exchange
	return e
}

// String - identity
func (e *egressCDC) String() string { return e.Uri() }

// Uri - agent identifier
func (e *egressCDC) Uri() string { return e.agentId }

// Message - message the agent
func (e *egressCDC) Message(m *messaging.Message) { e.ctrlC <- m }

// Add - add a shutdown function
func (e *egressCDC) Add(fn func()) { e.shutdownFunc = messaging.AddShutdown(e.shutdownFunc, fn) }

// Run - run the agent
func (e *egressCDC) Run() {
	if e.running {
		return
	}
	e.running = true
	go runEgressCDC(e, common.Guide)
}

// Shutdown - shutdown the agent
func (e *egressCDC) Shutdown() {
	if !e.running {
		return
	}
	e.running = false
	if e.shutdownFunc != nil {
		e.shutdownFunc()
	}
	msg := messaging.NewControlMessage(e.agentId, e.agentId, messaging.ShutdownEvent)
	if e.ctrlC != nil {
		e.ctrlC <- msg
	}
}

func (e *egressCDC) shutdown() { close(e.ctrlC) }

func runEgressCDC(e *egressCDC, guide *common.Guidance) {
	for {
		select {
		case msg := <-e.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				e.shutdown()
				e.handler.AddActivity(e.agentId, messaging.ShutdownEvent)
				return
			case messaging.ProcessEvent:
				e.handler.AddActivity(e.agentId, messaging.ProcessEvent)
				entries, status := guide.UpdatedEgressConfigs(e.handler, e.origin, e.lastId)
				if status.OK() {
					e.lastId = entries[len(entries)-1].EntryId
					for _, c := range entries {
						err := e.exchange.Send(newEgressMessage(e, c))
						if err != nil {
							e.handler.Handle(core.NewStatusError(core.StatusInvalidArgument, err))
						}
					}
				}
			default:
			}
		default:
		}
	}
}

func newEgressMessage(e *egressCDC, c resiliency1.EgressConfig) *messaging.Message {
	o := c.Origin()
	to := egress1.FieldOperativeUri(o)
	msg := messaging.NewControlMessage(to, e.agentId, messaging.DataChangeEvent)
	msg.SetContentType(common.ContentTypeEgressConfig)
	msg.Body = e
	return msg
}
