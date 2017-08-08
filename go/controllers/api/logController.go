package api

import (
	"github.com/zaaksam/dproxy/go/services"
)

// LogController 日志功能控制器
type LogController struct {
	BaseController
}

// List 获取白名单列表
func (c *LogController) List() {
	pageIndex, _ := c.GetInt("pageIndex")
	pageSize, _ := c.GetInt("pageSize")
	typ := c.GetString("typ")
	content := c.GetString("content")

	list, err := services.Log.Find(pageIndex, pageSize, typ, content)
	if err != nil {
		c.SetError(err)
		return
	}

	c.SetData("list", list)
}
