#!/bin/bash

# update
apt update && apt upgrade -y

# essential pkg
apt install build-essentail nginx git openssh

# install golang
wget https://golang.org/dl/go1.16.5.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.5.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> /etc/profile

# golang lis
go get -u github.com/gofiber/fiber/v2
