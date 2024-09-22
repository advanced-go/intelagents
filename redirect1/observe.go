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
// lookBackDuration = time.Minute *
)

type redirectObservation struct {
	Percentile  common1.Observation
	StatusCodes common1.Observation
}

func newObservation(percentile, statusCodes common1.Observation) *redirectObservation {
	o := new(redirectObservation)
	o.Percentile = percentile
	o.StatusCodes = statusCodes
	return o
}

func getObservation(h core.ErrorHandler, agentId string, msg *messaging.Message) *redirectObservation {
	if !msg.IsContentType(contentTypeRedirectObservation) {
		return nil
	}
	if p, ok := msg.Body.(*redirectObservation); ok {
		return p
	}
	h.Handle(observationTypeErrorStatus(agentId, msg.Body))
	return nil
}

func observationTypeErrorStatus(agentId string, t any) *core.Status {
	err := errors.New(fmt.Sprintf("error: redirect observation type:%v is invalid for agent:%v", reflect.TypeOf(t), agentId))
	return core.NewStatusError(core.StatusInvalidArgument, err)
}

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
	actual, status1 := observe.StatusCodeThresholdQuery(h, origin, from, to, statusCodes)
	if !status1.OK() {
		return common1.Observation{}, status1
	}
	return common1.Observation{Actual: actual, Limit: limit}, core.StatusOK()
}
