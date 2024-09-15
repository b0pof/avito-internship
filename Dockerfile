FROM golang:alpine as build

COPY . /project

WORKDIR /project

RUN apk add make && go build -o ./bin/main ./cmd/main.go

#========================================

FROM alpine:latest

COPY --from=build /project/bin/main /bin/

RUN apk update && apk add bash

WORKDIR /project

RUN apk add --no-cache tzdata
ENV TZ="Europe/Moscow"

CMD ["main"]