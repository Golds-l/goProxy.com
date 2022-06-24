#!/bin/bash
logpath="/var/log/proxy.log"
pidpath="/run/goproxy.pid"
if [ -e $logpath ]; then
    echo 
else
    touch $logpath
fi
if [ -e $pidpath ]; then
    echo 
else
    touch $pidpath
fi
nohup /path/to/goproxy/server -rP 2001 >> $logpath &
echo $! > $pidpath
sleep 0.1
echo "start.."