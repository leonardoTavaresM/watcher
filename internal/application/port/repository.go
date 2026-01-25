package port

import "github.com/leonardoTavaresM/watcher/internal/domain/entity"

type EventRepository interface {
	Save(event entity.FileEvent) error
	GetAll() map[int]entity.FileEvent
	GetByID(id int) entity.FileEvent
	Delete(id int)
}
