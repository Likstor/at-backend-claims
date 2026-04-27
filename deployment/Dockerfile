FROM golang:1.26-alpine AS builder

WORKDIR /builder

COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify

COPY ./cmd ./cmd
COPY ./internal ./internal

RUN GOOS=linux go build -o /builder/at-migrator ./cmd/migrator/

RUN GOOS=linux go build -o /builder/at-backend ./cmd/app/

FROM alpine:latest

WORKDIR /service

RUN apk add --no-cache curl

COPY --from=builder /builder/at-migrator .
COPY --from=builder /builder/at-backend .

COPY ./resources ./resources
COPY ./migrations ./migrations

COPY ./deployment/entrypoint.sh .

RUN sed -i 's/\r$//' entrypoint.sh && \
    chmod +x entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]