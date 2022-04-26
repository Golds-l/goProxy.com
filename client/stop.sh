string=$(ps -ef | grep ./client)
array=(`echo $string | tr ' ' ' '` )
echo ${array[1]}
kill -9 ${array[1]}
echo "stop: $(date "+%Y-%m-%d %H:%M:%S")" >> proxy.log