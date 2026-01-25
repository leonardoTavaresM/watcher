package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/leonardoTavaresM/watcher/internal/application/service"
	"github.com/leonardoTavaresM/watcher/internal/infrastructure/adapter/fsnotify"
	"github.com/leonardoTavaresM/watcher/internal/infrastructure/adapter/publisher"
	"github.com/leonardoTavaresM/watcher/internal/infrastructure/adapter/repository"
	"github.com/leonardoTavaresM/watcher/internal/infrastructure/config"
	"github.com/leonardoTavaresM/watcher/internal/infrastructure/handler"
)

func main() {
	app := fiber.New()

	// Infrastructure - Repository
	repo := repository.NewInMemoryRepository()

	// Infrastructure - Publishers
	consolePublisher := publisher.NewConsolePublisher(repo)

	rabbitConfig := config.GetRabbitMQConfig()
	rabbitPublisher, err := publisher.NewRabbitMQPublisher(
		rabbitConfig.URI,
		rabbitConfig.Exchange,
		rabbitConfig.Queue,
	)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ publisher: %v", err)
	}
	defer rabbitPublisher.Close()

	// Application - Service
	watcherService := service.NewWatcherService(repo, consolePublisher, rabbitPublisher)

	// Infrastructure - Adapters
	fsnotifyAdapter := fsnotify.NewFsnotifyAdapter(watcherService)

	// Infrastructure - HTTP Handlers
	httpHandler := handler.NewHTTPHandler(repo)

	// Routes
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(`{pong}`)
	})
	app.Get("/events", httpHandler.GetAllEvents)
	app.Get("/events/:id", httpHandler.GetEvent)

	// Watch path
	path := os.Getenv("WATCH_PATH")
	if path == "" {
		fmt.Println("fallback to current directory")
		path = "."
	}

	// Signal handling
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Start HTTP server
	go func() {
		err := app.Listen(":3000")
		if err != nil {
			log.Fatal("HTTP server error:", err)
		}
	}()

	// Start fsnotify watcher
	go func() {
		err := fsnotifyAdapter.Start(path)
		if err != nil {
			log.Fatal("fsnotify error:", err)
		}
	}()

	// Wait for shutdown signal
	<-c
	log.Println("Shutting down...")

	app.Shutdown()

	log.Println("Server stopped")
}
