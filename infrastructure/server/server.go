package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"net/http"
	"users/infrastructure/dependencies"
	"users/infrastructure/server/handlers"
	"users/infrastructure/server/routes"
)

func Setup(config *Config, actions *dependencies.Actions) *http.Server {
	ginServer := gin.New()
	ginServer.Use(otelgin.Middleware("app-server-gin"))

	router := ginServer.RouterGroup.Group(config.Prefix)
	router.Any("/health", handlers.HealthCheck)

	routes.Setup(router, actions)

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      ginServer,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}
}
