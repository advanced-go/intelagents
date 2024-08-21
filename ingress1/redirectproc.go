package ingress1

import (
	"fmt"
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	RedirectProcessClass  = "ingress-redirect1"
	processTickerDuration = time.Second * 60
)

type redirectProcess struct {
	running      bool
	agentId      string
	origin       core.Origin
	state        *resiliency1.IngressRedirectState
	ticker       *messaging.Ticker
	ctrlC        chan *messaging.Message
	handler      messaging.OpsAgent
	shutdownFunc func()
}

func redirectProcessAgentUri(origin core.Origin) string {
	if origin.SubZone == "" {
		return fmt.Sprintf("%v:%v.%v.%v", RedirectProcessClass, origin.Region, origin.Zone, origin.Host)
	}
	return fmt.Sprintf("%v:%v.%v.%v.%v", RedirectProcessClass, origin.Region, origin.Zone, origin.SubZone, origin.Host)
}

// newRedirectProcessAgent - create a new lead agent
func newRedirectProcessAgent(origin core.Origin, state *resiliency1.IngressRedirectState, handler messaging.OpsAgent) messaging.Agent {
	return newRedirectProcess(origin, state, handler, processTickerDuration)
}

func newRedirectProcess(origin core.Origin, state *resiliency1.IngressRedirectState, handler messaging.OpsAgent, tickerDur time.Duration) *redirectProcess {
	r := new(redirectProcess)
	r.agentId = redirectProcessAgentUri(origin)
	r.origin = origin
	r.state = state
	r.ticker = messaging.NewTicker(tickerDur)
	r.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	r.handler = handler
	return r
}

// String - identity
func (r *redirectProcess) String() string { return r.agentId }

// Uri - agent identifier
func (r *redirectProcess) Uri() string { return r.agentId }

// Message - message the agent
func (r *redirectProcess) Message(m *messaging.Message) { messaging.Mux(m, r.ctrlC, nil, nil) }

// Add - add a shutdown function
func (r *redirectProcess) Add(f func()) { r.shutdownFunc = messaging.AddShutdown(r.shutdownFunc, f) }

// Shutdown - shutdown the agent
func (r *redirectProcess) Shutdown() {
	if !r.running {
		return
	}
	r.running = false
	if r.shutdownFunc != nil {
		r.shutdownFunc()
	}
	msg := messaging.NewControlMessage(r.agentId, r.agentId, messaging.ShutdownEvent)
	if r.ctrlC != nil {
		r.ctrlC <- msg
	}
}

// Run - run the agent
func (r *redirectProcess) Run() {
	if r.running {
		return
	}
	go runRedirectProcess(r, redirection, common.Observe, localGuidance)
}

// startup - start tickers
func (r *redirectProcess) startup() {
	r.ticker.Start(-1)
}

// shutdown - close resources
func (r *redirectProcess) shutdown() {
	close(r.ctrlC)
	r.ticker.Stop()
}

func (r *redirect) updatePercentage() {
	switch r.state.Percentage {
	case 0:
		r.state.Percentage = 10
	case 10:
		r.state.Percentage = 20
	case 20:
		r.state.Percentage = 40
	case 40:
		r.state.Percentage = 70
	case 70:
		r.state.Percentage = 100
	default:
	}
}

func runRedirectProcess(r *redirectProcess, fn *redirectFunc, observe *common.Observation, guide *guidance) {
	//fn.startup(r, guide)

	for {
		select {
		case <-r.ticker.C():

		// control channel
		case msg := <-r.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				r.shutdown()
				r.handler.AddActivity(r.agentId, messaging.ShutdownEvent)
				return
			case messaging.DataChangeEvent:
				if msg.IsContentType(common.ContentTypeRedirectPlan) {
					r.handler.AddActivity(r.agentId, "onDataChange() - redirect plan")
					//r.updateRedirectPlan(guide)
					//r.updatePercentileSLO(guide)
				}
			default:
			}
		default:
		}
	}
}
