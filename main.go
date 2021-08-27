package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/mattn/go-sqlite3"

	"github.com/webview/webview"
	"github.com/zaaksam/dproxy/go/config"
	"github.com/zaaksam/dproxy/go/constant"
	_ "github.com/zaaksam/dproxy/go/db"
	_ "github.com/zaaksam/dproxy/go/routers"
	"github.com/zaaksam/dproxy/go/services"
)

func main() {
	config.AppConf.Version = VERSION

	beego.BConfig.AppName = config.AppConf.Name
	beego.BConfig.ServerName = config.AppConf.Name
	beego.BConfig.WebConfig.AutoRender = false
	beego.BConfig.WebConfig.DirectoryIndex = false
	beego.BConfig.CopyRequestBody = true

	beego.BConfig.Listen.HTTPAddr = config.AppConf.IP
	beego.BConfig.Listen.HTTPPort = config.AppConf.Port

	if config.AppConf.Debug {
		beego.BConfig.RunMode = beego.DEV

		// db.Engine.ShowSQL(true)
	} else {
		beego.BConfig.RunMode = beego.PROD
	}

	regionTip := ""
	if config.AppConf.Region != "" {
		regionTip = "，" + config.AppConf.Region + " 区域"
	}
	logs.Info("====== 欢迎使用 " + config.AppConf.Name + " " + config.AppConf.Version + " (" + config.AppConf.Mode.ToString() + "模式" + regionTip + ") ，关闭此窗口即可退出程序 ======")

	if config.AppConf.Mode == constant.MODE_APIUI || config.AppConf.Mode == constant.MODE_API {
		err := services.Proxy.StartAll()
		if err != nil {
			logs.Error("端口映射任务启动失败：", err)
		}

		beego.Run()

		services.Proxy.StopAll()
	} else if config.AppConf.Mode == constant.MODE_WEB {
		go openBrowser()

		beego.Run()
	} else if config.AppConf.Mode == constant.MODE_APP {
		go beego.Run()

		w := webview.New(config.AppConf.Debug)
		w.SetTitle(config.AppConf.Name)
		// Resizable: false
		w.SetSize(960, 600, webview.HintFixed)
		w.Navigate(fmt.Sprintf("http://%s:%d/web", config.AppConf.IP, config.AppConf.Port))

		w.Run()

		// err := webview.Open("DProxy", fmt.Sprintf("http://%s:%d/web", config.AppConf.IP, config.AppConf.Port), 960, 600, false)
		// if err != nil {
		// 	logs.Critical("webview启动失败：", err)
		// }
	}
}

func openBrowser() {
	time.Sleep(2 * time.Second)

	cmds := map[string][]string{
		"windows": []string{"cmd", "/c", "start", ""},
		"darwin":  []string{"open", ""},
		"linux":   []string{"xdg-open", ""},
	}

	osStr := runtime.GOOS

	args, ok := cmds[osStr]
	if !ok {
		logs.Error("非预设系统，无法打开浏览器：", osStr)
		return
	}

	args[len(args)-1] = fmt.Sprintf("http://%s:%d/web", config.AppConf.IP, config.AppConf.Port)

	c := exec.Command(args[0], args[1:]...)
	err := c.Start()
	if err != nil {
		logs.Critical("打开浏览器失败：", err)
	}
}
