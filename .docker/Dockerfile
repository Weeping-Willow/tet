FROM golang:1.24.9-alpine AS builder

WORKDIR app

COPY ../go.mod go.sum ./
RUN go mod download
COPY .. .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app ./cmd/app

FROM alpine:3

COPY --from=BUILDER /app /app