package common

import (
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	Peak              = "peak"
	OffPeak           = "off-peak"
	ScaleUp           = "scale-up"
	ScaleDown         = "scale-down"
	PeakDuration      = time.Minute * 1
	OffPeakDuration   = time.Minute * 5
	ScaleUpDuration   = time.Minute * 2
	ScaleDownDuration = time.Minute * 2
)

type Window struct {
	Hour int
	Tag  string // Peak,Off-Peak,Scale-Up,Scale-Down
	Rate int
}

func NewWindow(hour, rate int, tag string) *Window {
	w := new(Window)
	w.Rate = rate
	w.Hour = hour
	w.Tag = tag
	return w
}

func (w *Window) IsScaleUp() bool {
	return w.Tag == ScaleUp
}

type Profile struct {
	current *Window
	next    *Window
}

func NewProfile() *Profile {
	p := new(Profile)
	p.current = NewWindow(5, 10, OffPeak)
	p.next = NewWindow(6, 15, ScaleUp)
	return p
}

func (p *Profile) IsOffPeak() bool {
	return p.current.Tag == OffPeak
}

func (p *Profile) Next() *Window {
	return p.next
}

func (p *Profile) ResiliencyDuration(old time.Duration) time.Duration {
	switch p.current.Tag {
	case Peak:
		return PeakDuration
	case OffPeak:
		return OffPeakDuration
	case ScaleUp:
		return ScaleUpDuration
	default: // ScaleDown:
		return ScaleDownDuration
	}
}

func (p *Profile) CaseOfficerDuration() time.Duration {
	switch p.current.Tag {
	case OffPeak:
		return PeakDuration
	default: // ScaleUp, ScaleDown, Peak:
		return OffPeakDuration
	}
}

func GetProfile(h core.ErrorHandler, agentId string, msg *messaging.Message) *Profile {
	if !msg.IsContentType(ContentTypeProfile) {
		return nil
	}
	if p, ok := msg.Body.(*Profile); ok {
		return p
	}
	h.Handle(ProfileTypeErrorStatus(agentId, msg.Body), "")
	return nil
}
