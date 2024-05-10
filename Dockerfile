# docker build前需要先执行build.bat(或build.sh)脚本进行编译
FROM alpine:3.19.0
RUN apk update --no-cache
RUN apk add bash
WORKDIR /opt/go-admin-server
COPY .build/go-admin-server_linux-amd64 ./go-admin-server_linux-amd64
COPY nacos-config-example.json ./config/nacos-config.json
VOLUME ["./config"]
EXPOSE 8080
ENTRYPOINT ["./go-admin-server_linux-amd64", "server"]
CMD ["-C", "./config/nacos-config.json"]