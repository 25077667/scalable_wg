FROM ghcr.io/linuxserver/wireguard AS base
LABEL maintainer="zxc25077667@pm.me"

VOLUME [ "/lib/modules" ]

ENV PUID=1000
ENV PGID=1000
ENV PEERS=105
ENV TZ=Asia/Taipei
ENV SERVERPORT=51821
ENV PEERDNS=auto
ENV INTERNAL_SUBNET=10.13.13.0
ENV ALLOWEDIPS=0.0.0.0/0

# The monitor
FROM golang:alpine AS builder
RUN apk --no-cache add libcap ca-certificates
COPY . /src
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o monitor monitor.go

# Final stage
FROM base
COPY --from=builder /src/monitor /app
EXPOSE 8080