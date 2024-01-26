package consumer

type ConsumerRelay interface {
	AddConsumer(consumer *Consumer) bool
	Relay() error
}
