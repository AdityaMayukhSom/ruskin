package replica

import (
	consumer "github.com/AdityaMayukhSom/ruskin/consumer"
	websocket "github.com/gorilla/websocket"
)

type LoadPartition struct {
	topicChannel    <-chan string
	consumerChannel <-chan consumer.Consumer
	relays          map[string]consumer.ConsumerRelay
}

func (lp *LoadPartition) ProcessStream(topicName string) {
	consumerRelay := lp.relays[topicName]
	// relay over connection can take time, but need to consider if sending the message
	// actually takes some time, because we aren't waiting for the reply, we are only
	// publishing the message to the client
	go consumerRelay.Relay()
}

func (lp *LoadPartition) HandleSubscription(conn *websocket.Conn, topics []string) {
	for _, topic := range topics {
		_, found := lp.relays[topic]
		if !found {
			// no idea if we need to use mutex for synchronization
			lp.relays[topic] = consumer.NewWSConsumerRelay()
		}
	}
}

func (lp *LoadPartition) Start() {

	for {
		select {
		case topicName := <-lp.topicChannel:
			lp.ProcessStream(topicName)
		case consumer := <-lp.consumerChannel:
			lp.HandleSubscription(consumer.Conn, consumer.Topics)
		}
	}

}
