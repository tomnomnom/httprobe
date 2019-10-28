FROM golang:1.11.4-alpine3.7 AS build-env
RUN apk add --no-cache --upgrade git openssh-client ca-certificates
RUN go get -u github.com/golang/dep/cmd/dep
WORKDIR /go/src/app

COPY . /go/src/app

RUN go build -o httprobe main.go


FROM alpine:3.9

RUN apk add shadow bash && \
    useradd --create-home --shell /sbin/nologin httprobe && \
    mkdir /httprobe && \
    chown httprobe:httprobe /httprobe

COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-env /go/src/app/httprobe /httprobe/httprobe


USER httprobe
WORKDIR /httprobe

ENTRYPOINT ["/httprobe/httprobe"]
