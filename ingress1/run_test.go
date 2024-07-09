package ingress1

import (
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

func _ExampleRun_Reset() {
	origin := core.Origin{
		Region:     "us-central1",
		Zone:       "c",
		SubZone:    "",
		Host:       "www.host1.com",
		InstanceId: "",
	}
	msg := messaging.NewControlMessage("to", "from", messaging.ShutdownEvent)

	c := newControllerAgent(origin, newTestAgent())
	go run(c, nil, nil, nil, nil)
	time.Sleep(time.Second * 8)

	c.ctrlC <- msg
	time.Sleep(time.Second * 1)

	//Output:
	//test: activity1.Log() -> 2024-07-09T15:27:41.535Z : ingress-controller1:us-central1.c.www.host1.com : processing assignment
	//test: activity1.Log() -> 2024-07-09T15:27:42.528Z : ingress-controller1:us-central1.c.www.host1.com : processing assignment
	//test: activity1.Log() -> 2024-07-09T15:27:42.528Z : ingress-controller1:us-central1.c.www.host1.com : event:shutdown

}
