package ingress1

import (
	"fmt"
	"github.com/advanced-go/stdlib/core"
	test "github.com/advanced-go/stdlib/messaging/messagingtest"
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
	//msg := messaging.NewControlMessage("to", "from", messaging.ShutdownEvent)
	c := newControllerAgent(origin, test.NewAgent())
	go runControl(c, nil, nil, nil, nil, nil)

	fmt.Printf("test: run(c,nil,nil) -> %v\n", "OK")
	//time.Sleep(time.Second * 8)
	//c.ctrlC <- msg
	//time.Sleep(time.Second * 1)

	//Output:
	//test: run(c,nil,nil) -> OK

}
