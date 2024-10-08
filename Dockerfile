FROM golang:1.22-alpine AS base
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache nodejs-current=21.7.3-r0
RUN apk add --no-cache yarn
WORKDIR /app
COPY go.mod go.sum ./
COPY .env .
COPY . .
RUN go mod download
RUN cd /app/ui && yarn --force && yarn run build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/app /app/cmd/go-spa-example/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/migrate /app/cmd/migrate/main.go
EXPOSE 8080
RUN chmod +x /app/app
CMD ["sh", "-c", "/app/migrate && go test -v -p 1 ./... && /app/app"]
