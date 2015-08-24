# logmanager

##dependencies
	github.com/deepglint/glog
	github.com/deepglint/muses/util/ripple
##server
receive upload log file from client, and forward it to influxdb

###build
	go build logserver.go
###parameters
```
Usage of ./logserver:
  -alsologtostderr_deepglint=false: log to standard error as well as files
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

##client
upload log file to server

###build 
	go build logclient.go
###parameters
 ```
 Usage of ./logclient:
  -alsologtostderr_deepglint=false: log to standard error as well as files
  -client_listen_port=":1735": Log client server listening port
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