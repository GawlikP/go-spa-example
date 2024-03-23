FROM golang:1.22-alpine AS base
WORKDIR /app
COPY go.mod go.sum ./
COPY .env .
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/app /app/cmd/go-spa-example/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/migrate /app/cmd/migrate/main.go
EXPOSE 8080
CMD ["sh", "-c", "/app/migrate && /app/app"]
