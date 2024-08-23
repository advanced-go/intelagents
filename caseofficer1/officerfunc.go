package caseofficer1

import (
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/intelagents/egress1"
	"github.com/advanced-go/intelagents/ingress1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

// A nod to Linus Torvalds and plain C
type caseOfficerFunc struct {
	startup func(c *caseOfficer, guide *common.Guidance) *core.Status
	update  func(c *caseOfficer, guide *common.Guidance) *core.Status
}

var (
	newFieldOperative = func(ingress bool, profile *common.Profile, origin core.Origin, handler messaging.OpsAgent) messaging.Agent {
		if ingress {
			return ingress1.NewFieldOperative(origin, profile, handler)
		}
		return egress1.NewFieldOperative(origin, profile, handler)
	}

	officer = func() *caseOfficerFunc {
		return &caseOfficerFunc{
			startup: func(c *caseOfficer, guide *common.Guidance) *core.Status {
				entry, lastId, status := guide.Assignments(c.handler, c.origin)
				if status.OK() {
					c.lastId = lastId
					updateExchange(c, entry)
				}
				c.redirectAgent = newRedirectCDC(c.origin, c.lastId.Redirect, c.ingressAgents, c)
				c.failoverAgent = newEgressCDC(c.origin, c.lastId.Egress, c.egressAgents, c)
				c.startup()
				return core.StatusOK()
			},
			update: func(c *caseOfficer, guide *common.Guidance) *core.Status {
				entry, status := guide.NewAssignments(c.handler, c.origin, c.lastId.Entry)
				if status.OK() {
					c.lastId.Entry = entry[len(entry)-1].EntryId
					updateExchange(c, entry)
				}
				return core.StatusOK()
			},
		}
	}()
)

func updateExchange(c *caseOfficer, entries []resiliency1.HostEntry) {
	for _, e := range entries {
		o := core.Origin{
			Region:     e.Region,
			Zone:       e.Zone,
			SubZone:    "",
			Host:       "",
			InstanceId: "",
			Route:      "",
		}
		initAgent(c, o, true)
		initAgent(c, o, false)
	}
}

func initAgent(c *caseOfficer, origin core.Origin, ingress bool) {
	var agent messaging.Agent
	var err error

	if ingress {
		agent = newFieldOperative(ingress, c.profile, origin, c)
		err = c.ingressAgents.Register(agent)
	} else {
		agent = newFieldOperative(ingress, c.profile, origin, c)
		err = c.egressAgents.Register(agent)
	}
	if err != nil {
		c.handler.Handle(core.NewStatusError(core.StatusInvalidArgument, err), "")
	} else {
		agent.Run()
	}
}
