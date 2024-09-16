package redirect1

import "github.com/advanced-go/stdlib/messaging"

type Channel struct {
	C       chan *messaging.Message
	Enabled bool
}

func NewChannel(enable bool) Channel {
	c := Channel{}
	c.C = make(chan *messaging.Message, messaging.ChannelSize)
	c.Enabled = enable
	return c
}

func (c Channel) Close() {
	if c.C != nil {
		close(c.C)
	}
}

func (c Channel) Send(m *messaging.Message) {
	if m != nil && c.Enabled {
		c.C <- m
	}
}
