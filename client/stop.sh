#!/bin/bash
pid=`cat /run/goproxy.pid`
logpath="/home/golds/goproxy/proxy.log"
kill -9 $pid
echo "stop: $(date "+%Y-%m-%d %H:%M:%S")" >> proxy.log