#!/bin/sh
echo "sleep 20 seconds"
sleep 20
echo "starting haproxy"
cd /var/lib/haproxy || exit
haproxy -f /usr/local/etc/haproxy/haproxy.cfg

