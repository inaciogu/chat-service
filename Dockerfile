FROM golang:1.20.3-bullseye

ENV PATH="root/.cargo/bin:${PATH}"
ENV USER=root

WORKDIR /go/src

RUN ln -sf /bin/bash /bin/sh

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

CMD [ "go", "run", "cmd/chat-service/main.go" ]