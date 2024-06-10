FROM golang:1.21.5
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o app ./cmd/web
EXPOSE 8080
ENTRYPOINT ["./app"]