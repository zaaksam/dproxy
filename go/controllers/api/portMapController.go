package api

import (
	"encoding/json"
	"strconv"

	"github.com/zaaksam/dproxy/go/model"
	"github.com/zaaksam/dproxy/go/services"
)

// PortMapController 端口映射功能控制器
type PortMapController struct {
	BaseController
}

// Get 获取端口映射详情
func (c *PortMapController) Get() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		c.SetError("id参数错误")
		return
	}

	md, err := services.PortMap.Get(id)
	if err != nil {
		c.SetError(err)
		return
	}

	c.SetData("portMap", md)
}

// List 获取端口映射列表
func (c *PortMapController) List() {
	pageIndex, _ := c.GetInt("pageIndex")
	pageSize, _ := c.GetInt("pageSize")
	region := c.GetString("region")
	targetIP := c.GetString("targetIP")
	targetPort := c.GetString("targetPort")
	sourcePort := c.GetString("sourcePort")
	sortField := c.GetString("sortField")
	sortDesc := c.GetString("sortDesc")

	list, err := services.PortMap.Find(pageIndex, pageSize, region, targetIP, targetPort, sourcePort, sortField, sortDesc, nil)
	if err != nil {
		c.SetError(err)
		return
	}

	c.SetData("list", list)
}

// Put 创建或修改 端口映射
func (c *PortMapController) Put() {
	md := &model.PortMapModel{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, md)
	if err != nil {
		c.SetError("参数错误，请检查后重试：" + err.Error())
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	if idStr == "" {
		md, err = services.PortMap.Add(md)
	} else {
		md.ID, err = strconv.ParseInt(idStr, 10, 0)
		if err != nil {
			c.SetError("id参数错误")
			return
		}

		err = services.PortMap.Update(md)
	}
	if err != nil {
		c.SetError(err)
		return
	}

	c.SetData("portMap", md)
}

// Delete 删除端口映射记录
func (c *PortMapController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		c.SetError("id参数错误")
		return
	}

	err = services.PortMap.Delete(id)
	if err != nil {
		c.SetError(err)
	}
}
