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
	"ApiSmart/src/core/application/service"
	"ApiSmart/src/core/application/use_case"
	eventAdapter "ApiSmart/src/infrastructure/adapters/events"
	httpAdapter "ApiSmart/src/infrastructure/adapters/http"
	"ApiSmart/src/infrastructure/adapters/http/handlers"
	"ApiSmart/src/infrastructure/adapters/repositories/mysql"
	"ApiSmart/src/infrastructure/auth"
	"ApiSmart/src/infrastructure/database"
)

func main() {
	// Cargar configuración
	cfg := config.LoadConfig()

	// Inicializar conexión a base de datos
	dbConn, err := database.NewMySQLConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}
	defer dbConn.Close()

	db := dbConn.GetDB()

	// Inicializar servicio JWT
	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpiryHours)

	// Inicializar repositorios
	userRepo := mysql.NewUserRepository(db)
	sensorRepo := mysql.NewSensorRepository(db)

	// Inicializar servicios
	alertService := service.NewAlertService()

	// Intentar inicializar el sistema de eventos
	var eventDispatcher *eventAdapter.EventDispatcherAdapter
	var rabbitMQAdapter *eventAdapter.RabbitMQAdapter

	rabbitMQAdapter, err = eventAdapter.NewRabbitMQAdapter(cfg.RabbitMQ)
	if err != nil {
		log.Printf("Advertencia: Error inicializando RabbitMQ: %v", err)
		log.Println("Continuando sin sistema de eventos...")
	} else {
		defer rabbitMQAdapter.Close()
		eventDispatcher = eventAdapter.NewEventDispatcherAdapter(rabbitMQAdapter)
		log.Println("Sistema de eventos inicializado correctamente")
	}

	// Inicializar casos de uso
	authUseCase := use_case.NewAuthUseCase(userRepo, eventDispatcher, jwtService)
	sensorUseCase := use_case.NewSensorUseCase(sensorRepo, alertService, eventDispatcher)

	// Inicializar handlers HTTP
	authHandler := handlers.NewAuthHandler(authUseCase)
	sensorHandler := handlers.NewSensorHandler(sensorUseCase)

	// Configurar router HTTP
	router := httpAdapter.NewRouter(
		authHandler,
		sensorHandler,
		httpAdapter.RouterConfig{
			AllowedOrigins: []string{"http://localhost:3000", "http://127.0.0.1:8000"},
		},
	)

	// Si el sistema de eventos está disponible, configurar consumidores
	if rabbitMQAdapter != nil {
		// Crear manejadores de eventos
		sensorDataHandler := eventAdapter.NewSensorDataHandler(sensorUseCase)
		alertHandler := eventAdapter.NewAlertHandler(sensorUseCase)
		userEventHandler := eventAdapter.NewUserEventHandler(authUseCase)

		// Suscribir manejadores a topics
		if err := rabbitMQAdapter.Subscribe("sensor.data", sensorDataHandler); err != nil {
			log.Printf("Error suscribiendo al topic sensor.data: %v", err)
		}

		if err := rabbitMQAdapter.Subscribe("sensor.alerts", alertHandler); err != nil {
			log.Printf("Error suscribiendo al topic sensor.alerts: %v", err)
		}

		if err := rabbitMQAdapter.Subscribe("user.events", userEventHandler); err != nil {
			log.Printf("Error suscribiendo al topic user.events: %v", err)
		}
	}

	// Configurar servidor HTTP
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router.Setup(),
	}

	// Iniciar servidor en una goroutine
	go func() {
		log.Printf("Servidor iniciado en el puerto %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error iniciando servidor: %v", err)
		}
	}()

	// Configurar graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Apagando servidor...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Forzando cierre del servidor: %v", err)
	}

	log.Println("Servidor cerrado correctamente")
}
