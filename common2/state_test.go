package common2

import (
	"fmt"
	"github.com/advanced-go/events/threshold1"
)

func ExampleNewObservation() {
	o := NewObservation(threshold1.Entry{}, threshold1.Entry{})

	fmt.Printf("test: NewObservation() -> [%v]\n", o)

	//Output:
	//fail
}
