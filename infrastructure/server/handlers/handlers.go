package handlers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"io"
	"net/http"
	"strings"
	"users/domain/entities"
	"users/infrastructure/dependencies"
	"users/infrastructure/server/requests"
	"users/infrastructure/server/responses"
)

const xAppID = "X-Application-ID"

type Handlers struct {
	actions *dependencies.Actions
	tracer  trace.Tracer
}

func New(actions *dependencies.Actions, tracer trace.Tracer) *Handlers {
	return &Handlers{actions: actions, tracer: tracer}
}

// Get godoc
// @Summary     List users
// @Id          Get
// @Produce     json
// @Success     200 {array} responses.UserResponse
// @Failure     500
// @Router      /users [get]
func (h *Handlers) Get(ctx *gin.Context) {
	tracerCtx, span := h.tracer.Start(ctx.Request.Context(), "Get")
	defer span.End()

	headers := mapToString(ctx.Request.Header)

	result, err := h.actions.Get(tracerCtx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	span.SetAttributes(attribute.String(xAppID, ctx.Request.Header.Get(xAppID)))
	span.SetAttributes(attribute.String("http.headers", headers))

	ctx.JSON(http.StatusOK, gin.H{"data": responses.FromUserList(result)})
}

// GetSingle godoc
// @Summary     List a single user
// @Id          GetSingle
// @Produce     json
// @Param       id path int true "User ID"
// @Success     200 {object} responses.UserResponse
// @Failure     400 {object} error "error"
// @Failure     500 {object} error "error"
// @Router      /users/search/{id} [get]
func (h *Handlers) GetSingle(ctx *gin.Context) {
	tracerCtx, span := h.tracer.Start(ctx.Request.Context(), "GetSingle")
	defer span.End()

	headers := mapToString(ctx.Request.Header)

	id := ctx.Param("id")

	result, err := h.actions.GetByID(tracerCtx, []string{id})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	span.SetAttributes(attribute.String(xAppID, ctx.Request.Header.Get(xAppID)))
	span.SetAttributes(attribute.String("http.headers", headers))
	span.SetAttributes(attribute.String("http.path.id", id))

	ctx.JSON(http.StatusOK, gin.H{"data": responses.FromUserList(result)})
}

// GetMultiple godoc
// @Summary     List multiple users
// @Id          GetMultiple
// @Accept      json
// @Produce     json
// @Param       request body requests.MultipleIDRequest true "list might be empty"
// @Success     200 {array} responses.UserResponse
// @Failure     400 {object} error "error"
// @Failure     500 {object} error "error"
// @Router      /users/search [post]
func (h *Handlers) GetMultiple(ctx *gin.Context) {
	tracerCtx, span := h.tracer.Start(ctx.Request.Context(), "GetMultiple")
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	span.SetAttributes(attribute.String(xAppID, ctx.Request.Header.Get(xAppID)))
	span.SetAttributes(attribute.String("http.headers", headers))
	span.SetAttributes(attribute.StringSlice("http.body.users", body.Users))

	ctx.JSON(http.StatusOK, gin.H{"data": responses.FromUserList(result)})
}

// Save godoc
// @Summary     Create a user
// @Id          Save
// @Accept      json
// @Produce     json
// @Param       request body requests.SaveUser true
// @Success     200 {object} responses.UserResponse
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
	go func(u *entities.User) {
		var r result
		r.user, r.err = h.actions.Save(tracerCtx, u)
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

// Update godoc
// @Summary     Modify a user
// @Id          Update
// @Accept      json
// @Produce     json
// @Param       id path int true "User ID"
// @Param       request body requests.UpdateUser true
// @Success     200 {object} responses.UserResponse
// @Failure     400 {object} error "error"
// @Failure     500 {object} error "error"
// @Router      /users/{id} [put]
func (h *Handlers) Update(ctx *gin.Context) {
	tracerCtx, span := h.tracer.Start(ctx.Request.Context(), "Update")
	defer span.End()

	headers := mapToString(ctx.Request.Header)

	data, err := ctx.GetRawData()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	id := ctx.Param("id")

	ctx.Request.Body = io.NopCloser(bytes.NewReader(data))

	var body requests.UpdateUser
	if err := ctx.ShouldBindJSON(&body); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	user, err := body.ToUser()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	user.ID = id

	type result struct {
		user *entities.User
		err  error
	}

	resultChan := make(chan result, 1)
	go func(u *entities.User) {
		var r result
		r.user, r.err = h.actions.Update(tracerCtx, u)
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
	span.SetAttributes(attribute.String("http.path.id", id))

	ctx.JSON(http.StatusOK, responses.FromUser(r.user))
}

// Remove godoc
// @Summary     Delete a user.
// @Id          Remove
// @Accept      json
// @Produce     json
// @Param       id path int true "User ID"
// @Success     200 {object} interface{} "empty response"
// @Failure     400 {object} error "error"
// @Failure     500 {object} error "error"
// @Router      /users/{id} [delete]
func (h *Handlers) Remove(ctx *gin.Context) {
	tracerCtx, span := h.tracer.Start(ctx.Request.Context(), "Remove")
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

	ctx.JSON(http.StatusOK, render.Data{})
}

func mapToString(myMap map[string][]string) string {
	var sb strings.Builder

	for key, values := range myMap {
		for _, v := range values {
			sb.WriteString(fmt.Sprintf("%q:%q", key, v))
			sb.WriteString(" ")
		}
	}

	result := sb.String()
	if len(result) > 0 {
		result = result[:len(result)-1]
	}

	return result
}
