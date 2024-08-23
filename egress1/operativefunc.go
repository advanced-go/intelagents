package egress1

import (
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

type operativeFunc struct {
	startup      func(f *fieldOperative, guide *common.Guidance)
	onDataChange func(f *fieldOperative, guide *common.Guidance, msg *messaging.Message)
}

var (
	newAgent = func(origin core.Origin, profile *common.Profile, state resiliency1.EgressState, handler messaging.OpsAgent) messaging.Agent {
		return newResiliencyAgent(origin, profile, state, handler)
	}

	operative = func() *operativeFunc {
		return &operativeFunc{
			startup: func(f *fieldOperative, guide *common.Guidance) {
				s, status := guide.EgressState(f.handler, f.origin)
				if !status.OK() {
					return
				}
				updateExchange(f, s)
				return
			},
			onDataChange: func(f *fieldOperative, guide *common.Guidance, msg *messaging.Message) {
				plan, ok := msg.Body.(resiliency1.EgressConfig)
				if !ok {
					f.handler.Handle(common.MessageContentTypeErrorStatus(f.agentId, msg), "")
					return
				}
				switch plan.SQLCommand {
				case common.SQLUpdate:
					msg.SetFrom(f.agentId)
					msg.SetTo(plan.Origin().Uri(ResiliencyClass))
					f.agents.Send(msg)
				case common.SQLInsert:
					var state resiliency1.EgressState
					resiliency1.NewEgressState(&state)
					// Stale profile is OK, as the resiliency agent can handle the RPS mismatch between the
					// profile and actual.
					agent := newAgent(plan.Origin(), f.profile, state, f.handler)
					f.agents.Register(agent)
					agent.Run()
				case common.SQLDelete:
					a := f.agents.Get(plan.Origin().Uri(ResiliencyClass))
					if a != nil {
						a.Shutdown()
					}
				default:
				}
			},
		}
	}()
)

func updateExchange(f *fieldOperative, entries []resiliency1.EgressState) {
	for i, _ := range entries {
		e := entries[i]
		agent := newAgent(e.Origin(), f.profile, e, f.handler)
		f.agents.Register(agent)
		agent.Run()
	}
}
