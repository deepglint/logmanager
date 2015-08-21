package main

import (
	"flag"
	"github.com/deepglint/glog"
	"github.com/deepglint/logmanager/client/client"
	"github.com/deepglint/logmanager/client/controller"
	"github.com/deepglint/muses/util/ripple"
	"net/http"
	"time"
)

type LogClientConfig struct {
	Name       string
	Host       string
	Port       string
	Method     string
	Dir        string
	ListenPort string
	SleepTime  time.Duration
}

func main() {
	var config LogClientConfig
	flag.StringVar(&config.Name, "name", "log_client", "Log cilent name")
	flag.StringVar(&config.Host, "server_host", "http://localhost", "Log server host")
	flag.StringVar(&config.Port, "server_port", ":1734", "Log server listen port")
	flag.StringVar(&config.Method, "method", "/upload", "Log client method")
	flag.StringVar(&config.Dir, "dir", "./", "Upload Directory")
	flag.StringVar(&config.ListenPort, "client_listen_port", ":1735", "Log client server listening port")
	flag.DurationVar(&config.SleepTime, "sleep_duration", time.Duration(3)*time.Minute, "Sleep time between every upload action")
	flag.Parse()

	go func() {
		for {
			client.SendLog(config.Host+config.Port+config.Method, config.Dir)
			time.Sleep(config.SleepTime)
		}
	}()

	app := ripple.NewApplication()
	logclientcontroller, err := controller.NewLogClientController("logclientserver", config.Host+config.Port+config.Method, config.Dir)
	if err != nil {
		glog.Errorf("Log client controller init failed: %v", err)
	}
	app.RegisterController("logclient", logclientcontroller)
	app.AddRoute(ripple.Route{Pattern: "sync", Controller: "logclient", Action: "Sync"})
	app.SetBaseUrl("/")
	http.HandleFunc("/", app.ServeHTTP)
	glog.Infof("Starting log client server at ", config.ListenPort)
	http.ListenAndServe(config.ListenPort, nil)

}
