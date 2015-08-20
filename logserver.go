package main

import (
	"flag"
	"github.com/deepglint/glog"
	"github.com/deepglint/logmanager/controller"
	"github.com/deepglint/logmanager/models"
	"github.com/deepglint/muses/util/ripple"
	"net/http"
	"time"
)

type LogConfig struct {
	Name               string
	ListenPort         string
	logForwarderConfig models.Config
}

func main() {
	var config LogConfig
	flag.StringVar(&config.Name, "name", "log_server", "Server name")
	flag.StringVar(&config.ListenPort, "port", ":1734", "Listening port")
	flag.StringVar(&config.logForwarderConfig.URL.Scheme, "scheme", "http", "Set url scheme to http")
	flag.StringVar(&config.logForwarderConfig.URL.Host, "influxdb_url", "localhost:8086", "Influxdb host and port")
	flag.StringVar(&config.logForwarderConfig.Username, "influxdb_username", "", "Influxdb basic auth username")
	flag.StringVar(&config.logForwarderConfig.Password, "influxdb_password", "", "Influxdb basic auth password")
	flag.StringVar(&config.logForwarderConfig.UserAgent, "influxdb_user_agent", "", "User agent")
	flag.DurationVar(&config.logForwarderConfig.Timeout, "influxdb_timeout", time.Duration(10)*time.Second, "Influxdb request time out")
	flag.StringVar(&config.logForwarderConfig.RetentionPolicy, "influxdb_retention_policy", "3d", "Time to keep old data before clean it")
	flag.StringVar(&config.logForwarderConfig.Precision, "influxdb_timestamp_precision", "n", "Timestamp precision")
	flag.StringVar(&config.logForwarderConfig.Consistency, "influxdb_node_consistency", "one", "Influxdb nodes write consistency")
	flag.Parse()
	glog.Infof("Config: %v", config)

	app := ripple.NewApplication()
	logcontroller, err := controller.NewLogController(config.Name, config.logForwarderConfig)
	if err != nil {
		glog.Errorf("Log controller init failed: %v", err)
	}
	app.RegisterController("log", logcontroller)
	//app.AddRoute(ripple.Route{Pattern: "_controller/_action"})
	app.AddRoute(ripple.Route{Pattern: "upload", Controller: "log", Action: "Upload"})
	app.SetBaseUrl("/")
	http.HandleFunc("/", app.ServeHTTP)

	glog.Infof("Starting centralized log server at ", config.ListenPort)
	http.ListenAndServe(config.ListenPort, nil)
	defer glog.Flush()
}
