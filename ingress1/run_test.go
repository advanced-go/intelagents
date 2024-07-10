package ingress1

import (
	"fmt"
	"github.com/advanced-go/stdlib/core"
)

func ExampleRun_Nil() {
	origin := core.Origin{
		Region:     "us-central1",
		Zone:       "c",
		SubZone:    "",
		Host:       "www.host1.com",
		InstanceId: "",
	}
	//msg := messaging.NewControlMessage("to", "from", messaging.ShutdownEvent)
	c := newControllerAgent(origin, newTestAgent())
	go run(c, nil, nil)

	fmt.Printf("test: run(c,nil,nil) -> %v\n", "OK")
	//time.Sleep(time.Second * 8)
	//c.ctrlC <- msg
	//time.Sleep(time.Second * 1)

	//Output:
	//test: run(c,nil,nil) -> OK

}
