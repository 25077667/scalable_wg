#!/bin/bash

if [[ $EUID -ne 0 ]]; then
    echo "This script will install packages, please run this in root's permission"
    exit 1
fi

source .env

echo "192 vpn1" >>/etc/iproute2/rt_tables
echo "196 vpn2" >>/etc/iproute2/rt_tables

ip route flush table vpn1
ip route flush table vpn2

ip route add default via ${GATEWAY} enp3s0 src ${VPN1_IP} table vpn1
ip route add default via ${GATEWAY} enp4s0 src ${VPN2_IP} table vpn2

ip rule add from ${VPN1_IP} table vpn1
ip rule add from ${VPN2_IP} table vpn2

echo "
ip route flush table vpn1
ip route flush table vpn2

ip route add default via ${GATEWAY} enp3s0 src ${VPN1_IP} table vpn1
ip route add default via ${GATEWAY} enp4s0 src ${VPN2_IP} table vpn2

ip rule add from ${VPN1_IP} table vpn1
ip rule add from ${VPN2_IP} table vpn2
" >>/etc/rc.local
