FROM golang:1.24.3-alpine3.22

WORKDIR /app

COPY main.go ./
COPY go.mod ./
COPY rss_urls.json ./

EXPOSE 8080

CMD ["go", "run", "main.go"]

