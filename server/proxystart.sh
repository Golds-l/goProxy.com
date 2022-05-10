#!/bin/bash
if [ -e "./proxy.log" ]; then
    echo 
else
    touch proxy.log
fi
echo "start: $(date "+%Y-%m-%d %H:%M:%S")" >> proxy.log
nohup ../bin/server -lP 2000 -rP 2001 >> proxy.log &
sleep 0.1
echo "start.."