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
nohup /home/golds/goproxy/client -cS 110.40.167.60 -cSP 2001 -rH 127.0.0.1 -rHP 22 >> $logpath &
echo $! > $pidpath
sleep 0.1
echo "client start..."