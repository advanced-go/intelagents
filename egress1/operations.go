package egress1

import (
	"fmt"
	"github.com/advanced-go/intelagents/dependency1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	OperationsClass = "egress-operations1"
)

// Responsibilities:
//  1. Startup + Restart Events
//     a. Read all egress route configurations
//     b. If authority routing is configured, read all host names
//     c. Create all egress controller agents and a dependency agent if configured
//  2. Changeset Apply Event
//     a. Read new egress and dependency configurations, update controllers and dependency agent
//  3. Changeset Rollback Event
//     b. Read previous egress and dependency configurations, update controllers and dependency agent
//
// 4. Polling - What if an event is missed?? Need some way to save events in database.
type operations struct {
	running  bool
	uri      string
	origin   core.Origin
	version  string        // Current version of origin configuration, helps to stop duplicate updates of egress routes
	interval time.Duration // Needs to be configured dynamically during runtime
	ctrlC    chan *messaging.Message
	handler  messaging.OpsAgent
	manager  messaging.Agent
	shutdown func()
}

func OperationsAgentUri(origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v.%v", OperationsClass, origin.Region, origin.Zone, origin.Host)
	}
	return fmt.Sprintf("%v:%v.%v.%v.%v", OperationsClass, origin.Region, origin.Zone, origin.SubZone, origin.Host)
}

// NewOperationsAgent - create a new operations agent
func NewOperationsAgent(origin core.Origin, handler messaging.OpsAgent) messaging.OpsAgent {
	c := new(operations)
	c.uri = OperationsAgentUri(origin)
	c.origin = origin

	c.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.handler = handler
	c.manager = dependency1.NewDependencyAgent(origin, c)
	return c
}

// String - identity
func (a *operations) String() string {
	return a.uri
}

// Uri - agent identifier
func (a *operations) Uri() string {
	return a.uri
}

// Message - message the agent
func (a *operations) Message(m *messaging.Message) {
	messaging.Mux(m, a.ctrlC, nil, nil)
}

// Handle - error handler
func (a *operations) Handle(status *core.Status, requestId string) *core.Status {
	// TODO : Any operations specific processing ??  If not then forward to handler
	return a.handler.Handle(status, requestId)
}

// Add - add a shutdown function
func (a *operations) Add(f func()) {
	a.shutdown = messaging.AddShutdown(a.shutdown, f)

}

// Shutdown - shutdown the agent
func (a *operations) Shutdown() {
	if !a.running {
		return
	}
	a.running = false
	if a.shutdown != nil {
		a.shutdown()
	}
	msg := messaging.NewControlMessage(a.uri, a.uri, messaging.ShutdownEvent)
	if a.ctrlC != nil {
		a.ctrlC <- msg
	}
}

// Run - run the agent
func (a *operations) Run() {
	if a.running {
		return
	}
	go runOps(a)
}
