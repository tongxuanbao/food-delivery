# --- Build Stage ---
FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /backend

# --- Development Stage ---
FROM cosmtrek/air AS development

WORKDIR /app

EXPOSE 8080

CMD ["air"]

# --- Final Stage ---
FROM alpine:latest

COPY --from=builder /backend /backend

EXPOSE 8080

CMD ["/backend"]