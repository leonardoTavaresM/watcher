package service

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/leonardoTavaresM/watcher/internal/application/port"
	"github.com/leonardoTavaresM/watcher/internal/domain/entity"
)

type WatcherService struct {
	repository port.EventRepository
	publishers []port.Publisher
	lastEvent  map[string]time.Time
	mu         sync.Mutex
	debounce   time.Duration
}

func NewWatcherService(repository port.EventRepository, publishers ...port.Publisher) *WatcherService {
	return &WatcherService{
		repository: repository,
		publishers: publishers,
		lastEvent:  make(map[string]time.Time),
		debounce:   750 * time.Millisecond,
	}
}

func (s *WatcherService) HandleFileEvent(path, ext, evType string) error {
	if ok := s.ShouldProcess(path); !ok {
		return nil
	}

	event := entity.FileEvent{
		Timestamp: time.Now(),
		FilePath:  path,
		Ext:       ext,
		Event:     evType,
	}

	err := s.repository.Save(event)
	if err != nil {
		return errors.New(err.Error())
	}

	for _, publisher := range s.publishers {
		if err := publisher.Publish(event); err != nil {
			log.Printf("Failed to publish to %T: %v", publisher, err)
		}
	}

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
