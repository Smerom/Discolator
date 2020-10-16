package characterTracker

import (
	"sync"
	"unicode/utf8"
)

type memoryTracker struct {
	sync.RWMutex
	characterCount int
}

func NewMemoryTracker() Tracker {
	return &memoryTracker{}
}

func (mt *memoryTracker) AddCharacters(s string) int {
	mt.Lock()
	defer mt.Unlock()

	count := utf8.RuneCountInString(s)
	mt.characterCount += count

	return mt.characterCount
}

func (mt *memoryTracker) CountAfterString(s string) int {
	mt.RLock()
	defer mt.RUnlock()

	count := utf8.RuneCountInString(s)
	return mt.characterCount + count
}
