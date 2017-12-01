FROM golang:1.9 as builder
WORKDIR /go/src/github.com/go-ignite/ignite
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ignite .


FROM alpine
LABEL maintainer="iwendellsun@gmail.com"
RUN apk --no-cache add ca-certificates tzdata \
			&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
			&& echo "Asia/Shanghai" >  /etc/timezone \
			&& apk del tzdata

WORKDIR /root/ignite
COPY --from=builder /go/src/github.com/go-ignite/ignite/ignite ./
COPY --from=builder /go/src/github.com/go-ignite/ignite/templates ./templates
COPY --from=builder /go/src/github.com/go-ignite/ignite/static ./static

EXPOSE 5000
CMD ["/bin/sh", "-c", "./ignite"]
