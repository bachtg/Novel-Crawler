FROM golang:1.21.5
WORKDIR /app
COPY . .
RUN go mod init novel_crawler
RUN go mod tidy
COPY . .
RUN go mod download
RUN go build -o app ./cmd/web

EXPOSE 8080
ENTRYPOINT ["./app"]