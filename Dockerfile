FROM golang:alpine

WORKDIR /opt/app

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    GIN_MODE=release
COPY . .
RUN go build -o nacos-sync && chmod 777 ./nacos-sync
CMD ./nacos-sync