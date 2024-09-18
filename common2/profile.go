package common2

import (
	"github.com/advanced-go/events/threshold1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

func GetProfile(h core.ErrorHandler, agentId string, msg *messaging.Message) *threshold1.Profile {
	if !msg.IsContentType(ContentTypeProfile) {
		return nil
	}
	if p, ok := msg.Body.(*threshold1.Profile); ok {
		return p
	}
	h.Handle(ProfileTypeErrorStatus(agentId, msg.Body))
	return nil
}
