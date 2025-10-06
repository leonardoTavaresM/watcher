package memory

import (
	"errors"
	"fmt"

	"github.com/leonardoTavaresM/watcher/internal/domain"
)

type InMemoryEvent struct {
	events map[int]domain.FileEvent
}

func NewInMemoryEvent() *InMemoryEvent {
	return &InMemoryEvent{
		events: make(map[int]domain.FileEvent), // Inicializar o mapa!
	}
}

func (m *InMemoryEvent) SaveInMemory(event domain.FileEvent) error {
	if event == (domain.FileEvent{}) {
		return errors.New("event is empty")
	}

	fmt.Println("evento chegando", event)
	m.events[len(m.events)] = event
	fmt.Println("sendo salvo", m.events)
	return nil
}

func (m *InMemoryEvent) GetEvents() map[int]domain.FileEvent {
	fmt.Println("pegando eventos", m.events)
	return m.events
}

func (m *InMemoryEvent) GetEvent(id int) domain.FileEvent {
	fmt.Println("pegando evento", m.events[id])
	return m.events[id]
}

func (m *InMemoryEvent) DeleteEvent(id int) {
	delete(m.events, id)
}
