FROM golang:1.24-alpine

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
COPY .env .env

RUN go build -o main main.go

EXPOSE 8080

ENTRYPOINT ["/app/main"]
