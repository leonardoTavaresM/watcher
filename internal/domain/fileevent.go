package domain

import "time"

type FileEvent struct {
	Timestamp time.Time
	FilePath  string
	Ext       string
	Event     string
}
