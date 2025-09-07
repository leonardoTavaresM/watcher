package consolepub

import (
	"encoding/json"
	"fmt"

	"github.com/leonardoTavaresM/watcher/internal/domain"
)

type ConsolePublisher struct{}

func NewConsolePublisher() *ConsolePublisher {
	return &ConsolePublisher{}
}

func (c *ConsolePublisher) Publish(event domain.FileEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	fmt.Println("publish", string(data))
	return nil
}
