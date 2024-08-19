package common

import (
	"errors"
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"reflect"
)

func ProfileTypeErrorStatus(t any) *core.Status {
	return core.NewStatusError(core.StatusInvalidArgument, errors.New(fmt.Sprintf("error: data change profile type is invalid [%v]", reflect.TypeOf(t))))
}

func PercentileTypeErrorStatus(t any) *core.Status {
	return core.NewStatusError(core.StatusInvalidArgument, errors.New(fmt.Sprintf("error: data change percentile type is invalid [%v]", reflect.TypeOf(t))))
}

func RedirectPlanTypeErrorStatus(t any) *core.Status {
	return core.NewStatusError(core.StatusInvalidArgument, errors.New(fmt.Sprintf("error: data change reditect plan type is invalid [%v]", reflect.TypeOf(t))))
}
