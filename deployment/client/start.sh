#!/bin/sh
echo "Client"
echo "sleep 30 seconds"
sleep 30
echo "starting client"
./bin/client -URL http://lab10_nginx01:8080/getname

