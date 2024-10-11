FROM golang:1.22.1-alpine AS base

FROM base AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-s' -x -o /service ./cmd

FROM gcr.io/distroless/base-debian11 AS runner

WORKDIR /

COPY --from=builder ./service ./service

ENV ENVIRONMENT=container

EXPOSE 8080

CMD ["./service"]