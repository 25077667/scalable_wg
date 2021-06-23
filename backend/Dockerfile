FROM golang:alpine AS builder
LABEL maintainer="zxc25077667@pm.me"
RUN apk --no-cache add libcap ca-certificates
COPY . /src
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o manager manager.go

# Final stage
FROM alpine
WORKDIR /app
COPY --from=builder /src/manager /app
EXPOSE 3000
CMD ["./manager"]