FROM golang:1.23.4-alpine as build

WORKDIR /app

COPY . .

RUN apk update && apk upgrade && apk add git tzdata ffmpeg

COPY go.mod ./
COPY go.sum ./
RUN go mod download


RUN go build main.go

FROM alpine:3.14
WORKDIR /app

COPY --from=build /app/main /app/main
EXPOSE 8080

CMD [ "./main" ]
