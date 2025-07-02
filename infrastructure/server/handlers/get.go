package handlers

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"net/http"
	"users/infrastructure/server/responses"
)

// Get godoc
// @Summary     List active users
// @Id          Get
// @Produce     json
// @Success     200 {array} responses.UserResponse
// @Failure     500
// @Router      /users [get]
func (h *Handlers) Get(ctx *gin.Context) {
	tracerCtx, span := h.tracer.Start(ctx.Request.Context(), "Handler-Get")
	defer span.End()

	headers := mapToString(ctx.Request.Header)

	result, err := h.actions.Get(tracerCtx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	span.SetAttributes(attribute.String(xAppID, ctx.Request.Header.Get(xAppID)))
	span.SetAttributes(attribute.String("http.headers", headers))

	ctx.JSON(http.StatusOK, gin.H{"data": responses.FromUserList(result)})
}
