if [ -e "./proxy.log" ]; then
    echo 
else
    touch proxy.log
fi
echo "start: $(date "+%Y-%m-%d %H:%M:%S")" >> proxy.log
nohup ./client -cS 110.40.167.60 -cSP 2001 -rH 127.0.0.1 -rHP 22 >> proxy.log &
sleep 0.1
echo "client start..."