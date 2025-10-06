package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/leonardoTavaresM/watcher/internal/adapter/consolepub"
	"github.com/leonardoTavaresM/watcher/internal/adapter/fsnotify"
	"github.com/leonardoTavaresM/watcher/internal/adapter/httppub"
	"github.com/leonardoTavaresM/watcher/internal/domain/repository/memory"
	"github.com/leonardoTavaresM/watcher/internal/domain/service/watcher"
)

func main() {

	app := fiber.New()

	repository := memory.NewInMemoryEvent()
	publisher := consolepub.NewConsolePublisher(repository)
	service := watcher.NewWatcherService(repository, publisher)

	adapter := fsnotify.NewFsnotifyAdapter(service)

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(`{pong}`)
	})

	path := os.Getenv("WATCH_PATH")
	if path == "" {
		fmt.Println("fallback")
		path = "/app/dev"
	}

	httpAdapter := httppub.NewHTTPAdapter(repository)

	app.Get("/events", httpAdapter.GetAllEvents)
	app.Get("/events/:id", httpAdapter.GetEvent)

	// Canal para sinais do sistema
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Iniciar HTTP server em goroutine
	go func() {
		err := app.Listen(":3000")
		if err != nil {
			log.Fatal("HTTP server error:", err)
		}
	}()

	// Iniciar fsnotify em goroutine
	go func() {
		err := adapter.Start(path)
		if err != nil {
			log.Fatal("fsnotify error:", err)
		}
	}()

	// Aguardar sinal de shutdown
	<-c
	log.Println("Shutting down...")

	// Shutdown graceful do HTTP server
	app.Shutdown()

	log.Println("Server stopped")
}
