package handlers

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"net/http"
	"users/infrastructure/server/requests"
	"users/infrastructure/server/responses"
)

// GetMultiple godoc
// @Summary     List multiple users
// @Id          GetMultiple
// @Accept      json
// @Produce     json
// @Param       request body requests.MultipleIDRequest true "Enter the IDs of the users to list."
// @Success     200 {array} responses.UserResponse
// @Failure     400 {object} error "error"
// @Failure     500 {object} error "error"
// @Router      /users/search [post]
func (h *Handlers) GetMultiple(ctx *gin.Context) {
	tracerCtx, span := h.tracer.Start(ctx.Request.Context(), "Handler-GetMultiple")
	defer span.End()

	headers := mapToString(ctx.Request.Header)

	var body requests.MultipleIDRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	result, err := h.actions.GetByID(tracerCtx, body.Users)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	span.SetAttributes(attribute.String(xAppID, ctx.Request.Header.Get(xAppID)))
	span.SetAttributes(attribute.String("http.headers", headers))
	span.SetAttributes(attribute.StringSlice("http.body.users", body.Users))

	ctx.JSON(http.StatusOK, gin.H{"data": responses.FromUserList(result)})
}
