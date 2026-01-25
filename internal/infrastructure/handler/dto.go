package handler

import (
	"time"

	"github.com/leonardoTavaresM/watcher/internal/domain/entity"
)

type FileEventDTO struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	FilePath  string    `json:"file_path"`
	Ext       string    `json:"ext"`
	Event     string    `json:"event"`
}

type FileEventResponse struct {
	Events map[int]FileEventDTO `json:"events"`
}

func ToEventsResponse(events map[int]entity.FileEvent) (FileEventResponse, error) {
	response := FileEventResponse{
		Events: make(map[int]FileEventDTO),
	}

	for i, e := range events {
		response.Events[i] = FileEventDTO{
			ID:        i,
			Timestamp: e.Timestamp,
			FilePath:  e.FilePath,
			Ext:       e.Ext,
			Event:     e.Event,
		}
	}

	return response, nil
}

func ToEventResponse(id int, event entity.FileEvent) (FileEventDTO, error) {
	response := FileEventDTO{
		ID:        id,
		Timestamp: event.Timestamp,
		FilePath:  event.FilePath,
		Ext:       event.Ext,
		Event:     event.Event,
	}
	return response, nil
}
