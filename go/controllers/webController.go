package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/zaaksam/dproxy/go/config"
	"github.com/zaaksam/dproxy/go/views"
)

type WebController struct {
	beego.Controller
}

func (c *WebController) Get() {
	var unix int64
	if config.AppConf.Debug {
		unix = time.Now().Unix()
	} else {
		unix = config.AppConf.Started
	}
	unixStr := "?t=" + strconv.FormatInt(unix, 10)

	index := strings.Replace(views.Index, "{{.unix}}", unixStr, -1)
	index = strings.Replace(index, "{{.appname}}", beego.BConfig.AppName, -1)

	if c.Ctx.ResponseWriter.Header().Get("Content-Type") == "" {
		c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	}

	c.Ctx.Output.Body([]byte(index))
}
