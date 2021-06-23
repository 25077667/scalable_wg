FROM nginx:alpine AS builder
LABEL maintainer="zxc25077667@pm.me"
RUN apk --no-cache add libcap ca-certificates
COPY html/ /usr/share/nginx/html

EXPOSE 80
EXPOSE 443