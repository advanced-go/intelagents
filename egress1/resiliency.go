package egress1

import (
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	ResiliencyClass = "egress-resiliency1"
)

type resiliency struct {
	running bool
	agentId string
	origin  core.Origin

	interval time.Duration // Needs to be configured dynamically during runtime
	ctrlC    chan *messaging.Message
	handler  messaging.OpsAgent
	shutdown func()
}

func ResiliencyAgentUri(origin core.Origin) string {
	return origin.Uri(ResiliencyClass)
}

// NewResiliencyAgent - create a new Resiliency agent
func NewResiliencyAgent(origin core.Origin, handler messaging.OpsAgent) messaging.OpsAgent {
	c := new(resiliency)
	c.agentId = ResiliencyAgentUri(origin)
	c.origin = origin

	c.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.handler = handler
	return c
}

// String - identity
func (a *resiliency) String() string { return a.Uri() }

// Uri - agent identifier
func (a *resiliency) Uri() string { return a.agentId }

// Message - message the agent
func (a *resiliency) Message(m *messaging.Message) {
	messaging.Mux(m, a.ctrlC, nil, nil)
}

// Handle - error handler
func (a *resiliency) Handle(status *core.Status, requestId string) *core.Status {
	// TODO : Any resiliency specific processint ??  If not then forward to handler
	return a.handler.Handle(status, requestId)
}

// AddActivity - add activity
func (a *resiliency) AddActivity(agentId string, content any) {
	// TODO : Any operations specific processing ??  If not then forward to handler
	//return a.handler.Handle(status, requestId)
}

// Add - add a shutdown function
func (a *resiliency) Add(f func()) {
	a.shutdown = messaging.AddShutdown(a.shutdown, f)

}

// Shutdown - shutdown the agent
func (a *resiliency) Shutdown() {
	if !a.running {
		return
	}
	a.running = false
	if a.shutdown != nil {
		a.shutdown()
	}
	msg := messaging.NewControlMessage(a.agentId, a.agentId, messaging.ShutdownEvent)
	if a.ctrlC != nil {
		a.ctrlC <- msg
	}
}

// Run - run the agent
func (a *resiliency) Run() {
	if a.running {
		return
	}
	//go runController(a, access1.EgressQuery, inference1.EgressQuery, nil, nil)
}
