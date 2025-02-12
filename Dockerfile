FROM golang:1.23.4-alpine as build

WORKDIR /app

COPY . .

RUN apk update && apk upgrade && apk add git tzdata ffmpeg

COPY go.mod ./
COPY go.sum ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go -o service

FROM alpine:3.14
WORKDIR /app

COPY --from=build /app/service /app/service
EXPOSE 8080

CMD [ "./service" ]
