#!/bin/bash

if [[ $EUID -ne 0 ]]; then
    echo "This script will install packages, please run this in root's permission"
    exit 1
fi

echo "192 vpn1" >>/etc/iproute2/rt_tables
echo "196 vpn2" >>/etc/iproute2/rt_tables

ip route flush table vpn1
ip route flush table vpn2

ip route add default via 140.117.169.254 enp3s0 src 140.117.169.212 table vpn1
ip route add default via 140.117.169.254 enp4s0 src 140.117.169.213 table vpn2

ip rule add from 140.117.169.212 table vpn1
ip rule add from 140.117.169.213 table vpn2

echo "
ip route flush table vpn1
ip route flush table vpn2

ip route add default via 140.117.169.254 enp3s0 src 140.117.169.212 table vpn1
ip route add default via 140.117.169.254 enp4s0 src 140.117.169.213 table vpn2

ip rule add from 140.117.169.212 table vpn1
ip rule add from 140.117.169.213 table vpn2
" >>/etc/rc.local
