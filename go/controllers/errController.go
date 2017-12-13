package controllers

import (
	"net/http"

	"github.com/astaxie/beego"
)

// ErrController 错误处理控制器
type ErrController struct {
	beego.Controller
}

// Error404 404页面
func (c *ErrController) Error404() {
	// c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	// c.Ctx.Output.SetStatus(404)
	// c.Ctx.Output.Body([]byte("404"))
	http.Error(c.Ctx.ResponseWriter, "", 404)
}

// Error500 500页面
func (c *ErrController) Error500() {
	// c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	// c.Ctx.Output.SetStatus(500)
	// c.Ctx.Output.Body([]byte("500"))
	http.Error(c.Ctx.ResponseWriter, "", 500)
}
