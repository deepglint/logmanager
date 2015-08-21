package controller

import (
	"fmt"
	//"github.com/deepglint/glog"
	"github.com/deepglint/logmanager/client/client"
	"github.com/deepglint/muses/util/ripple"
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
	client.Flush(this.url, this.dir)
}
