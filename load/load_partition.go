package load

import (
	src "github.com/AdityaMayukhSom/ruskin/connector/store_relay"
	consumer "github.com/AdityaMayukhSom/ruskin/consumer"
	messagequeue "github.com/AdityaMayukhSom/ruskin/messagequeue"
)

type LoadPartition struct {
	/*
		Whenever a new message is generated in a topic, that topic is pushed to this
		channel so that LoadPartition can take appropriate actions and invoke the
		consumer relay to relay that message to the connected consumers.

		TODO: we need some kind of interface named `TopicIdentifier` which
		can be used to identify topics over a cluster of computers and then
		replace the *messagequeue.Store in topicChannel with that interface
		which can be used to return a unique identifier, preferebly a UUID or
		IP address with port of the computer in which that topic exists.

		Also an interface `Connector` needs to be implemented which can be used to
		establish connection both over the nerwork and via pointer if same address
		space is used (can be either logically same or physically same)
	*/
	topicChannel <-chan messagequeue.TopicIdentifier

	connnectionChannel <-chan struct {
		consumer consumer.Consumer
		topic    messagequeue.TopicIdentifier
	}

	// A map to store the in which relay the consumer exists.
	relays map[messagequeue.TopicIdentifier]consumer.ConsumerRelay
}

func (lp *LoadPartition) ProcessStream(topic messagequeue.TopicIdentifier) {
	consumerRelay := lp.relays[topic]
	// relay over connection can take time, but need to consider if sending the message
	// actually takes some time, because we aren't waiting for the reply, we are only
	// publishing the message to the client
	go consumerRelay.Relay()
}

func (lp *LoadPartition) HandleSubscription(cnsmr consumer.Consumer, topic messagequeue.TopicIdentifier) {
	relay, found := lp.relays[topic]
	if !found {
		// TODO: create the topic in message queue and then get the pointer to that newly created topic
		storeRelayConnector := src.NewPointerStoreConnector(nil)

		relay = consumer.NewWSConsumerRelay(storeRelayConnector)
		// TODO: no idea if we need to use mutex for synchronization
		lp.relays[topic] = relay
	}
	// addition of that consumer can be done in another go routine
	go relay.AddConsumer(&cnsmr)
}

func (lp *LoadPartition) Start() {
	for {
		select {
		// TODO: we need to depend on an interface to give us the TopicIdentifier
		// rather than depending on the topic channel directly. We can create
		// an interface LoadPartition with GetTopic() which internally handles
		// how to get the notification of the arrival of a new message under a topic.
		case topicName := <-lp.topicChannel:
			go lp.ProcessStream(topicName)
		case con := <-lp.connnectionChannel:
			go lp.HandleSubscription(con.consumer, con.topic)
		}
	}
}
