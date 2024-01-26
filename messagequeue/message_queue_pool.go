package messagequeue

type MessageQueuePool interface {
	// This guarantees to provide a new message queue whenever this method
	// is invoked on a pool
	Get() (*MessageQueue, error)

	// Creates message queues and populates the pool.
	//
	// A new message queue is added to the pool each time this
	// method is invoked on an instance of message queue pool.
	Create() error

	Add(*MessageQueue) error

	// Deletes the specified message queue from the message queue pool if
	// that message queue existed in the pool previously. Otherwise does not
	// change the state of the pool.
	//
	// Consider: currently it used pointer to the message queue for identifying
	// the message queue but in future, we may need to use unique ID for a particular
	// message queue if we want to scale it in a distributed manner as different
	// computers in a queue swarm do not share the same address space.
	Remove(*MessageQueue) error
}
