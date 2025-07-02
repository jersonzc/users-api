package handlers

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"net/http"
	"users/infrastructure/server/responses"
)

// GetSingle godoc
// @Summary     List a single user
// @Id          GetSingle
// @Produce     json
// @Param       id path string true "User ID"
// @Success     200 {object} responses.UserResponse
// @Failure     400 {object} error "error"
// @Failure     500 {object} error "error"
// @Router      /users/search/{id} [get]
func (h *Handlers) GetSingle(ctx *gin.Context) {
	tracerCtx, span := h.tracer.Start(ctx.Request.Context(), "Handler-GetSingle")
	defer span.End()

	headers := mapToString(ctx.Request.Header)

	id := ctx.Param("id")

	result, err := h.actions.GetByID(tracerCtx, []string{id})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	span.SetAttributes(attribute.String(xAppID, ctx.Request.Header.Get(xAppID)))
	span.SetAttributes(attribute.String("http.headers", headers))
	span.SetAttributes(attribute.String("http.path.id", id))

	ctx.JSON(http.StatusOK, gin.H{"data": responses.FromUserList(result)})
}
