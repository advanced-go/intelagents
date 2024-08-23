package ingress1

import (
	"errors"
	"fmt"
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

type operativeFunc struct {
	newRedirectAgent       func(f *fieldOperative, state resiliency1.IngressRedirectState)
	processRedirectStartup func(f *fieldOperative, fn *operativeFunc, guide *common.Guidance)
	processRedirectMessage func(f *fieldOperative, fn *operativeFunc, msg *messaging.Message)
}

var operative = func() *operativeFunc {
	return &operativeFunc{
		newRedirectAgent: func(f *fieldOperative, state resiliency1.IngressRedirectState) {
			f.redirect = newRedirectAgent(f.origin, state, f)
		},
		processRedirectStartup: processRedirectStartup,
		processRedirectMessage: processRedirectMessage,
	}
}()

func processRedirectMessage(f *fieldOperative, fn *operativeFunc, msg *messaging.Message) {
	if f.redirect != nil {
		err := errors.New(fmt.Sprintf("error: currently active redirect agent:%v", f.agentId))
		f.handler.Handle(core.NewStatusError(core.StatusInvalidArgument, err), "")
	}
	if r, ok := msg.Body.(resiliency1.RedirectPlan); ok {
		switch r.SQLCommand {
		case "insert":
			f.state.Location = r.Location
			f.state.Status = r.Status
			f.state.EntryId = r.EntryId
			f.state.RouteName = r.RouteName
			if f.state.IsConfigured() {
				fn.newRedirectAgent(f, f.state)
				f.redirect.Run()
			}
		// TODO : how are these handled
		case "update":
		case "delete":
		default:
		}
	} else {
		f.handler.Handle(common.MessageContentTypeErrorStatus(f.agentId, msg), "")
	}
}

func processRedirectStartup(f *fieldOperative, fn *operativeFunc, guide *common.Guidance) {
	if f.redirect != nil {
		err := errors.New(fmt.Sprintf("error: currently active redirect agent:%v", f.agentId))
		f.handler.Handle(core.NewStatusError(core.StatusInvalidArgument, err), "")
	}
	state, status := guide.RedirectState(f.handler, f.origin)
	if !status.OK() {
		return
	}
	f.state = state
	if f.state.IsConfigured() {
		fn.newRedirectAgent(f, f.state)
		f.redirect.Run()
	}
}
