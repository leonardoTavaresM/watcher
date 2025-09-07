package watcher

import (
	"sync"
	"time"

	"github.com/leonardoTavaresM/watcher/internal/domain"
)

type WatcherService struct {
	publisher domain.EventPublisher
	lastEvent map[string]time.Time
	mu        sync.Mutex
	debounce  time.Duration
}

func NewWatcherService(pub domain.EventPublisher) *WatcherService {
	return &WatcherService{
		publisher: pub,
		lastEvent: make(map[string]time.Time),
		debounce:  750 * time.Millisecond,
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

	return s.publisher.Publish(event)
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
