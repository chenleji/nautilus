FROM golang:1.10 as builder
MAINTAINER chenleji@gmail.com

COPY ./ /go/src/github.com/chenleji/nautilus/
WORKDIR /go/src/github.com/chenleji/nautilus/
RUN go build -o nautilus

FROM centos:7
MAINTAINER chenleji@gmail.com
RUN mkdir -p /conf && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
WORKDIR /
COPY --from=builder /go/src/github.com/chenleji/nautilus/nautilus .
COPY --from=builder /go/src/github.com/chenleji/nautilus/conf/app.conf ./conf/

CMD /nautilus