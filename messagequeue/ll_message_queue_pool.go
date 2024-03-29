package messagequeue

type node struct {
	mq   *MessageQueue
	next *node
}

type LinkedListMessageQueuePool struct {
	length int
	head   *node
	tail   *node
}

func (mqp *LinkedListMessageQueuePool) Get() (*MessageQueue, error) {
	if mqp.head == nil {
		// create a new empty MessageQueue and add that to the pool and
		// then return the pointer to the newly created MessageQueue

	}

	return nil, nil
}

// Creates message queues and populates the pool.
//
// A new message queue is added to the pool each time this
// method is invoked on an instance of message queue pool.
func (mqp *LinkedListMessageQueuePool) Create() error {
	return nil
}

func (mqp *LinkedListMessageQueuePool) Add(mq *MessageQueue) error {

	return nil
}

// Deletes the specified message queue from the message queue pool if
// that message queue existed in the pool previously. Otherwise does not
// change the state of the pool.
//
// Consider: currently it used pointer to the message queue for identifying
// the message queue but in future, we may need to use unique ID for a particular
// message queue if we want to scale it in a distributed manner as different
// computers in a queue swarm do not share the same address space.
func (mqp *LinkedListMessageQueuePool) Delete(mq *MessageQueue) error {
	if mqp.length == 1 && mqp.head.mq == mq {
		// Garbage collector will clean up the wasted memory.
		mqp.head = nil
		return nil
	}

	// stores the current node that is being examined. keep in mind that
	// to delete a node from a linked list, we actually need a pointer
	// to the node previous to that node, i.e.
	tempNode := mqp.head

	for tempNode != nil {
		if tempNode.next != nil && tempNode.next.mq == mq {
			// The required message queue is deleted, so we can return early.
			// If for most of the cases, the message queue is found near the
			// head, early return can lead to significant performance boosts.
			tempNode.next = tempNode.next.next
			return nil
		}
		tempNode = tempNode.next
	}

	return nil
}
