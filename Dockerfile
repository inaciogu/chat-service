FROM golang:1.20.3-alpine3.16 as base
WORKDIR /src/chatservice
COPY go.mod go.sum ./
COPY . .
RUN go build -o chatservice ./cmd/chatservice

FROM alpine:3.16 as binary
COPY --from=base /src/chatservice/chatservice .
EXPOSE 8080
EXPOSE 50051
CMD ["./chatservice"]