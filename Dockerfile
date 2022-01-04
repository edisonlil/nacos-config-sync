FROM golang:alpine

WORKDIR /opt/app

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    GIN_MODE=release \
    project_addr=/opt/config \
    config_addr=/opt/config \
    nacos_username=nacos \
    nacos_password=nacos

COPY . .
COPY ./sync-config.yml /etc/sync-config/sync-config.yml

RUN go build -o nacos-sync && chmod 777 ./nacos-sync
CMD ./nacos-sync