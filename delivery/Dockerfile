# --- Build Stage ---
FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /delivery

# --- Final Stage ---
FROM alpine:latest

COPY --from=builder /delivery /delivery

EXPOSE 8080

CMD ["/delivery"]