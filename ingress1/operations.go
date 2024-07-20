package ingress1

import (
	"fmt"
	"github.com/advanced-go/stdlib/core"
)

const ()

type operations struct {
	log func(uri string, content any)
}

func newOperations(handler func(status *core.Status, _ string) *core.Status) *operations {
	if handler == nil {
		handler = func(status *core.Status, _ string) *core.Status {
			return status
		}
	}
	return &operations{
		log: func(uri string, content any) {
			//ctx, cancel := context.WithTimeout(context.Background(), logDuration)
			//defer cancel()
			status := core.StatusOK()
			fmt.Printf("%v - %v", uri, content)
			if !status.OK() {
				handler(status, "")
			}
		},
	}
}
