package consolepub

import (
	"encoding/json"
	"fmt"

	"github.com/leonardoTavaresM/watcher/internal/domain"
	"github.com/leonardoTavaresM/watcher/internal/domain/repository/memory"
)

type ConsolePublisher struct {
	repository *memory.InMemoryEvent
}

// Close implements domain.Publisher.
func (c *ConsolePublisher) Close() error {
	panic("unimplemented")
}

func NewConsolePublisher(repository *memory.InMemoryEvent) *ConsolePublisher {
	return &ConsolePublisher{
		repository: repository,
	}
}

func (c *ConsolePublisher) Publish(event domain.FileEvent) error {
	data, err := json.Marshal(c.repository.GetEvents())
	if err != nil {
		return err
	}

	fmt.Println("publish", string(data))
	return nil
}
