package handlers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"io"
	"net/http"
	"users/domain/entities"
	"users/infrastructure/server/requests"
	"users/infrastructure/server/responses"
)

// Save godoc
// @Summary     Create a user
// @Id          Save
// @Accept      json
// @Produce     json
// @Param       payload body requests.SaveUser true "Create a user: 'name' field is required; all other fields are optional."
// @Success     201 {object} responses.UserResponse
// @Failure     400 {object} error "error"
// @Failure     500 {object} error "error"
// @Router      /users [post]
func (h *Handlers) Save(ctx *gin.Context) {
	tracerCtx, span := h.tracer.Start(ctx.Request.Context(), "Handler-Save")
	defer span.End()

	headers := mapToString(ctx.Request.Header)

	data, err := ctx.GetRawData()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	ctx.Request.Body = io.NopCloser(bytes.NewReader(data))

	var body requests.SaveUser
	if bindErr := ctx.ShouldBindJSON(&body); bindErr != nil {
		span.RecordError(bindErr)
		span.SetStatus(codes.Error, bindErr.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": bindErr.Error()})
		return
	}

	user, err := body.ToUser()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	type result struct {
		user *entities.User
		err  error
	}

	resultChan := make(chan result, 1)
	go func(user *entities.User) {
		var r result
		r.user, r.err = h.actions.Save(tracerCtx, user)
		resultChan <- r
	}(user)

	r := <-resultChan

	if r.err != nil {
		span.RecordError(r.err)
		span.SetStatus(codes.Error, r.err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": r.err.Error()})
		return
	}

	span.SetAttributes(attribute.String(xAppID, ctx.Request.Header.Get(xAppID)))
	span.SetAttributes(attribute.String("http.headers", headers))
	span.SetAttributes(attribute.String("http.body", string(data)))

	ctx.JSON(http.StatusCreated, gin.H{"data": responses.FromUser(r.user)})
}
