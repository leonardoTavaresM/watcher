package repository

import (
	"errors"
	"fmt"

	"github.com/leonardoTavaresM/watcher/internal/domain/entity"
)

type InMemoryRepository struct {
	events map[int]entity.FileEvent
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		events: make(map[int]entity.FileEvent),
	}
}

func (m *InMemoryRepository) Save(event entity.FileEvent) error {
	if event == (entity.FileEvent{}) {
		return errors.New("event is empty")
	}

	fmt.Println("evento chegando", event)
	m.events[len(m.events)] = event
	fmt.Println("sendo salvo", m.events)
	return nil
}

func (m *InMemoryRepository) GetAll() map[int]entity.FileEvent {
	fmt.Println("pegando eventos", m.events)
	return m.events
}

func (m *InMemoryRepository) GetByID(id int) entity.FileEvent {
	fmt.Println("pegando evento", m.events[id])
	return m.events[id]
}

func (m *InMemoryRepository) Delete(id int) {
	delete(m.events, id)
}
