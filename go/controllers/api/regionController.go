package api

import (
	"github.com/zaaksam/dproxy/go/services"
)

// RegionController 区域控制器
type RegionController struct {
	BaseController
}

// List 获取区域列表
func (c *RegionController) List() {
	pageIndex, _ := c.GetInt("pageIndex")
	pageSize, _ := c.GetInt("pageSize")

	list, err := services.Region.Find(pageIndex, pageSize)
	if err != nil {
		c.SetError(err)
		return
	}

	c.SetData("list", list)
}
