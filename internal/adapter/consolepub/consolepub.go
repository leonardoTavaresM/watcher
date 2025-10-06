package consolepub

import (
	"encoding/json"
	"fmt"

	"github.com/leonardoTavaresM/watcher/internal/domain/repository/memory"
)

type ConsolePublisher struct {
	repository *memory.InMemoryEvent
}

func NewConsolePublisher(repository *memory.InMemoryEvent) *ConsolePublisher {
	return &ConsolePublisher{
		repository: repository,
	}
}

func (c *ConsolePublisher) Publish() error {
	data, err := json.Marshal(c.repository.GetEvents())
	if err != nil {
		return err
	}

	fmt.Println("publish", string(data))
	return nil
}
