FROM golang:1.23.4-alpine as build

WORKDIR /app

COPY /root/* /root/

RUN apk update && apk upgrade && apk add git

RUN git config --global url."https://"$( cat /root/token )":x-oauth-basic@github.com/".insteadOf "https://github.com/"

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build main.go

FROM alpine:3.14
WORKDIR /app

COPY --from=build /app/main /app/main
RUN apk add --no-cache tzdata && apk add --no-cache ffmpeg
EXPOSE 8080

CMD [ "./main" ]
