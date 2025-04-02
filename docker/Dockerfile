FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-s -w' -o main cmd/api/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/main /app/main

ENV config=docker

EXPOSE 8080

CMD ["/app/main"]