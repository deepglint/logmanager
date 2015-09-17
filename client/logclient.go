package main

import (
	"flag"
	"fmt"
	"github.com/deepglint/glog"
	"github.com/deepglint/logmanager/client/client"
	"github.com/deepglint/logmanager/client/controller"
	"github.com/deepglint/muses/util/ripple"
	"net/http"
	"time"
)

type LogClientConfig struct {
	Name           string
	Host           string
	Port           string
	Method         string
	Dir            string
	ListenPort     string
	UploadInterval time.Duration
	KeepInterval   time.Duration
	SleepInterval  time.Duration
}

func main() {
	var config LogClientConfig
	flag.StringVar(&config.Name, "name", "log_client", "Log cilent name")
	flag.StringVar(&config.Host, "server_host", "http://localhost", "Log server host")
	flag.StringVar(&config.Port, "server_port", ":1734", "Log server listen port")
	flag.StringVar(&config.Method, "method", "/upload", "Log client method")
	flag.StringVar(&config.Dir, "dir", "/tmp/", "Upload Directory")
	flag.StringVar(&config.ListenPort, "client_listen_port", ":1735", "Log client server listening port")
	flag.DurationVar(&config.UploadInterval, "upload_interval", time.Duration(15)*time.Minute, "Upload file created before upload interval (better be smaller than keep_interval, has to be bigger than log file create interval)")
	flag.DurationVar(&config.KeepInterval, "keep_interval", time.Duration(30)*time.Minute, "Log file kept time (better be bigger than sleep_interval and upload_interval)")
	flag.DurationVar(&config.SleepInterval, "sleep_interval", time.Duration(10)*time.Minute, "Sleep time interval between every upload action (better smaller than keep_interval)")
	flag.Parse()

	// if config.KeepInterval < config.UploadInterval || config.KeepInterval < config.SleepInterval {
	// 	if config.UploadInterval >= config.SleepInterval {
	// 		config.KeepInterval = config.UploadInterval * 2
	// 		// fmt.Printf("%v\n", config.KeepInterval)
	// 	} else {
	// 		config.KeepInterval = config.SleepInterval * 2
	// 		// fmt.Printf("%v\n", config.KeepInterval)
	// 	}
	// }
	if config.Dir[len(config.Dir)-1] != '/' {
		config.Dir += "/"
	}

	glog.Infof("Log file will be kept for %v, log file will be uploaded %v after created, log client will run every %v", config.KeepInterval, config.UploadInterval, config.SleepInterval)
	fmt.Printf("Log file will be kept for %v, log file will be uploaded %v after created, log client will run every %v\n", config.KeepInterval, config.UploadInterval, config.SleepInterval)

	go func() {
		for {
			client.SendLog(config.Host+config.Port+config.Method, config.Dir, config.UploadInterval, config.KeepInterval)
			time.Sleep(config.SleepInterval)
		}
	}()

	app := ripple.NewApplication()
	logclientcontroller, err := controller.NewLogClientController("logclientserver", config.Host+config.Port+config.Method, config.Dir)
	if err != nil {
		glog.Errorf("Log client controller init failed: %v", err)
	}
	app.RegisterController("logclient", logclientcontroller)
	app.AddRoute(ripple.Route{Pattern: "sync", Controller: "logclient", Action: "Sync"})
	app.AddRoute(ripple.Route{Pattern: "favicon.ico", Controller: "logclient", Action: "Favicon"})
	app.AddRoute(ripple.Route{Pattern: "locallog", Controller: "logclient", Action: "LocalLogFile"})
	app.SetBaseUrl("/")
	http.HandleFunc("/", app.ServeHTTP)
	glog.Infof("Starting log client server at ", config.ListenPort)
	http.ListenAndServe(config.ListenPort, nil)

}
