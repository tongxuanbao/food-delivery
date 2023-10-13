FROM golang:1.21

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /use

EXPOSE 8080

CMD ["/user"]