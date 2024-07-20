package guidance1

import (
	"context"
	"fmt"
	"github.com/advanced-go/stdlib/core"
	fmt2 "github.com/advanced-go/stdlib/fmt"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

// run - schedule
func runSchedule(s *schedule) {
	if s == nil {
		return
	}
	s.ticker.Start(0)

	for {
		select {
		case <-s.ticker.C():
			// TODO : poll database for global calendar change

		case msg := <-s.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				close(s.ctrlC)
				s.ticker.Stop()
				testLog(nil, s.uri, messaging.ShutdownEvent)
				return
			default:
			}
		default:
		}
	}
}

func testLog(_ context.Context, agentId string, content any) *core.Status {
	fmt.Printf("test: activity1.Log() -> %v : %v : %v\n", fmt2.FmtRFC3339Millis(time.Now().UTC()), agentId, content)
	return core.StatusOK()
}

func proc(origin core.Origin, h core.ErrorHandler) bool {
	_, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	//_, status := schedule1.Get(ctx, origin)
	//h.Handle(status, "")
	return true
}
