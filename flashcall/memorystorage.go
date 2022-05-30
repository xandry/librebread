package flashcall

import (
	"sync"
	"time"
)

type memRecord struct {
	at   time.Time
	to   string
	from string
	code string
}

type MemoryStorage struct {
	mu      sync.Mutex
	records []memRecord
}

func (stor *MemoryStorage) addRecord(to string, from string, code string) {
	stor.mu.Lock()
	defer stor.mu.Unlock()

	stor.records = append(stor.records, memRecord{
		at:   time.Now(),
		to:   to,
		from: from,
		code: code,
	})
}

func (stor *MemoryStorage) allRecordsReversed() []memRecord {
	allRecords := make([]memRecord, 0, len(stor.records))

	for i := len(stor.records) - 1; i >= 0; i-- {
		allRecords = append(allRecords, stor.records[i])
	}

	return allRecords
}
