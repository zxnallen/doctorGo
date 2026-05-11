FROM golang:1.25-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/doctor-go ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/doctor-go-migrate ./cmd/migrate

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/doctor-go /app/doctor-go
COPY --from=builder /out/doctor-go-migrate /app/doctor-go-migrate
COPY configs /app/configs
COPY docs/swagger /app/docs/swagger

EXPOSE 8080

ENV APP_ENV=prod

CMD ["/app/doctor-go"]
