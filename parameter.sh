#!/bin/bash
sysctl -w net.ipv4.ip_local_port_range="1024 65535"
sysctl -w net.ipv4.tcp_tw_reuse="1"
sysctl -w net.ipv4.tcp_rmem="1024 2048 4096"
sysctl -w net.ipv4.tcp_wmem="1024 2048 4096"

sysctl -w net.core.netdev_max_backlog=30000
sysctl -w net.ipv4.tcp_max_syn_backlog=2048
sysctl -w net.core.somaxconn="10000"

ifconfig lo txqueuelen 5000
