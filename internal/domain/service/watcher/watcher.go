package watcher

import (
	"errors"
	"sync"
	"time"

	"github.com/leonardoTavaresM/watcher/internal/adapter/consolepub"
	"github.com/leonardoTavaresM/watcher/internal/domain"
	"github.com/leonardoTavaresM/watcher/internal/domain/repository/memory"
)

type WatcherService struct {
	repository *memory.InMemoryEvent
	publisher  *consolepub.ConsolePublisher
	lastEvent  map[string]time.Time
	mu         sync.Mutex
	debounce   time.Duration
}

func NewWatcherService(memory *memory.InMemoryEvent, publisher *consolepub.ConsolePublisher) *WatcherService {
	return &WatcherService{
		repository: memory,
		publisher:  publisher,
		lastEvent:  make(map[string]time.Time),
		debounce:   750 * time.Millisecond,
	}
}

func (s *WatcherService) HandleFileEvent(path, ext, evType string) error {
	if ok := s.ShouldProcess(path); !ok {
		return nil
	}

	event := domain.FileEvent{
		Timestamp: time.Now(),
		FilePath:  path,
		Ext:       ext,
		Event:     evType,
	}

	err := s.repository.SaveInMemory(event)
	if err != nil {
		return errors.New(err.Error())
	}

	s.publisher.Publish()

	return nil
}

func (s *WatcherService) ShouldProcess(path string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	if last, ok := s.lastEvent[path]; ok {
		if now.Sub(last) < s.debounce {
			return false
		}
	}
	s.lastEvent[path] = now

	return true
}
