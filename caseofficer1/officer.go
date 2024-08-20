package caseofficer1

import (
	"fmt"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	CaseOfficerClass = "case-officer1"
)

type caseOfficer struct {
	running        bool
	agentId        string
	origin         core.Origin
	traffic        string
	lastEntryId    int
	lastRedirectId int
	lastFailoverId int
	profile        *common.Profile
	ticker         *messaging.Ticker
	ctrlC          chan *messaging.Message
	handler        messaging.OpsAgent
	ingressAgents  *messaging.Exchange
	egressAgents   *messaging.Exchange
	shutdownFunc   func()
}

func AgentUri(traffic string, origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf(core.RegionZoneHostFmt, CaseOfficerClass, traffic, origin.Region, origin.Zone)
	}
	return fmt.Sprintf(core.RegionZoneSubZoneHostFmt, CaseOfficerClass, traffic, origin.Region, origin.Zone, origin.SubZone)
}

// NewAgent - create a new case officer agent
func NewAgent(traffic string, origin core.Origin, profile *common.Profile, handler messaging.OpsAgent) messaging.OpsAgent {
	return newAgent(traffic, origin, profile, handler)
}

// newAgent - create a new case officer agent
func newAgent(traffic string, origin core.Origin, profile *common.Profile, handler messaging.OpsAgent) *caseOfficer {
	c := new(caseOfficer)
	c.agentId = AgentUri(traffic, origin)
	c.traffic = traffic
	c.origin = origin
	c.profile = profile
	c.ticker = messaging.NewTicker(profile.CaseOfficerDuration(-1))

	c.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.handler = handler
	c.ingressAgents = messaging.NewExchange()
	c.egressAgents = messaging.NewExchange()
	return c
}

// String - identity
func (c *caseOfficer) String() string { return c.Uri() }

// Uri - agent identifier
func (c *caseOfficer) Uri() string { return c.agentId }

// Message - message the agent
func (c *caseOfficer) Message(m *messaging.Message) { messaging.Mux(m, c.ctrlC, nil, nil) }

// Handle - error handler
func (c *caseOfficer) Handle(status *core.Status, requestId string) *core.Status {
	// TODO : do we need any processing specific to a case officer? If not then forward to handler
	return c.handler.Handle(status, requestId)
}

// AddActivity - add activity
func (c *caseOfficer) AddActivity(agentId string, content any) {
	// TODO : Any operations specific processing ??  If not then forward to handler
	c.handler.AddActivity(agentId, content)
}

// Add - add a shutdown function
func (c *caseOfficer) Add(f func()) { c.shutdownFunc = messaging.AddShutdown(c.shutdownFunc, f) }

// Run - run the agent
func (c *caseOfficer) Run() {
	if c.running {
		return
	}
	c.running = true
	go runCaseOfficer(c, officer, guide)
}

// Shutdown - shutdown the agent
func (c *caseOfficer) Shutdown() {
	if !c.running {
		return
	}
	c.running = false
	// Is this needed or called in the right place??
	if c.shutdownFunc != nil {
		c.shutdownFunc()
	}
	msg := messaging.NewControlMessage(c.agentId, c.agentId, messaging.ShutdownEvent)
	if c.ctrlC != nil {
		c.ctrlC <- msg
	}
	c.ingressAgents.Broadcast(msg)
	c.egressAgents.Broadcast(msg)
}

func (c *caseOfficer) startup() {
	c.ticker.Start(-1)
}

func (c *caseOfficer) shutdown() {
	close(c.ctrlC)
	c.ticker.Stop()
}

func (c *caseOfficer) reviseTicker(newDuration time.Duration) {
	c.ticker.Start(newDuration)
}

func runCaseOfficer(c *caseOfficer, fn *caseOfficerFunc, guide *guidance) {
	//percentile, _ := fn.startup(c, nil, guide)

	for {
		// main processing : query for and add new assignments
		select {
		case <-c.ticker.C():
			c.handler.AddActivity(c.agentId, "onTick")
			//fn.process(r, percentile, observe, exp)
		default:
		}
		// control channel processing
		select {
		case msg := <-c.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				c.shutdown()
				c.handler.AddActivity(c.agentId, messaging.ShutdownEvent)
				return
			case messaging.DataChangeEvent:
				if msg.IsContentType(common.ContentTypeProfile) {
					// TODO : broadcast to all field operatives
				}
			default:
			}
		default:
		}
	}
}
