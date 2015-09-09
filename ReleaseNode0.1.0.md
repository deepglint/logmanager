# logmanager
```
服务名称：logmanager
服务版本号：0.1.0
提交者：kaixian hu
日期：2015-08-27
```

##依赖包
	github.com/deepglint/muses/util/ripple
##服务器端
	receive upload log file from client, and forward it to influxdb

###编译命令
	go build logserver.go
###参数
```
Usage of ./logserver:
  -alsologtostderr_deepglint=false: log to standard error as well as files
  -debug=false: set debug to true to output local debug file
  -influxdb_node_consistency="one": Influxdb nodes write consistency
  -influxdb_password="": Influxdb basic auth password
  -influxdb_retention_policy="3d": Time to keep old data before clean it
  -influxdb_timeout=10s: Influxdb request time out 
  -influxdb_timestamp_precision="n": Timestamp precision
  -influxdb_url="localhost:8086": Influxdb host and port
  -influxdb_user_agent="": User agent
  -influxdb_username="": Influxdb basic auth username
  -log_backtrace_at_deepglint=:0: when logging hits line file:N, emit a stack trace
  -log_dir_deepglint="": If non-empty, write log files in this directory
  -log_file_name_interval=5m0s: log file name interval, create a new log file every interval
  -logtostderr_deepglint=false: log to standard error instead of files
  -name="log_server": Server name
  -port=":1734": Listening port
  -scheme="http": Set url scheme to http
  -stderr=0: logs at or above this threshold go to stderr
  -v_deepglint=0: log level for V logs
  -vmodule_deepglint=: comma-separated list of pattern=N settings for file-filtered logging
```
###启动示例
	sudo ./logserver -log_dir_deepglint /tmp/ -influxdb_url 192.168.5.46:8088
##客户端
	upload log file to server

###编译命令
	go build logclient.go
###参数
```
 Usage of ./logclient:
  -alsologtostderr_deepglint=false: log to standard error as well as files
  -client_listen_port=":1735": Log client server listening port
  -debug=false: set debug to true to output local debug file
  -dir="./": Upload Directory
  -keep_interval=10m0s: Log file kept time (better be bigger than sleep_interval and upload_interval)
  -log_backtrace_at_deepglint=:0: when logging hits line file:N, emit a stack trace
  -log_dir_deepglint="": If non-empty, write log files in this directory
  -log_file_name_interval=5m0s: log file name interval, create a new log file every interval
  -logtostderr_deepglint=false: log to standard error instead of files
  -method="/upload": Log client method
    -name="log_client": Log cilent name
  -server_host="http://localhost": Log server host
  -server_port=":1734": Log server listen port
  -sleep_interval=3m0s: Sleep time interval between every upload action (better smaller than keep_interval)
  -stderr=0: logs at or above this threshold go to stderr
  -upload_interval=5m0s: Upload file created before upload interval (better be smaller than keep_interval)
  -v_deepglint=0: log level for V logs
  -vmodule_deepglint=: comma-separated list of pattern=N settings for file-filtered logging
```
###api
	同步客户端最新的日志到网管服务器
	URL: http://clienthost:1735/sync
	method: HTTP GET

  URL: http://clienthost:1735/locallog
  method: HTTP GET
  params: interval 
	
	
###启动示例
	sudo ./logclient -dir /tmp/ -keep_interval 1h -upload_interval 10m -sleep_interval 10m -server_host http://192.168.5.46
	
##测试方法
	用logmanager目录下的glog替换程序原来使用的glog，重新编译后启动logclient并设置参数为demo地址即可

##Demo地址
	server: 192.168.5.46:8088
##查询已上传日志方法
```
程序生成日志文件名称: 	LOG.XXX.MST2006-08-27-12:00:00 XXX即为database名称
查询已有数据库		http://192.168.5.46:8088/query?q=show databases

程序名称即为表名
查询数据库下表名	http://192.168.5.46:8088/query?db=xxx&q=show measurements

查询详细		http://192.168.5.46:8088/query?db=xxx&q=select * from xxxx
```
