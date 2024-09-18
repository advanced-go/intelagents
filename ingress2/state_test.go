package ingress2

import (
	"fmt"
	"github.com/advanced-go/stdlib/core"
)

func ExampleSetRateLimiting() {
	status := core.StatusOK()

	fmt.Printf("test: setRateLimiting() -> [%v]\n", status)

	//Output:
	//fail
}
