package load

import (
	consumer "github.com/AdityaMayukhSom/ruskin/consumer"
	messagequeue "github.com/AdityaMayukhSom/ruskin/messagequeue"
	websocket "github.com/gorilla/websocket"
)

type LoadPartition struct {

	// TODO: we need some kind of interface named topic identifier which
	// can be used to identify topics over a cluster of computers, and then
	// replace the *messagequeue.Store in topicChannel with that interface
	// which can be used to return a unique identifier, preferebly a UUID or
	// IP address with port of the computer in which that port runs, and also
	// an interface Connector needs to be implemented which can be used to
	// establish connection both over the nerwork and via pointer if same address
	// space is used (can be logically same or physically same)
	topicChannel    <-chan *messagequeue.Store
	consumerChannel <-chan consumer.Consumer
	relays          map[*messagequeue.Store]consumer.ConsumerRelay
}

func (lp *LoadPartition) ProcessStream(topicName *messagequeue.Store) {
	consumerRelay := lp.relays[topicName]
	// relay over connection can take time, but need to consider if sending the message
	// actually takes some time, because we aren't waiting for the reply, we are only
	// publishing the message to the client
	go consumerRelay.Relay()
}

func (lp *LoadPartition) HandleSubscription(conn *websocket.Conn, topic *messagequeue.Store) {
	relay, found := lp.relays[topic]
	if !found {
		relay = consumer.NewWSConsumerRelay(topic)
		// no idea if we need to use mutex for synchronization
		lp.relays[topic] = relay
	}

}

func (lp *LoadPartition) Start() {
	for {
		select {
		case topicName := <-lp.topicChannel:
			go lp.ProcessStream(topicName)
		case consumer := <-lp.consumerChannel:
			go lp.HandleSubscription(consumer.Conn, consumer.Topic)
		}
	}
}
