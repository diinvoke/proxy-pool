FROM golang:1.13.5 AS builder
LABEL maintainer="Ming Cheng"

# Using 163 mirror for Debian Strech
#RUN sed -i 's/deb.debian.org/mirrors.163.com/g' /etc/apt/sources.list
#RUN apt-get update

ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV PACKAGE github.com/mingcheng/proxypool
ENV BUILD_DIR ${GOPATH}/src/${PACKAGE}
ENV GOPROXY "https://goproxy.cn"

# Print go version
#RUN echo "GOROOT is ${GOROOT}"
#RUN echo "GOPATH is ${GOPATH}"
#RUN ${GOROOT}/bin/go version

# Build
COPY . ${BUILD_DIR}
WORKDIR ${BUILD_DIR}
RUN make clean && make build && mv ${BUILD_DIR}/proxypool /usr/bin/proxypool

# Stage2
FROM alpine:3.10.3

# @from https://mirrors.ustc.edu.cn/help/alpine.html
#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

COPY --from=builder /usr/bin/proxypool /bin/proxypool
ENV PORT 8080
EXPOSE ${PORT}
HEALTHCHECK --interval=5s --timeout=3s \
              CMD curl -fs http://localhost:${PORT}/ || exit 1
ENTRYPOINT ["/bin/proxypool"]