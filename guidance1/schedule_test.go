package guidance1

import (
	"fmt"
	"time"
)

func ExampleNewScheduleAgent() {
	a := newScheduleAgent(time.Second*5, nil)

	fmt.Printf("test: newScheduleAgent() -> %v\n", a)

	//Output:
	//test: newScheduleAgent() -> schedule1

}
