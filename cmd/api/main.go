package main

import (
	"fmt"
	"log"
	"os"

	"github.com/leonardoTavaresM/watcher/internal/adapter/consolepub"
	"github.com/leonardoTavaresM/watcher/internal/adapter/fsnotify"
	"github.com/leonardoTavaresM/watcher/internal/service/watcher"
)

func main() {

	publisher := consolepub.NewConsolePublisher()

	service := watcher.NewWatcherService(publisher)

	adapter := fsnotify.NewFsnotifyAdapter(service)

	path := os.Getenv("WATCH_PATH")
	if path == "" {
		fmt.Println("fallback")
		path = "/app/dev" // fallback
	}

	err := adapter.Start(path)
	if err != nil {
		log.Fatal(err)
	}
}
