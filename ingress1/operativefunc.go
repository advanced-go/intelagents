package ingress1

import (
	"errors"
	"fmt"
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/stdlib/core"
)

type operativeFunc struct {
	setRedirectState func(f *fieldOperative, guide *guidance) *core.Status
	newRedirectAgent func(f *fieldOperative, state *resiliency1.IngressRedirectState)
	processRedirect  func(f *fieldOperative, fn *operativeFunc, guide *guidance)
}

var operative = func() *operativeFunc {
	return &operativeFunc{
		setRedirectState: func(f *fieldOperative, guide *guidance) *core.Status {
			s, status := guide.redirectState(f.handler, f.origin)
			if status.OK() {
				*f.state = *s
			}
			return status
		},
		newRedirectAgent: func(f *fieldOperative, state *resiliency1.IngressRedirectState) {
			f.redirect = newRedirectAgent(f.origin, state, f)
		},
		processRedirect: processRedirect,
	}
}()

func processRedirect(f *fieldOperative, fn *operativeFunc, guide *guidance) {
	if f.redirect != nil {
		err := errors.New(fmt.Sprintf("error: currently active redirect agent:%v", f.agentId))
		f.handler.Handle(core.NewStatusError(core.StatusInvalidArgument, err), "")
	}
	fn.setRedirectState(f, guide)
	if f.state.IsConfigured() {
		fn.newRedirectAgent(f, f.state)
		f.redirect.Run()
	}
}
