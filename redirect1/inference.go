package redirect1

import "github.com/advanced-go/stdlib/core"

// TODO: Given an observation of actual and limit for both percentiles and status codes
//       determine if the redirect has failed.

func runInference(r *redirect, obs *observation) *core.Status {
	return core.StatusOK()
}
