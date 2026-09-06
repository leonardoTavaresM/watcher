package entity

import "time"

type FileEvent struct {
	Timestamp time.Time `json:"timestamp"`
	FilePath  string    `json:"file_path"`
	Ext       string    `json:"extension"`
	Event     string    `json:"event"`
}
