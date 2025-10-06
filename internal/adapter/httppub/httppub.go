package httppub

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/leonardoTavaresM/watcher/internal/domain/repository/memory"
)

type HttpAdapter struct {
	repository memory.InMemoryEvent
}

func NewHTTPAdapter(repository *memory.InMemoryEvent) *HttpAdapter {
	return &HttpAdapter{
		repository: *repository,
	}
}

func (h *HttpAdapter) GetAllEvents(c *fiber.Ctx) error {
	events := h.repository.GetEvents()

	response, err := ToEventsResponse(events)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	fmt.Println("response GetAllEvents", response)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *HttpAdapter) GetEvent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	event := h.repository.GetEvent(id)

	response, err := ToEventResponse(id, event)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	fmt.Println("response GetEvent", response)
	return c.Status(fiber.StatusOK).JSON(response)
}
