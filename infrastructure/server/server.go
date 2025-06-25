package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"users/infrastructure/dependencies"
	"users/infrastructure/server/handlers"
)

func Setup(config *Config, actions *dependencies.Actions, tracer *trace.Tracer) *http.Server {
	ginServer := gin.New()

	ginServer.Use(otelgin.Middleware("app-server-gin"))

	router := ginServer.RouterGroup.Group(config.Prefix)
	router.Any("/health", handlers.HealthCheck)

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      ginServer,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}
}
