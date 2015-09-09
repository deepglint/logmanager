package controller

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	//"github.com/deepglint/glog"
	"github.com/deepglint/glog"
	"github.com/deepglint/logmanager/client/client"
	"github.com/deepglint/muses/util/ripple"
	// "time"
)

type LogClientController struct {
	name string
	url  string
	dir  string
}

func (this *LogClientController) String() string {
	return fmt.Sprintf("LogClientControlloer: %s", this.name)
}

func NewLogClientController(name, hosturl, logdir string) (*LogClientController, error) {
	ctl := &LogClientController{
		name: name,
		url:  hosturl,
		dir:  logdir,
	}
	return ctl, nil
}

func (this *LogClientController) GetSync(ctx *ripple.Context) {
	err := client.Flush(this.url, this.dir)
	if err != nil {
		ctx.Response.Status = 500
		ctx.Response.Body = "Latest log file uploading failed."
	} else {
		ctx.Response.Status = 200
		ctx.Response.Body = "Latest log file has been uploaded."
	}
}

func (this *LogClientController) GetFavicon(ctx *ripple.Context) {
	ctx.Response.Status = 200
	ctx.Response.Body = nil
}

func (this *LogClientController) GetLocalLogFile(ctx *ripple.Context) {
	interval_str := ctx.Params["interval"]
	// fmt.Println(interval_str)
	if interval_str == "" {
		ctx.Response.Status = 400
		ctx.Response.Body = "Unable to find parameter \"interval\"."
		return
	}
	interval, err := time.ParseDuration(interval_str)
	if err != nil {
		ctx.Response.Status = 400
		ctx.Response.Body = "Wrong time interval format( time.Duration format e.g. 10m or 1h or 5s)"
		return
	}
	logs, err := this.uploadlocal(interval)
	if err != nil {
		ctx.Response.Status = 500
		ctx.Response.Body = "No log file to upload."
		return
	}
	ctx.Response.Status = 200
	ctx.Response.Body = string(logs)
}

func (this *LogClientController) uploadlocal(interval time.Duration) ([]byte, error) {
	var logs []byte
	err := filepath.Walk(this.dir, func(filename string, f os.FileInfo, err error) error {
		if f == nil {
			glog.Errorf("No log file found")
			return nil
		}
		if f.IsDir() {
			return nil
		}
		if b, _ := path.Match("LOG.*.*????-??-??T??:??:??Z", f.Name()); b {
			fields := strings.Split(f.Name(), ".")
			t, _ := time.Parse("MST2006-01-02T15:04:05Z", fields[2])
			time_int := t.Unix()
			if time.Now().Unix()-int64(interval.Seconds()) < time_int {
				// fmt.Println(filename)
				log, err := ioutil.ReadFile(filename)
				if err != nil {
					glog.Errorf("%v", err)
					return nil
				}
				logs = append(logs, log...)
				// fmt.Println(string(logs))
				// logs = tmp
			}
		}
		return nil
	})
	return logs, err
}
