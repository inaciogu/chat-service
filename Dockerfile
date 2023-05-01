FROM golang:1.20.3-bullseye

# Set destination for COPY
WORKDIR /go/src

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

# Build
RUN go build -o /chat-service

EXPOSE 8080

# Run
CMD ["/chat-service"]
