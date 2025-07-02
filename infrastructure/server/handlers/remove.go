package handlers

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"net/http"
)

// Remove godoc
// @Summary     Delete a user.
// @Id          Remove
// @Accept      json
// @Produce     json
// @Param       id path string true "User ID"
// @Success     204
// @Failure     400 {object} error "error"
// @Failure     500 {object} error "error"
// @Router      /users/{id} [delete]
func (h *Handlers) Remove(ctx *gin.Context) {
	tracerCtx, span := h.tracer.Start(ctx.Request.Context(), "Handler-Remove")
	defer span.End()

	headers := mapToString(ctx.Request.Header)

	id := ctx.Param("id")

	chErr := make(chan error, 1)
	go func(id string) {
		chErr <- h.actions.Remove(tracerCtx, id)
	}(id)

	if err := <-chErr; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	span.SetAttributes(attribute.String(xAppID, ctx.Request.Header.Get(xAppID)))
	span.SetAttributes(attribute.String("http.headers", headers))
	span.SetAttributes(attribute.String("http.path.id", id))

	ctx.JSON(http.StatusNoContent, gin.H{})
}
