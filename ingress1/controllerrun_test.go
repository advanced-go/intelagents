package ingress1

import (
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	test "github.com/advanced-go/stdlib/messaging/messagingtest"
	"time"
)

func testHandler(status *core.Status, _ string) *core.Status {
	fmt.Printf("test: testHandler() -> [status:%v]\n", status)
	return status
}

func ExampleRun_Nil() {
	origin := core.Origin{
		Region:     "us-central1",
		Zone:       "c",
		SubZone:    "",
		Host:       "www.host1.com",
		InstanceId: "",
	}
	dur := time.Second * 2
	msg := messaging.NewControlMessage("to", "from", messaging.ShutdownEvent)
	c := newController(origin, test.NewAgent(), dur, dur, dur)
	go controllerRun(c, nil, nil, nil, nil, nil, nil)

	fmt.Printf("test: run(c,nil,nil) -> %v\n", "OK")
	time.Sleep(time.Second * 8)
	c.ctrlC <- msg
	time.Sleep(time.Second * 1)

	//Output:
	//test: run(c,nil,nil) -> OK

}
