package ingress1

import (
	"errors"
	"fmt"
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"reflect"
	"time"
)

const (
	RedirectProcessClass  = "ingress-redirect-process1"
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
	return origin.Uri(RedirectProcessClass)
}

// newRedirectProcessAgent - create a new lead agent
func newRedirectProcessAgent(origin core.Origin, handler messaging.OpsAgent) messaging.Agent {
	return newRedirectProcess(origin, handler, processTickerDuration)
}

func newRedirectProcess(origin core.Origin, handler messaging.OpsAgent, tickerDur time.Duration) *redirectProcess {
	r := new(redirectProcess)
	r.agentId = redirectProcessAgentUri(origin)
	r.origin = origin
	r.state = resiliency1.NewIngressRedirectState()
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
			case messaging.ProcessEvent:
				r.handler.AddActivity(r.agentId, fmt.Sprintf("%v - %v", msg.Event(), msg.ContentType()))
				processProcessEvent(r, msg)
			default:
				r.handler.Handle(common.MessageEventErrorStatus(r.agentId, msg), "")
			}
		default:
		}
	}
}

func processProcessEvent(r *redirectProcess, msg *messaging.Message) {
	switch msg.ContentType() {
	case common.ContentTypeRedirectState:
		// If current state is active, then error
		// TODO : reconcile with redirect state??
		if r.state.IsActive() {
			err := errors.New(fmt.Sprintf("error: currently active state for agent:%v", r.agentId))
			r.handler.Handle(core.NewStatusError(core.StatusExecError, err), "")
		}
		if rs, ok := msg.Body.(*resiliency1.IngressRedirectState); ok {
			r.state = rs
			return
		}
		err := errors.New(fmt.Sprintf("error: redirect state process type:%v is invalid for agent:%v", reflect.TypeOf(msg.Body), r.agentId))
		r.handler.Handle(core.NewStatusError(core.StatusInvalidArgument, err), "")
		r.state = nil
	default:
		r.handler.Handle(common.MessageContentTypeErrorStatus(r.agentId, msg), "")
	}
}
