package caseofficer1

import (
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

const (
	EntryCDCClass = "entry-cdc1"
)

type entryCDC struct {
	running       bool
	agentId       string
	lastId        int
	origin        core.Origin
	ctrlC         chan *messaging.Message
	handler       messaging.OpsAgent
	ingressAgents *messaging.Exchange
	egressAgents  *messaging.Exchange
	shutdownFunc  func()
}

func entryCDCUri(origin core.Origin) string {
	return origin.Uri(EntryCDCClass)
}

// entryCDCAgent - create a new redirect CDC agent
func entryCDCAgent(origin core.Origin, lastId int, ingressAgents, egressAgents *messaging.Exchange, handler messaging.OpsAgent) messaging.Agent {
	return newEntryCDC(origin, lastId, ingressAgents, egressAgents, handler)
}

// newEntryCDC - create a new entryCDC struct
func newEntryCDC(origin core.Origin, lastId int, ingressAgents, egressAgents *messaging.Exchange, handler messaging.OpsAgent) *entryCDC {
	e := new(entryCDC)
	e.agentId = entryCDCUri(origin)
	e.lastId = lastId
	e.origin = origin
	e.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	e.handler = handler
	e.ingressAgents = ingressAgents
	e.egressAgents = egressAgents
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
	go runEntryCDC(e, guide)
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

func (e *entryCDC) shutdown() { close(e.ctrlC) }

func runEntryCDC(e *entryCDC, guide *guidance) {
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
				entry, status := guide.newAssignments(e.handler, e.origin, e.lastId)
				if status.OK() {
					e.lastId = entry[len(entry)-1].EntryId
					//updateExchange()
					for _, e1 := range entry {
						// TODO : update the ingress and egress agent exchanges
						if e1.EntryId != 0 {
						}
					}
				}
			default:
			}
		default:
		}
	}
}
