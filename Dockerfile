FROM golang:1.9 as builder
ARG VERSION
WORKDIR /go/src/github.com/go-ignite/ignite
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o ignite .


FROM alpine
LABEL maintainer="go-ignite"
RUN apk --no-cache add ca-certificates tzdata sqlite \
			&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
			&& echo "Asia/Shanghai" >  /etc/timezone \
			&& apk del tzdata
# See https://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
VOLUME /root/ignite/data

WORKDIR /root/ignite
COPY --from=builder /go/src/github.com/go-ignite/ignite/ignite ./
COPY --from=builder /go/src/github.com/go-ignite/ignite/templates ./templates
COPY --from=builder /go/src/github.com/go-ignite/ignite/static ./static
COPY --from=builder /go/src/github.com/go-ignite/ignite/conf ./conf
RUN mv ./conf/config-temp.toml ./conf/config.toml

EXPOSE 5000
ENTRYPOINT ["./ignite"]
