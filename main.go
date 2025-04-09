package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ApiSmart/config"
	eventHandlers "ApiSmart/internal/adapters/event"
	"ApiSmart/internal/adapters/handlers"
	"ApiSmart/internal/adapters/repositories/mysql"
	"ApiSmart/internal/core/services"
	"ApiSmart/pkg/database"
	"ApiSmart/pkg/event"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	// Conexión a base de datos
	db, err := database.NewMySQLConnection(cfg.DBConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Inicializar sistema de eventos (RabbitMQ)
	eventDispatcher, broker, err := event.InitEventSystem(cfg.RabbitMQConfig)
	if err != nil {
		log.Printf("Warning: Event system initialization failed: %v", err)
		log.Println("Continuing without event system...")
	} else {
		defer broker.Close()
	}

	// Inicializar repositorios
	userRepo := mysql.NewUserRepository(db)
	sensorRepo := mysql.NewSensorRepository(db)

	// Inicializar servicios
	alertService := services.NewAlertService()

	var authService services.AuthService
	var sensorService services.SensorService

	// Utilizar servicios basados en eventos si el sistema de eventos está disponible
	if eventDispatcher != nil {
		authService = services.NewEventDrivenAuthService(userRepo, eventDispatcher)
		sensorService = services.NewEventDrivenSensorService(sensorRepo, alertService, eventDispatcher)

		// Inicializar handlers de eventos
		sensorDataHandler := eventHandlers.NewSensorDataHandler(sensorService)
		alertHandler := eventHandlers.NewAlertHandler(sensorService)
		userEventHandler := eventHandlers.NewUserEventHandler(authService)

		// Configurar consumidores de eventos
		eventConsumerHandlers := map[string]map[string]event.EventHandler{
			event.TopicSensorData: {
				event.EventTypeSensorDataCreated: sensorDataHandler,
			},
			event.TopicSensorAlerts: {
				event.EventTypeSensorThresholdAlert: alertHandler,
			},
			event.TopicUserEvents: {
				event.EventTypeUserRegistered:    userEventHandler,
				event.EventTypeUserAuthenticated: userEventHandler,
			},
		}

		if err := event.InitConsumers(broker, eventConsumerHandlers); err != nil {
			log.Printf("Warning: Failed to initialize event consumers: %v", err)
		}
	} else {
		// Fallback a servicios tradicionales si el sistema de eventos no está disponible
		authService = services.NewAuthService(userRepo)
		sensorService = services.NewSensorService(sensorRepo, alertService)
	}

	// Inicializar handlers HTTP
	authHandler := handlers.NewAuthHandler(authService)
	sensorHandler := handlers.NewSensorHandler(sensorService)

	// Configurar router
	router := gin.Default()

	// Configurar middleware CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:8000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rutas de autenticación
	router.POST("/api/register", authHandler.Register)
	router.POST("/api/login", authHandler.Login)

	// Ruta para enviar datos de sensores (sin autenticación para facilitar la integración con dispositivos IoT)
	router.POST("/sensores", sensorHandler.CreateSensorData)

	// Rutas que requieren autenticación
	authorized := router.Group("/api")
	{
		authorized.GET("/sensors", sensorHandler.GetAllSensorData)
		authorized.GET("/sensors/latest", sensorHandler.GetLatestSensorData)
		authorized.GET("/sensors/alerts", sensorHandler.GetAlerts)
		authorized.PUT("/sensors/alerts/:id/read", sensorHandler.MarkAlertAsRead)
	}

	// Configurar servidor HTTP
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	// Iniciar servidor en una goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.ServerPort)
	log.Printf("Event system initialized: %v", eventDispatcher != nil)

	// Configurar graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
