package common1

import (
	"fmt"
	"github.com/advanced-go/events/timeseries1"
)

func ExampleNewObservation() {
	o := NewObservation(timeseries1.Threshold{}, timeseries1.Threshold{})

	fmt.Printf("test: NewObservation() -> [%v]\n", o)

	//Output:
	//fail
}
