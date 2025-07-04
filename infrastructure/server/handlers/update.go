package handlers

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"io"
	"net/http"
	"users/domain/entities"
	"users/infrastructure/server/requests"
	"users/infrastructure/server/responses"
)

// Update godoc
// @Summary     Modify a user
// @Id          Update
// @Accept      json
// @Produce     json
// @Param       id path string true "The ID of the user."
// @Param       request body requests.UpdateUser true "The info to update."
// @Success     204
// @Failure     400 {object} error "error"
// @Failure     500 {object} error "error"
// @Router      /users/{id} [put]
func (h *Handlers) Update(ctx *gin.Context) {
	tracerCtx, span := h.tracer.Start(ctx.Request.Context(), "Handler-Update")
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

	var body requests.UpdateUser
	if err = ctx.ShouldBindJSON(&body); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	id := ctx.Param("id")

	fields, err := body.ToMap()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	if len(fields) == 0 {
		err = errors.New("at least one field is required")
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
	go func(id string, fields map[string]interface{}) {
		var r result
		r.user, r.err = h.actions.Update(tracerCtx, id, fields)
		resultChan <- r
	}(id, fields)

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
	span.SetAttributes(attribute.String("http.path.id", id))

	ctx.JSON(http.StatusOK, gin.H{"data": responses.FromUser(r.user)})
}
