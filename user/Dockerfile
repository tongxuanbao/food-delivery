# --- Build Stage ---
FROM golang:alpine3.18 AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=1 GOOS=linux go build -o /user

# --- Development Stage ---
FROM cosmtrek/air AS development

WORKDIR /app

EXPOSE 8080

CMD ["air"]

# --- Final Stage ---
FROM alpine:latest

RUN apk add --no-cache libc6-compat

COPY --from=builder /user /user

EXPOSE 8080

CMD ["/user"]
