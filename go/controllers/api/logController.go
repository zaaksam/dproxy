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
	created, _ := c.GetInt64("created")
	typ := c.GetString("type")
	content := c.GetString("content")

	list, err := services.Log.Find(pageIndex, pageSize, created, typ, content)
	if err != nil {
		c.SetError(err)
		return
	}

	c.SetData("list", list)
}

// Delete 删除白名单记录
func (c *LogController) Delete() {
	created, _ := c.GetInt64("created")
	typ := c.GetString("type")
	content := c.GetString("content")

	err := services.Log.Delete(created, typ, content)
	if err != nil {
		c.SetError(err)
	}
}
