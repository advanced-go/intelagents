package common

import (
	"errors"
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"reflect"
)

func ProfileTypeErrorStatus(agentId string, t any) *core.Status {
	err := errors.New(fmt.Sprintf("error: profile data change type:%v is invalid for agent:%v", reflect.TypeOf(t), agentId))
	return core.NewStatusError(core.StatusInvalidArgument, err)
}

/*
func PercentileTypeErrorStatus(t any) *core.Status {
	return core.NewStatusError(core.StatusInvalidArgument, errors.New(fmt.Sprintf("error: data change percentile type is invalid [%v]", reflect.TypeOf(t))))
}


*/

func RedirectPlanTypeErrorStatus(agentId string, t any) *core.Status {
	err := errors.New(fmt.Sprintf("error: redirect data change type:%v is invalid for agent:%v", reflect.TypeOf(t), agentId))
	return core.NewStatusError(core.StatusInvalidArgument, err)
}

func MessageEventErrorStatus(agentId string, msg *messaging.Message) *core.Status {
	err := errors.New(fmt.Sprintf("error: message event:%v is invalid for agent:%v", msg.Event(), agentId))
	return core.NewStatusError(core.StatusInvalidArgument, err)
}
