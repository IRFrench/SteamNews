# Create the website as a binary
FROM golang:1.23-alpine AS binary

RUN apk update; apk add git make

COPY . /opt
WORKDIR /opt


RUN CGO_ENABLED=0 go build -o ./build/steamnews ./cmd/main.go

# Create the container
FROM alpine:3.20.2

RUN apk update; apk add wget curl

RUN update-ca-certificates --fresh

COPY --from=binary /opt/build/steamnews /usr/bin/steamnews

WORKDIR /etc

VOLUME /etc

CMD ["steamnews"]