package messagequeue

type MessageQueuePool interface {
	Get() (*MessageQueue, error)
	Create() error
	Add(*MessageQueue) error
	Remove(*MessageQueue) error
}
