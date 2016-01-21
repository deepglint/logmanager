# logmanager客户端
```
服务名称：logmanager/logclient
服务版本号：0.1.3
提交者：kaixian hu
日期：2016-01-21
```

##logclient
	upload log file to server

###改动日志
改动 | 单元测试
:----------:|:----------:
ADD: log upload load control |Y
ADD: log quatity control |Y
ADD: new database table structure |Y
CHANGE: log tailing method |Y
CHANGE: log cleaning method |Y

###参数
```
 Usage of ./logclient:
  -alsologtostderr_deepglint=false: log to standard error as well as files
  -client_listen_port=":1735": Log client server listening port
  -debug=false: set debug to true to output local debug file
  -dir="/tmp/": Upload Directory
  -keep_interval=30m0s: Log file kept time (better be bigger than sleep_interval and upload_interval)
  -line_num_limit=5000: Line number limit for every log file
  -log_backtrace_at_deepglint=:0: when logging hits line file:N, emit a stack trace
  -log_dir_deepglint="/tmp/": If non-empty, write log files in this directory
  -log_file_name_interval=5m0s: Log file name interval, create a new log file every interval
  -log_level=1: 0 for all log, 1 for log above warning, 2 for log above error, 3 for fatal only, 4 for no log, default set to 1
  -logtostderr_deepglint=false: log to standard error instead of files
  -method="/upload": Log client method
  -name="log_client": Log cilent name
  -server_host="http://localhost": Log server host
  -server_port=":1734": Log server listen port
  -sleep_interval=10m0s: Sleep time interval between every upload action (better smaller than keep_interval)
  -stderr=0: logs at or above this threshold go to stderr
  -table_name="libra": table name for database
  -tailed_file_interval=24h0m0s: Tailed file clean interval
  -upload_interval=15m0s: Upload file created before upload interval (better be smaller than keep_interval, has to be bigger than log file create interval)
  -v_deepglint=0: log level for V logs
  -vmodule_deepglint=: comma-separated list of pattern=N settings for file-filtered logging
```
###api
	同步客户端最新的日志到网管服务器
	URL: http://clienthost:1735/sync
	method: HTTP GET
	
	从web端查看本地日志
	URL: http://clienthost:1735/locallog
  	method: HTTP GET
  	params: interval 查看interval到现在的本地日志 
  	示例: http://clienthost:1735/locallog?interval=5m
	
	
###arm正式启动示例
  sudo docker run -i -t -d -p 1735:1735 --restart=always --net=host --name=logclient -v=/etc/localtime:/etc/localtime:ro -v=/tmp:/tmp 192.168.5.46:5000/armhf-logclient:0.1.3 ./logclient.arm -keep_interval=30m -upload_interval=5m -sleep_interval=5m -server_host http://192.168.5.46 -dir /tmp/ -log_dir_deepglint /tmp/

###arm测试启动示例
  sudo docker run -i -t -d -p 1735:1735 --restart=always --net=host --name=logclient -v=/etc/localtime:/etc/localtime:ro -v=/tmp:/tmp 192.168.5.46:5000/armhf-logclient:0.1.3 ./logclient.arm -keep_interval=30m -upload_interval=5m -sleep_interval=5m -server_host http://192.168.5.46 -dir /tmp/ -log_dir_deepglint /tmp/ -log_level=0
	
###amd64正式启动示例
  sudo docker run -i -t -d -p 1735:1735 --restart=always --net=host --name=logclient -v=/etc/localtime:/etc/localtime:ro -v=/tmp:/tmp 192.168.5.46:5000/logclient:0.1.3 ./logclient.linux -keep_interval=30m -upload_interval=5m -sleep_interval=5m -server_host http://192.168.5.46 -dir /tmp/ -log_dir_deepglint /tmp/

###amd64测试启动示例
  sudo docker run -i -t -d -p 1735:1735 --restart=always --net=host --name=logclient -v=/etc/localtime:/etc/localtime:ro -v=/tmp:/tmp 192.168.5.46:5000/logclient:0.1.3 ./logclient.linux -keep_interval=30m -upload_interval=5m -sleep_interval=5m -server_host http://192.168.5.46 -dir /tmp/ -log_dir_deepglint /tmp/ -log_level=0

##Demo地址

##查询已上传日志方法
```
查询已有数据库		http://192.168.5.46:8088/query?q=show databases

查询数据库下表名	http://192.168.5.46:8088/query?db=xxx&q=show measurements

查询详细		http://192.168.5.46:8088/query?db=xxx&q=select * from xxxx
```
