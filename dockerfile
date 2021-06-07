FROM golang:1.16.5-alpine3.13
RUN apk add --no-cache libc6-compat
RUN apk add --no-cache iproute2
COPY . ~/
WORKDIR ~/