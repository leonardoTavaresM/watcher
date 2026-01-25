package publisher

import (
	"encoding/json"
	"fmt"

	"github.com/leonardoTavaresM/watcher/internal/application/port"
	"github.com/leonardoTavaresM/watcher/internal/domain/entity"
)

type ConsolePublisher struct {
	repository port.EventRepository
}

func NewConsolePublisher(repository port.EventRepository) *ConsolePublisher {
	return &ConsolePublisher{
		repository: repository,
	}
}

func (c *ConsolePublisher) Publish(event entity.FileEvent) error {
	data, err := json.Marshal(c.repository.GetAll())
	if err != nil {
		return err
	}

	fmt.Println("publish", string(data))
	return nil
}

func (c *ConsolePublisher) Close() error {
	return nil
}
