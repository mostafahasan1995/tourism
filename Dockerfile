#first stage - builder
FROM golang:1.22 as builder

WORKDIR /app

ARG GITLAB_ACCESS_TOKEN
RUN git config --global url."https://oauth2:${GITLAB_ACCESS_TOKEN}@git.larsa.io/".insteadOf "https://git.larsa.io/"

COPY go.mod go.sum ./


ENV GO111MODULE=on

RUN CGO_ENABLED=0 GOOS=linux go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build main.go

#second stage
FROM debian:buster-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/.env .
#COPY --from=builder /app/public/ ./public

CMD ["./main"]