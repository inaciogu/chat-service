FROM golang:1.20.3-bullseye

# Set destination for COPY
WORKDIR /go/src

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080
EXPOSE 50051

# Run
CMD ["/chat-service"]
