package messagequeue

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"

	transport "github.com/AdityaMayukhSom/ruskin/transport"
)

type MessageQueueConfig struct {
	// stores the maximum number of topics which can be included in one message queue
	topicThreshold int

	// Factory to produce new stores
	factory StoreFactory
}

type MessageQueue struct {
	*MessageQueueConfig

	// It is used to create stores for a corresponding topic on demand.
	storeFactory StoreFactory

	// It is used when a goroutine is trying to check if a store
	// corresponding to a topic already exists or not.
	getStoreStateMut sync.RWMutex

	// It is used in when a goroutine is willing to update topicStores.
	changeStoreStateMut sync.Mutex

	// To update topic store in a goroutine safe manner, the developer must
	// lock the changeStoreStateMut before testing for the condition and
	// must write lock getStoreStateMut before actually updating the state
	// and then they must uplock them in the reverse manner.
	topicStores map[string]Store

	produceChannel <-chan transport.Message
}

func NewMessageQueue(config *MessageQueueConfig) (*MessageQueue, error) {
	mq := &MessageQueue{
		MessageQueueConfig: config,
		storeFactory:       config.factory,
		topicStores:        make(map[string]Store),
	}

	return mq, nil
}

// Returns true if a store for the topic exists, otherwise returns false.
func (mq *MessageQueue) checkStore(topicName string) (bool, error) {
	if len(topicName) == 0 {
		// no store can exist with an empty topic
		return false, errors.New("cannot check for store with empty name")
	}

	// Checks whether a topic with the name exists in the map or not.
	// ok will be true if it exists, if it doesn't, ok will be false
	mq.getStoreStateMut.RLock()
	_, found := mq.topicStores[topicName]
	mq.getStoreStateMut.RUnlock()

	// errors.New("unable to check if store exists or not")
	return found, nil
}

// Creates a new topic if the topic does not already exists.
//
// Returns true if topic is successfully created, otherwise returns false.
//
// This is for internally creating a store, with the help of mutexes for
// concurrency management, this function should not be exposed outside of
// the module, instead -
//
//	store, err := mq.GetStore(topicName)
//
// should always be used as it will check if a store for that topic already
// exists and will return that otherwise create a new store and return that.
func (mq *MessageQueue) createStore(topicName string) (Store, error) {
	if len(topicName) == 0 {
		// we cannot create a topic with length zero
		return nil, errors.New("cannot create topic store with empty name")
	}

	var topicStore Store

	// As we may require to modify the state of topicStores, we first lock
	// the changeStoreStateMut variable. Now no other function shall be able
	// to modify the state of topicStores (considering createStore is the
	// only way to modify topicStores)
	mq.changeStoreStateMut.Lock()

	found, err := mq.checkStore(topicName)
	if err == nil && !found {
		// err is nil here
		topicStore = mq.storeFactory.Produce(topicName)

		// we need to modify the state of the store
		// hence stop read access to check store
		mq.getStoreStateMut.Lock()

		// add the newly generated store for the given topic
		mq.topicStores[topicName] = topicStore

		mq.getStoreStateMut.Unlock()

		slog.Info("created store", "topic", topicName)
	} else if err != nil {
		// here the first condition fails because there has
		// been an error, in that case, we return nil as no new
		// store has been created and also return the not nil err
		// for the user to either try again or handle as per the
		// requirement by inspecting the error
		topicStore = nil
	} else if found {
		// err is nil here too
		topicStore = mq.topicStores[topicName]
	}

	// unlock the mutex so that any subsequent threads willing to create a
	// new topic can modify the topicStores map
	mq.changeStoreStateMut.Unlock()

	return topicStore, err
}

// This function gets the the corresponding store for a particular topic name.
// This function guarentees to return the store for that particular topic name
// in case any error does not occur while checking the existence or creating a
// new store as this will create a new store if previously store for a given
// did not exist.
//
// Returns nil as store when an error occurs.
func (mq *MessageQueue) GetStore(topicName string) (Store, error) {
	if len(topicName) == 0 {
		return nil, errors.New("cannot get topic store with empty name")
	}

	// this will store the reference to the store which has to be returned
	var topicStore Store

	found, err := mq.checkStore(topicName)

	if err == nil && !found {
		// no error and not found means the store does not exist
		if topicStore, err = mq.createStore(topicName); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else if found {
		topicStore = mq.topicStores[topicName]
	}

	return topicStore, nil
}

// Publishes the message to the topic mentioned in `Message.Topic` field.
//
// If the mentioned topic does not exist, this function will create a
// new topic and then publish the data into that topic.
//
// Returns the offset of the message inside the topic store as first value.
//
// Returns an error if the topic name is empty string or unable
// to insert data into the topic store. If an error is thrown then
// offset value returned is -1.
func (mq *MessageQueue) PublishMessage(message transport.Message) (int, error) {
	store, err := mq.GetStore(message.Topic)
	if err != nil {
		return -1, fmt.Errorf("cannot find store for topic %s", message.Topic)
	}

	offset, err := store.Insert(message.Data)
	if err != nil {
		return -1, fmt.Errorf("could not insert message=%s to topic=%s",
			string(message.Data), message.Topic)
	}

	return offset, nil
}

// Returns a boolean indicating whether the message queue is
// completely filled or not.
func (mq *MessageQueue) IsFull() bool {
	mq.getStoreStateMut.Lock()
	defer mq.getStoreStateMut.Unlock()

	topicsCount := len(mq.topicStores)

	return topicsCount == mq.topicThreshold
}
