FROM golang:1.23 AS build

WORKDIR /cli


COPY go.mod .
COPY go.sum .

RUN go mod download 

COPY cmd cmd
COPY main.go .

RUN go build

FROM alpinelinux/docker-cli

RUN apk add gcompat

COPY --from=build /cli/gone-cli /gone-cli

CMD ["/gone-cli"]
