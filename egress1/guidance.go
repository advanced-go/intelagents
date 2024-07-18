package egress1

import "time"

const (
	logisticsDuration = time.Second * 2
	officerDuration   = time.Second * 5
	scheduleDuration  = time.Second * 5
)

type guidance struct {
	//logisticsInterval   func() time.Duration
	//caseOfficerInterval func() time.Duration
	scheduleInterval func() time.Duration
}

func newGuidance() *guidance {
	return &guidance{
		//logisticsInterval:   func() time.Duration { return logisticsDuration },
		//caseOfficerInterval: func() time.Duration { return officerDuration },
		scheduleInterval: func() time.Duration { return scheduleDuration },
	}
}
