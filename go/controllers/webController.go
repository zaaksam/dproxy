package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/Zaaksam/dproxy/go/views"
	"github.com/astaxie/beego"
)

type WebController struct {
	beego.Controller
}

func (c *WebController) Get() {
	unixStr := ""
	if beego.BConfig.RunMode == beego.DEV {
		unixStr = "?t=" + strconv.FormatInt(time.Now().Unix(), 10)
	}

	index := strings.Replace(views.Index, "{{.unix}}", unixStr, -1)
	index = strings.Replace(index, "{{.appname}}", beego.BConfig.AppName, -1)

	if c.Ctx.ResponseWriter.Header().Get("Content-Type") == "" {
		c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	}

	c.Ctx.Output.Body([]byte(index))
}
