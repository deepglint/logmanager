package controller

import (
	//"encoding/json"
	"fmt"
	"github.com/deepglint/glog"
	"github.com/deepglint/logmanager/models"
	"github.com/deepglint/muses/util/ripple"
	//"io/ioutil"
	//"os"
	//"path/filepath"
	//"mime/multipart"
	//"time"
)

type LogController struct {
	name               string
	logForwarderConfig models.Config
}

func (this *LogController) String() string {
	return fmt.Sprintf("LogControlloer: %s", this.name)
}

func NewLogController(name string, c models.Config) (*LogController, error) {
	//fmt.Println("New LogController")
	ctl := &LogController{
		name: name,
		logForwarderConfig: models.Config{
			URL:             c.URL,
			Username:        c.Username,
			Password:        c.Password,
			UserAgent:       c.UserAgent,
			Timeout:         c.Timeout,
			RetentionPolicy: c.RetentionPolicy,
			Precision:       c.Precision,
			Consistency:     c.Consistency,
		},
	}
	ctl.name = name
	return ctl, nil
}

func (this *LogController) PostUpload(ctx *ripple.Context) {
	ctx.Request.ParseMultipartForm(32 << 20)
	file, handler, err := ctx.Request.FormFile("uploadlogfile")
	if err != nil {
		glog.Errorf("Log file upload failed: %v", err)
		ctx.Response.Status = 400
		ctx.Response.Body = err
		return
	}
	defer file.Close()

	logforwarder, err := models.NewLogForwarder(this.logForwarderConfig)
	status_code, err := logforwarder.SaveFile(file, handler) // decode file and forward to database
	if err != nil {
		ctx.Response.Status = status_code
		ctx.Response.Body = err
		return
	} else {
		glog.Infof("Save log file to database succeed")
		ctx.Response.Status = status_code
		ctx.Response.Body = "Save log file to database succeed"
		return
	}
}
