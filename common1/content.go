package common1

import (
	"errors"
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"reflect"
)

const (
	ContentTypeProfile            = "application/profile"
	ContentTypePercentileSLO      = "application/percentileSLO"
	ContentTypeRedirectConfig     = "application/redirect-config"
	ContentTypeEgressConfig       = "application/egress-config"
	ContentTypeIngressTimeseries1 = "application/ingress-timeseries1"

	ContentTypeEgressTimeseries1 = "application/egress-timeseries1"

	SQLInsert = "insert"
	SQLUpdate = "update"
	SQLDelete = "delete"
)

func ProfileTypeErrorStatus(agentId string, t any) *core.Status {
	err := errors.New(fmt.Sprintf("error: profile data change type:%v is invalid for agent:%v", reflect.TypeOf(t), agentId))
	return core.NewStatusError(core.StatusInvalidArgument, err)
}

func MessageEventErrorStatus(agentId string, msg *messaging.Message) *core.Status {
	err := errors.New(fmt.Sprintf("error: message event:%v is invalid for agent:%v", msg.Event(), agentId))
	return core.NewStatusError(core.StatusInvalidArgument, err)
}

func MessageContentTypeErrorStatus(agentId string, msg *messaging.Message) *core.Status {
	err := errors.New(fmt.Sprintf("error: message content:%v is invalid for agent:%v and event:%v", msg.ContentType(), agentId, msg.Event()))
	return core.NewStatusError(core.StatusInvalidArgument, err)
}
