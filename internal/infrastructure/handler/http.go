package handler

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/leonardoTavaresM/watcher/internal/application/port"
)

type HTTPHandler struct {
	repository port.EventRepository
}

func NewHTTPHandler(repository port.EventRepository) *HTTPHandler {
	return &HTTPHandler{
		repository: repository,
	}
}

func (h *HTTPHandler) GetAllEvents(c *fiber.Ctx) error {
	events := h.repository.GetAll()

	response, err := ToEventsResponse(events)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	fmt.Println("response GetAllEvents", response)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *HTTPHandler) GetEvent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	event := h.repository.GetByID(id)

	response, err := ToEventResponse(id, event)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	fmt.Println("response GetEvent", response)
	return c.Status(fiber.StatusOK).JSON(response)
}
