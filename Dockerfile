#FROM golang:1.13.15-alpine3.12
FROM alpine:3.12

EXPOSE 8080

COPY etc /opt/pitbull/etc
COPY pitbull /opt/pitbull/pitbull

WORKDIR /opt/pitbull

CMD ["/opt/pitbull/pitbull"]