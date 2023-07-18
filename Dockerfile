FROM golang:latest

WORKDIR /app

# download the required Go dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./

EXPOSE 8090

RUN go build -o ./cmd/main ./cmd/main.go
CMD ["./cmd/main"]