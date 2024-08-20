package common

import (
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

type Profile struct {
	Hour int
	Tag  string // Peak,Off-Peak,Scale-Up,Scale-Down
	Rate int
}

func (p *Profile) ResiliencyDuration(old time.Duration) time.Duration {
	switch p.Tag {
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
	switch p.Tag {
	case OffPeak:
		return PeakDuration
	default: // ScaleUp, ScaleDown, Peak:
		return OffPeakDuration
	}
}
