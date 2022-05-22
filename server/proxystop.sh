#!/bin/bash
pid=`cat /run/goproxy.pid`
logpath="/var/log/proxy.log"
kill -9 $pid
echo "stop: $(date "+%Y-%m-%d %H:%M:%S")" >> $logpath