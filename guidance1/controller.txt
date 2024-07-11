package guidance

import (
	"github.com/advanced-go/stdlib/core"
	"golang.org/x/time/rate"
	"time"
)

// This is only used by the client envoy on startup.
//
// Ingress and egress default timeouts can be set for any host
// If host timeouts are not set, then the default is used.
// Timeouts are stored as a string, so a conversion is needed for time.Duration
// These are read only and are provided as a generic start when a controller is initialized
// These can be changed dynamically way in memory but are never reset in the database.
// Storing current values, which reflect current processing could cause an issue as current
// processing does not equal startup processing.
//
// No Rate Limits or Rate Bursts are stored as this can cause a processing error.
// Only defaults

const (
	IngressDefault   = time.Millisecond * 2500
	EgressDefault    = time.Millisecond * 2500
	RateLimitDefault = rate.Limit(100)
	RateBurstDefault = 25
)

type Entry struct {
	Region    string `json:"region"`
	Zone      string `json:"zone"`
	SubZone   string `json:"sub-zone"`
	Host      string `json:"host"`
	Status    string `json:"status"`
	CreatedTS string `json:"created-ts"`
	UpdatedTS string `json:"updated-ts"`

	// Timeout
	IngressTimeout string `json:"ingress-timeout"`
	EgressTimeout  string `json:"egress-timeout"`
}

func IngressTimeout(origin core.Origin) time.Duration {
	return IngressDefault
}

func EgressTimeout(origin core.Origin) time.Duration {
	return EgressDefault
}
