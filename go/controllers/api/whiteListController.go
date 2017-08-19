package api

import (
	"encoding/json"
	"strconv"

	"github.com/zaaksam/dproxy/go/model"
	"github.com/zaaksam/dproxy/go/services"
)

// WhiteListController 白名单功能控制器
type WhiteListController struct {
	BaseController
}

// Get 获取白名单详情
func (c *WhiteListController) Get() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		c.SetError("id参数错误")
		return
	}

	md, err := services.WhiteList.Get(id)
	if err != nil {
		c.SetError(err)
		return
	}

	c.SetData("whiteList", md)
}

// List 获取白名单列表
func (c *WhiteListController) List() {
	pageIndex, _ := c.GetInt("pageIndex")
	pageSize, _ := c.GetInt("pageSize")
	ip := c.GetString("ip")
	userName := c.GetString("userName")
	isExpired, _ := c.GetBool("isExpired")

	list, err := services.WhiteList.Find(pageIndex, pageSize, ip, userName, isExpired)
	if err != nil {
		c.SetError(err)
		return
	}

	c.SetData("list", list)
}

// Put 创建或修改 白名单
func (c *WhiteListController) Put() {
	md := &model.WhiteListModel{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, md)
	if err != nil {
		c.SetError("参数错误，请检查后重试：" + err.Error())
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	if idStr == "" {
		// md.IP = c.Ctx.Input.IP()
		md, err = services.WhiteList.Add(md)
	} else {
		md.ID, err = strconv.ParseInt(idStr, 10, 0)
		if err != nil {
			c.SetError("id参数错误")
			return
		}

		err = services.WhiteList.Update(md)
	}
	if err != nil {
		c.SetError(err)
		return
	}

	c.SetData("whiteList", md)
}

// Delete 删除白名单记录
func (c *WhiteListController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		c.SetError("id参数错误")
		return
	}

	err = services.WhiteList.Delete(id)
	if err != nil {
		c.SetError(err)
	}
}

// GetIP 获取请求者IP
func (c *WhiteListController) GetIP() {
	c.SetData("ip", c.Ctx.Input.IP())
}
