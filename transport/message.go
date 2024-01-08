package transport

import "fmt"

type Message struct {
	Topic string
	Data  []byte
}

func (m Message) String() string {
	return fmt.Sprintf("(topic :: %s, data :: %s)", m.Topic, string(m.Data))
}
