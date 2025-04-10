package http

import (
	"time"

	"ApiSmart/src/infrastructure/adapters/http/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Router maneja la configuración de las rutas HTTP
type Router struct {
	authHandler   *handlers.AuthHandler
	sensorHandler *handlers.SensorHandler
	corsConfig    cors.Config
}

// RouterConfig contiene la configuración para el router
type RouterConfig struct {
	AllowedOrigins []string
}

// NewRouter crea un nuevo router HTTP
func NewRouter(
	authHandler *handlers.AuthHandler,
	sensorHandler *handlers.SensorHandler,
	config RouterConfig,
) *Router {
	// Configurar CORS
	corsConfig := cors.Config{
		AllowOrigins:     config.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	return &Router{
		authHandler:   authHandler,
		sensorHandler: sensorHandler,
		corsConfig:    corsConfig,
	}
}

// Setup configura las rutas en el engine de Gin
func (r *Router) Setup() *gin.Engine {
	router := gin.Default()

	// Middleware CORS
	router.Use(cors.New(r.corsConfig))

	// Rutas de autenticación (públicas)
	router.POST("/api/register", r.authHandler.Register)
	router.POST("/api/login", r.authHandler.Login)

	// Ruta para enviar datos de sensores (sin autenticación para dispositivos IoT)
	router.POST("/sensores", r.sensorHandler.CreateSensorData)

	// Rutas protegidas por autenticación
	authorized := router.Group("/api")
	authorized.Use(r.authHandler.AuthMiddleware())
	{
		authorized.GET("/sensors", r.sensorHandler.GetAllSensorData)
		authorized.GET("/sensors/latest", r.sensorHandler.GetLatestSensorData)
		authorized.GET("/sensors/alerts", r.sensorHandler.GetAlerts)
		authorized.PUT("/sensors/alerts/:id/read", r.sensorHandler.MarkAlertAsRead)
	}

	return router
}
