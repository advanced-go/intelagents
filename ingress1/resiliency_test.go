package ingress1

import (
	"context"
	"fmt"
	"github.com/advanced-go/stdlib/core"
	fmt2 "github.com/advanced-go/stdlib/fmt"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

type testAgent struct{}

func newTestAgent() *testAgent                    { return new(testAgent) }
func (t *testAgent) Uri() string                  { return "testAgent" }
func (t *testAgent) Message(m *messaging.Message) { fmt.Printf("test: testAgent.Message() -> %v\n", m) }
func (t *testAgent) Handle(status *core.Status, _ string) *core.Status {
	fmt.Printf("test: opsAgent.Handle() -> [status:%v]\n", status)
	status.Handled = true
	return status
}
func (t *testAgent) Run()      {}
func (t *testAgent) Shutdown() {}

func testLog2(_ context.Context, agentId string, content any) *core.Status {
	fmt.Printf("test: activity1.Log() -> %v : %v : %v\n", fmt2.FmtRFC3339Millis(time.Now().UTC()), agentId, content)
	return core.StatusOK()
}

func ExampleControllerAgentUri() {
	origin := core.Origin{
		Region:     "us-central1",
		Zone:       "c",
		SubZone:    "sub-zone",
		Host:       "host",
		InstanceId: "",
	}
	u := resiliencyAgentUri(origin)
	fmt.Printf("test: AgentUri() -> [%v]\n", u)

	origin.Region = "us-west1"
	origin.Zone = "a"
	origin.SubZone = ""
	u = resiliencyAgentUri(origin)
	fmt.Printf("test: AgentUri() -> [%v]\n", u)

	//Output:
	//test: AgentUri() -> [ingress-controller1:us-central1.c.sub-zone.host]
	//test: AgentUri() -> [ingress-controller1:us-west1.a.host]

}
