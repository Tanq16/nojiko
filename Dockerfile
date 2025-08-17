FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -ldflags="-s -w" -o nojiko cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/nojiko .
EXPOSE 8080

CMD ["/app/nojiko"]
