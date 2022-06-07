# 内网穿透
> 通过公网的服务器转发流量实现内网穿透，需要有一台具有公网IP的服务器。
> 支持多台内网机器同时穿透
#### 可执行文件仅linux平台，目前仅测试过穿透SSH服务!!!
## 参数说明
### 公网服务端

-rP 通信端口，与内网客户端通信，需与内网客户端中的cSP保持一致
### 内网客户端

-cS 公网服务端的IP(ipv4)

-cSP 公网服务端的通信端口，与公网服务端通信，与公网服务端中的rP保持一致

-rH 需穿透的内网机器的ip

-rHP 需穿透的内网机器的服务的端口

-lP 对应的公网端口

## 启停
可通过脚本及配置文件配合systemd使用,需要修改配置文件中脚本路径和脚本中的执行文件路径。
也可单独使用脚本启停服务，同样需修改脚本中可执行文件的路径和相关参数。
服务端默认使用2001端口与内网机器通信。
### 使用systemd管理服务
需先修改脚本中可执行文件的路径及参数和配置文件中脚本的路径！
#### 内网机器
加载配置文件:
```shell
sudo ln -S goproxy/client/goproxy.service /usr/lib/systemd/system/goproxy.service
sudo systemctl daemon-reload
```
管理服务:
```shell
sudo systemctl start goproxy.service // 启动
sudo systemctl stop goproxy.service // 停止
sudo systemctl restart goproxy.service // 重启
sudo systemctl status goproxy.service // 查看服务状态
sudo systemctl enable goproxy.service // 开机启动
sudo systemctl disable goproxy.service // 取消开机启动
```
#### 公网服务器
加载配置文件:
```shell
sudo ln -S goproxy/client/goproxy.service /usr/lib/systemd/system/goproxy.service
sudo systemctl daemon-reload
```
启停服务:
```shell
sudo systemctl start goproxy.service // 启动
sudo systemctl stop goproxy.service // 停止
sudo systemctl restart goproxy.service // 重启
sudo systemctl status goproxy.service // 查看服务状态
sudo systemctl enable goproxy.service // 开机启动
sudo systemctl disable goproxy.service // 取消开机启动
```
### 使用脚本启停服务
需先修改脚本中可执行文件的路径及设置各项参数！
#### 内网机器
```shell
sudo goproxy/client/start.sh
sudo goproxy/client/stop.sh
```
#### 公网服务器
```shell
sudo goproxy/server/proxystart.sh
sudo goproxy/server/proxystop.sh
```
