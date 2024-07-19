package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"message_service/internal/caching"
	"message_service/internal/config"
	"message_service/internal/db"
	"message_service/internal/events"
	"message_service/internal/handlers"
	"message_service/internal/messaging"
	"message_service/internal/repositories"
	"message_service/internal/services"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	app := &AppDIContainer{}
	cfg := config.Load()

	bootstrapServer(app, cfg)
}

func startServer(app *AppDIContainer, srvPort string) {
	log.Printf("Server listening on port %s", srvPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", srvPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func startEventConsumersInBackground(app *AppDIContainer, rabbitMQ *messaging.RabbitMQConn) {
	fmt.Println("Starting event consumers...")
	consumerManager := messaging.NewConsumerManager(rabbitMQ)

	eventQueue := events.EventQueues[events.MessageCreationRequestedQueue]
	consumerManager.AddConsumer(eventQueue.Name, app.EventHandlers.HandleMessageCreated)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := consumerManager.StartConsumers(ctx)
	if err != nil {
		log.Fatalf("Failed to start consumers: %v", err)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down...")
	cancel()
	time.Sleep(2 * time.Second)
}

func bootstrapServer(app *AppDIContainer, cfg *config.Config) {
	dbConn, err := db.NewMySQLConn(cfg.MySQLDSN)
	if err != nil {
		log.Panic(fmt.Errorf("failed to connect to MySQL: %v", err))
	}

	// repositories
	messageRepository := repositories.NewMySQLMessageRepository(dbConn)
	chatRepository := repositories.NewMySQLChatRepository(dbConn)

	// messaging
	rabbitMQ, err := messaging.NewRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		log.Panic(fmt.Errorf("failed to connect to RabbitMQ: %v", err))
	}

	// services
	messageService := &services.MessageService{
		MessageRepository:     messageRepository,
		ChatRepository:        chatRepository,
		EventPublisherManager: rabbitMQ,
	}

	// event handlers
	eventHandler := handlers.NewEventHandler(messageService)
	app.EventHandlers = eventHandler

	// caching
	cache := caching.NewRedisCache(cfg.RedisURL, "", 0)

	// API handlers
	app.ApiCmdHandlers = &handlers.ApiCmdHandlers{
		MessageService:        messageService,
		Cache:                 cache,
		EventPublisherManager: rabbitMQ,
	}

	go startServer(app, cfg.ServerPort)
	startEventConsumersInBackground(app, rabbitMQ)
}
