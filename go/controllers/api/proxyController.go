package api

import (
	"strconv"

	"github.com/zaaksam/dproxy/go/services"
)

// ProxyController 代理请求功能控制器
type ProxyController struct {
	BaseController
}

// Start 启动代理请求
func (c *ProxyController) Start() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		c.SetError("id参数错误")
		return
	}

	err = services.Proxy.Start(id)
	if err != nil {
		c.SetError(err)
	}
}

// Stop 停止代理请求
func (c *ProxyController) Stop() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		c.SetError("id参数错误")
		return
	}

	services.Proxy.Stop(id)
}
