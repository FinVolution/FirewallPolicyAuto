ARG BASE_IMAGE=golang
ARG BASE_IMAGE_VERSION=1.22-alpine

FROM ${BASE_IMAGE}:${BASE_IMAGE_VERSION} AS builder

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk update && \
    apk --no-cache --no-progress add --virtual -U \
    make \
    ca-certificates

# copy src to the proper build path
WORKDIR /build
ADD . .

RUN make dist

# runtime container
FROM alpine:3.15

ENV TZ=Asia/Shanghai

WORKDIR /app/firewall-policy-auto
COPY --from=builder /build/bin/* .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

ENTRYPOINT ["./service_api"]
