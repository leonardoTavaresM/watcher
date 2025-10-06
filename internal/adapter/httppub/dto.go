package httppub

import (
	"time"

	"github.com/leonardoTavaresM/watcher/internal/domain"
)

type FileEvent struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	FilePath  string    `json:"file_path"`
	Ext       string    `json:"ext"`
	Event     string    `json:"event"`
}

type FileEventResponse struct {
	Events map[int]FileEvent `json:"events"`
}

func ToEventsResponse(event map[int]domain.FileEvent) (FileEventResponse, error) {
	response := FileEventResponse{
		Events: make(map[int]FileEvent),
	}

	for i, e := range event {
		response.Events[i] = FileEvent{
			ID:        i,
			Timestamp: e.Timestamp,
			FilePath:  e.FilePath,
			Ext:       e.Ext,
			Event:     e.Event,
		}
	}

	return response, nil
}

func ToEventResponse(id int, event domain.FileEvent) (FileEvent, error) {
	response := FileEvent{
		ID:        id,
		Timestamp: event.Timestamp,
		FilePath:  event.FilePath,
		Ext:       event.Ext,
		Event:     event.Event,
	}
	return response, nil
}
