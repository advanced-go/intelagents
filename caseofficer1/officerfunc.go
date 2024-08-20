package caseofficer1

import (
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/intelagents/egress1"
	"github.com/advanced-go/intelagents/ingress1"
	"github.com/advanced-go/stdlib/access"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

// A nod to Linus Torvalds and plain C
type caseOfficerFunc struct {
	startup           func(c *caseOfficer, guide *guidance) *core.Status
	update            func(c *caseOfficer, guide *guidance) *core.Status
	newFieldOperative func(traffic string, profile *common.Profile, origin core.Origin, handler messaging.OpsAgent) messaging.Agent
}

var officer = func() *caseOfficerFunc {
	return &caseOfficerFunc{
		startup: func(c *caseOfficer, guide *guidance) *core.Status {
			entry, lastId, status := guide.assignments(c.handler, c.origin)
			if status.OK() {
				c.lastId = lastId
				updateExchange(c, entry)
			}
			c.redirectAgent = newRedirectCDC(c.origin, c.lastId.Redirect, c.ingressAgents, c)
			c.failoverAgent = newFailoverCDC(c.origin, c.lastId.Failover, c.egressAgents, c)
			c.startup()
			return core.StatusOK()
		},
		update: func(c *caseOfficer, guide *guidance) *core.Status {
			entry, status := guide.newAssignments(c.handler, c.origin, c.lastId.Entry)
			if status.OK() {
				c.lastId.Entry = entry[len(entry)-1].EntryId
				updateExchange(c, entry)
			}
			return core.StatusOK()
		},
		newFieldOperative: func(traffic string, profile *common.Profile, origin core.Origin, handler messaging.OpsAgent) messaging.Agent {
			if traffic == access.IngressTraffic {
				return ingress1.NewFieldOperative(origin, profile, handler)
			}
			return egress1.NewFieldOperative(origin, profile, handler)
		},
	}
}()

func updateExchange(c *caseOfficer, entries []resiliency1.Entry) {
	for _, e := range entries {
		o := core.Origin{
			Region:     e.Region,
			Zone:       e.Zone,
			SubZone:    "",
			Host:       "",
			InstanceId: "",
			Route:      "",
		}
		initAgent(c, access.IngressTraffic, o)
		initAgent(c, access.EgressTraffic, o)
	}
}

func initAgent(c *caseOfficer, traffic string, origin core.Origin) {
	var agent messaging.Agent
	var err error

	if traffic == access.IngressTraffic {
		agent = officer.newFieldOperative(access.IngressTraffic, c.profile, origin, c)
		err = c.ingressAgents.Register(agent)
	} else {
		agent = officer.newFieldOperative(access.EgressTraffic, c.profile, origin, c)
		err = c.egressAgents.Register(agent)
	}
	if err != nil {
		c.handler.Handle(core.NewStatusError(core.StatusInvalidArgument, err), "")
	} else {
		agent.Run()
	}
}
