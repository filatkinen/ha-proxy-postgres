#!/bin/sh
echo "sleep 10 seconds"
sleep 10
echo "starting haproxy"
cd /var/lib/haproxy || exit
haproxy -f /usr/local/etc/haproxy/haproxy.cfg

