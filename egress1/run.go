package egress1

import (
	"context"
	"github.com/advanced-go/observation/access1"
	"github.com/advanced-go/observation/inference1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"net/http"
	"time"
)

type queryAccessFunc func(ctx context.Context, origin core.Origin) ([]access1.Entry, *core.Status)
type queryInferenceFunc func(ctx context.Context, origin core.Origin) ([]inference1.Entry, *core.Status)
type getGuidanceFunc func()
type insertInferenceFunc func(ctx context.Context, h http.Header, e inference1.Entry) *core.Status

// run - egress controller
func run(c *controller, access queryAccessFunc, inference queryInferenceFunc, guidance getGuidanceFunc, insert insertInferenceFunc) {
	if c == nil {
		return
	}
	tick := time.Tick(c.interval)

	for {
		select {
		case <-tick:
			status := core.StatusOK()
			if !status.OK() && !status.NotFound() {
				c.handler.Message(messaging.NewStatusMessage(c.handler.Uri(), c.uri, status))
			}
		case msg, open := <-c.ctrlC:
			if !open {
				return
			}
			switch msg.Event() {
			case messaging.ShutdownEvent:
				return
			default:
			}
		default:
		}
	}
}

/*
func updateAssignments(partition landscape1.Entry) ([]assignment1.Entry, *core.Status) {
	values := make(url.Values)
	values.Add(core.RegionKey, partition.Region)
	values.Add(core.ZoneKey, partition.Zone)
	values.Add(core.SubZoneKey, partition.SubZone)
	entries, _, status := assignment1.Get(nil, nil, values)
	return entries, status
}

func logActivity(body []activity1.Entry) *core.Status {
	_, status := activity1.Put(nil, body)
	return status
}

func processAssignments(c *caseOfficer, log func(body []activity1.Entry) *core.Status, update func(partition landscape1.Entry) ([]assignment1.Entry, *core.Status)) *core.Status {
	status := log([]activity1.Entry{{AgentId: c.uri}})
	if !status.OK() {
		return status
	}
	entries, status1 := update(c.partition)
	if !status1.OK() {
		return status
	}
	if len(entries) > 0 {
	}
	return status
}


*/
