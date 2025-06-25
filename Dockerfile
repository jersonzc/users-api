# syntax=docker/dockerfile:1

FROM golang:1.23.9-alpine AS builder

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app

FROM gcr.io/distroless/base-debian12 AS runner

WORKDIR /

COPY --from=builder /app /app
COPY --from=builder /go/src/migrations /migrations

USER nonroot:nonroot

ENTRYPOINT ["/app"]