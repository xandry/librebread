package push

import (
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var ErrMessageNotFound = errors.New("message not found in store")

type SentMessage struct {
	ID       string
	Time     time.Time
	Provider string
	Msg      BatchMessage
}

type MemoryStorage struct {
	nextID func() string

	mutex    sync.Mutex
	byID     map[string]*SentMessage
	messages []SentMessage
}

func NewMemoryStorage() *MemoryStorage {
	lastID := int64(0)

	return &MemoryStorage{
		nextID: func() string {
			id := atomic.AddInt64(&lastID, 1)
			return strconv.FormatInt(id, 10)
		},
		byID: map[string]*SentMessage{},
	}
}

func (store *MemoryStorage) AddBatchMessage(provider string, msg BatchMessage) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	id := store.nextID()

	storeMessage := SentMessage{
		ID:       id,
		Time:     time.Now(),
		Provider: provider,
		Msg:      msg,
	}

	// reverse message order with fewer allocations
	store.messages = append(store.messages, SentMessage{})
	copy(store.messages[1:], store.messages)
	store.messages[0] = storeMessage

	store.byID[id] = &store.messages[0]

	return nil
}

func (store *MemoryStorage) AllMessages() ([]SentMessage, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	return store.messages, nil
}

func (store *MemoryStorage) ByID(id string) (SentMessage, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	msg, ok := store.byID[id]
	if !ok {
		return SentMessage{}, ErrMessageNotFound
	}

	return *msg, nil
}
