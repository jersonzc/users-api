package routes

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"users/infrastructure/dependencies"
	"users/infrastructure/server/handlers"
)

func Setup(baseRouter *gin.RouterGroup, actions *dependencies.Actions, tracer trace.Tracer) *gin.RouterGroup {
	handler := handlers.New(actions, tracer)

	prefix := baseRouter.Group("/users")

	prefix.GET("", handler.Get)
	prefix.POST("", handler.Save)
	prefix.PUT(":id", handler.Update)
	prefix.DELETE(":id", handler.Remove)

	prefix.POST("/search", handler.GetMultiple)
	prefix.GET("/search/:id", handler.GetSingle)

	return prefix
}
