# --- Build Stage ---
FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /order

# --- Final Stage ---
FROM alpine:latest

COPY --from=builder /order /order

EXPOSE 8080

CMD ["/order"]