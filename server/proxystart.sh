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
echo "start: $(date "+%Y-%m-%d %H:%M:%S")" >> $logpath
nohup /path/to/goproxy/server -lP 2000 -rP 2001 >> $logpath &
echo $! > $pidpath
sleep 0.1
echo "start.."