package messagequeue

type MessageQueueConfig struct {
	// stores the maximum number of topics which can be included in one message queue
	topicThreshold int

	// Factory to produce new stores
	factory StoreFactory
}
