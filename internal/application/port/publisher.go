package port

import "github.com/leonardoTavaresM/watcher/internal/domain/entity"

type Publisher interface {
	Publish(event entity.FileEvent) error
	Close() error
}
