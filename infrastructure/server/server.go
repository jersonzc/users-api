package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"net/http"
	"users/infrastructure/dependencies"
	"users/infrastructure/server/handlers"
	"users/infrastructure/server/routes"
)

func Setup(config *Config, actions *dependencies.Actions) *http.Server {
	ginServer := gin.New()
	ginServer.Use(otelgin.Middleware("app-server-gin"))

	router := ginServer.Group(config.Prefix)
	router.GET("/health", handlers.HealthCheck)
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.Setup(router, actions)

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      ginServer,
		IdleTimeout:  config.IdleTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}
}
