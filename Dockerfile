FROM golang:1.23.4-alpine AS build

WORKDIR /app

COPY . .

RUN apk update && apk upgrade git

COPY go.mod ./
COPY go.sum ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service

FROM alpine:latest
WORKDIR /app
RUN apk update && apk upgrade && apk add git tzdata ffmpeg

COPY --from=build /app/service /app/service
EXPOSE 8080

ENTRYPOINT [ "./service" ]
