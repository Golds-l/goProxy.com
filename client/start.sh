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
nohup /path/to/goproxy/client -cS x.x.x.x -cSP 2001 -rH x.x.x.x -rHP xx -lP xx >> $logpath &
echo $! > $pidpath
sleep 0.1
echo "client start..."