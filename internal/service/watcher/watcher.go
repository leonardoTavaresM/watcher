package watcher

import (
	"time"

	"github.com/leonardoTavaresM/watcher/internal/domain"
)

type WatcherService struct {
	publisher domain.EventPublisher
}

func NewWatcherService(pub domain.EventPublisher) *WatcherService {
	return &WatcherService{publisher: pub}
}

func (s *WatcherService) HandleFileEvent(path, ext, evType string) error {
	event := domain.FileEvent{
		Timestamp: time.Now(),
		FilePath:  path,
		Ext:       ext,
		Event:     evType,
	}

	return s.publisher.Publish(event)
}
