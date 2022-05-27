package flashcall

import "sync"

type MemoryStorage struct {
	mu     sync.Mutex
	phones []string
}

func (stor *MemoryStorage) AddPhone(phone string) {
	stor.mu.Lock()
	defer stor.mu.Unlock()

	stor.phones = append(stor.phones, phone)
}
