FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .

# This line installs dependencies and builds the Go application
RUN go mod tidy && go build -o server .

EXPOSE 9000

FROM scratch
COPY --from=builder /app/server /app/server

CMD ["./app/server"]
