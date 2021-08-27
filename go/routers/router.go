package routers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/rakyll/statik/fs"
	"github.com/zaaksam/dproxy/go/config"
	"github.com/zaaksam/dproxy/go/constant"
	"github.com/zaaksam/dproxy/go/controllers"
	"github.com/zaaksam/dproxy/go/controllers/api"

	// 加载静态资源
	_ "github.com/zaaksam/dproxy/go/statik"
)

func init() {
	beego.ErrorController(&controllers.ErrController{})

	api := beego.NewNamespace("/api",
		beego.NSRouter("/whitelist/clear", &api.WhiteListController{}, "delete:Clear"),
		beego.NSRouter("/whitelist/list", &api.WhiteListController{}, "get:List"),
		beego.NSRouter("/whitelist/getip", &api.WhiteListController{}, "get:GetIP"),
		beego.NSRouter("/whitelist/?:id:int", &api.WhiteListController{}),
		beego.NSRouter("/portmap/list", &api.PortMapController{}, "get:List"),
		beego.NSRouter("/portmap/?:id:int", &api.PortMapController{}),
		beego.NSRouter("/proxy/start/:id:int", &api.ProxyController{}, "get:Start"),
		beego.NSRouter("/proxy/stop/:id:int", &api.ProxyController{}, "get:Stop"),
		beego.NSRouter("/log/list", &api.LogController{}, "get:List"),
		beego.NSRouter("/log", &api.LogController{}),
		beego.NSRouter("/region/list", &api.RegionController{}, "get:List"),
	)

	if config.AppConf.Mode == constant.MODE_API {
		//只开启API服务
		beego.AddNamespace(api)
		return
	}

	web := beego.NewNamespace("/web", beego.NSRouter("/*", &controllers.WebController{}))
	var staticHandler http.Handler
	if config.AppConf.Debug {
		// dir := path.Dir(os.Args[0])
		dir, _ := os.Getwd()
		dir = filepath.ToSlash(dir) // .../dproxy/go
		// dirs := strings.Split(dir, "/")
		// dir = strings.Join(dirs[0:len(dirs)-1], "/") // .../dproxy
		staticHandler = http.StripPrefix("/static/", http.FileServer(http.Dir(dir+"/web/static")))
	} else {
		statikFS, err := fs.New()
		if err != nil {
			logs.Error("statik获取失败：", err)
		}
		staticHandler = http.StripPrefix("/static/", http.FileServer(statikFS))
	}

	static := beego.NewNamespace("/static", beego.NSGet("/*", func(ctx *context.Context) {
		staticHandler.ServeHTTP(ctx.ResponseWriter, ctx.Request)
	}))

	beego.AddNamespace(static, api, web)
}
