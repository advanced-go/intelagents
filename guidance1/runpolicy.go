package guidance1

import "github.com/advanced-go/stdlib/messaging"

// run - policy change
func runPolicy(p *policy) {
	if p == nil {
		return
	}
	p.ticker.Start(0)

	for {
		select {
		case <-p.ticker.C():
			// TODO : poll database for controller policy changes

		case msg := <-p.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				close(p.ctrlC)
				p.ticker.Stop()
				testLog(nil, p.uri, messaging.ShutdownEvent)
				return
			default:
			}
		default:
		}
	}
}
