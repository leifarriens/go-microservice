FROM golang:1.22.1-alpine AS build-base

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go mod download

FROM build-base AS dev

RUN go install github.com/cosmtrek/air@v1.43.0

COPY . .

CMD ["air", "-c", ".air.toml"]

FROM build-base AS build-prod

COPY . .

RUN go run github.com/swaggo/swag/cmd/swag@v1.16.4 init \
  --dir ./cmd,./model/,./handler/ \
  --generalInfo main.go \
  --requiredByDefault \
  --outputTypes yaml,go

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-s' -x -o /service ./cmd

FROM gcr.io/distroless/base-debian11 AS prod

WORKDIR /

COPY --from=build-prod ./service ./service

ENV ENVIRONMENT=container

EXPOSE 8080

CMD ["./service"]