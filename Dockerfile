FROM golang:1.24.3-alpine3.22

WORKDIR /app

COPY go.mod ./
RUN go mod download

EXPOSE 8080

CMD ["go", "run", "main.go"]

