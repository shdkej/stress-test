#!/bin/bash
sysctl -w net.ipv4.ip_local_port_range="1024 61000"
sysctl -w net.ipv4.tcp_tw_reuse="0"
sysctl -w net.ipv4.tcp_rmem="1024 4096 4096"
sysctl -w net.ipv4.tcp_wmem="1024 4096 4096"
sysctl -w net.core.somaxconn="10000"

sysctl -w net.core.netdev_max_backlog=2000
sysctl -w net.ipv4.tcp_max_syn_backlog=2048
ifconfig lo txqueuelen 5000

sysctl -w net.core.rmem_default=
