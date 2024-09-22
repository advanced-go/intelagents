package redirect1

import (
	"errors"
	"fmt"
	"github.com/advanced-go/events/timeseries1"
	"github.com/advanced-go/intelagents/common1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"reflect"
	"time"
)

const (
	lookBackDuration = time.Minute * 30
)

type observation struct {
	percentile common1.Observation
	statusCode common1.Observation
}

func newObservation(percentile, statusCode common1.Observation) *observation {
	o := new(observation)
	o.percentile = percentile
	o.statusCode = statusCode
	return o
}

func castObservation(h core.ErrorHandler, agentId string, msg *messaging.Message) *observation {
	if !msg.IsContentType(contentTypeRedirectObservation) {
		return nil
	}
	if p, ok := msg.Body.(*observation); ok {
		return p
	}
	h.Handle(observationTypeErrorStatus(agentId, msg.Body))
	return nil
}

func observationTypeErrorStatus(agentId string, t any) *core.Status {
	err := errors.New(fmt.Sprintf("error: redirect observation type:%v is invalid for agent:%v", reflect.TypeOf(t), agentId))
	return core.NewStatusError(core.StatusInvalidArgument, err)
}

type observer struct {
	limitPercent    timeseries1.Threshold
	limitStatusCode timeseries1.Threshold
}

func newObserver(r *redirect, events *common1.Events) *observer {
	o := new(observer)
	common1.SetPercentileThreshold(r.handler, r.origin, &o.limitPercent, events)
	to := time.Now().UTC()
	common1.SetStatusCodesThreshold(r.handler, r.origin, &o.limitStatusCode, to.Add(-lookBackDuration), to, common1.DefaultStatusCodes, events)
	return o
}

func (o *observer) observation(r *redirect, origin core.Origin, events *common1.Events) (*observation, *core.Status) {
	to := time.Now().UTC()
	from := to.Add(-(r.state.StepDuration - time.Minute))
	actualPercent, status := events.PercentThresholdQuery(r.handler, origin, from, to)
	if !status.OK() {
		return nil, status
	}
	actualStatusCode, status1 := events.StatusCodeThresholdQuery(r.handler, origin, from, to, common1.DefaultStatusCodes)
	if !status1.OK() {
		return nil, status1
	}
	return newObservation(common1.Observation{
		Actual: actualPercent,
		Limit:  o.limitPercent,
	}, common1.Observation{
		Actual: actualStatusCode,
		Limit:  o.limitStatusCode,
	},
	), core.StatusOK()

}

/*

func percentileObservation(h core.ErrorHandler, origin core.Origin, limit timeseries1.Threshold, stepDuration time.Duration, observe *common1.Events) (common1.Observation, *core.Status) {
	to := time.Now().UTC()
	from := to.Add(-(stepDuration - time.Minute))
	actual, status := observe.PercentThresholdQuery(h, origin, from, to)
	if !status.OK() {
		return common1.Observation{}, status
	}
	return common1.Observation{Actual: actual, Limit: limit}, core.StatusOK()
}

func statusCodeObservation(h core.ErrorHandler, origin core.Origin, limit timeseries1.Threshold, stepDuration time.Duration, statusCodes string, observe *common1.Events) (common1.Observation, *core.Status) {
	to := time.Now().UTC()
	from := to.Add(-(stepDuration - time.Minute))

	return common1.Observation{Actual: actual, Limit: limit}, core.StatusOK()
}





*/
