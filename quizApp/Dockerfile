#Build Stage - Create a go standalone binary using fat container
FROM golang:1.21-alpine AS builder

RUN apk --update add \
    go \
    musl-dev
RUN apk --update add \
    util-linux-dev
RUN apk add --no-cache tzdata
RUN apk --update --no-cache add curl
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache gcc g++ git openssh-client

COPY go.mod /project/go.mod
COPY go.sum /project/go.sum
COPY pkg /project/pkg

COPY /api/ /project/api

WORKDIR /project/api
RUN CGO_CFLAGS="-D_LARGEFILE64_SOURCE" go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o api .

#Create the actual container no with just the binary
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /project/api
COPY --from=builder /project/api/ /project/api/

EXPOSE 5000

ENTRYPOINT ["/project/api/api"]
