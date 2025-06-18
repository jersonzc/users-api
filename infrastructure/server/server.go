package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"users/infrastructure/server/handlers"
)

func Setup(config *Config) *http.Server {
	ginServer := gin.New()

	router := ginServer.RouterGroup.Group(config.Prefix)
	router.Any("/health", handlers.HealthCheck)

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      ginServer,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}
}
