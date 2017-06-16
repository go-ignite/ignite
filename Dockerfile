FROM alpine
MAINTAINER Wendell Sun <iwendellsun@gmail.com>

WORKDIR /ignite

ADD ignite /ignite/ignite
ADD templates /ignite/templates
ADD static /ignite/static

EXPOSE 5000
CMD ["/bin/sh", "-c", "./ignite"]
