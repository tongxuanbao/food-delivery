# --- Build Stage ---
FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /restaurant

# --- Final Stage ---
FROM alpine:latest

COPY --from=builder /restaurant /restaurant

EXPOSE 8080

CMD ["/restaurant"]