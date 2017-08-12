package config

import (
	"flag"
	"time"
)

// AppConf 应用全局参数
var AppConf appConf

func init() {
	flag.BoolVar(&AppConf.Debug, "debug", false, "调试模式，默认：false")
	flag.BoolVar(&AppConf.AutoOpen, "ao", true, "自动打开浏览器，默认：true")
	flag.BoolVar(&AppConf.AutoStart, "as", false, "自动启动端口映射，默认：false")
	flag.StringVar(&AppConf.IP, "ip", "", "监听的IP地址，默认：127.0.0.1")
	flag.IntVar(&AppConf.Port, "port", 0, "服务端口，默认：8080")
	flag.BoolVar(&AppConf.UI, "ui", true, "是否开启WebUI管理服务，默认：true")
	flag.StringVar(&AppConf.PrefixPath, "prefix", "", "WebUI的路径前缀，默认为空")
	flag.StringVar(&AppConf.Token, "token", "", "API授权令牌，为空时不校验，默认为空")
	flag.Parse()

	// authURL?redirect_uri=http://127.0.0.1:8080/web/auth&code=abc
	// redirect时，须在uri上附上 code、userID、userName 参数
	// flag.StringVar(&authURL, "authURL", "", "授权URL地址，不为空时开启鉴权，跳转往该地址获取授权")

	if AppConf.IP == "" {
		AppConf.IP = "127.0.0.1"
	}

	if AppConf.Port <= 0 {
		AppConf.Port = 8080
	}

	AppConf.Name = "dproxy"
	AppConf.Version = "0.2.0"
	AppConf.Started = time.Now().Unix()
}

type appConf struct {
	Name       string
	Version    string
	Debug      bool
	AutoOpen   bool
	AutoStart  bool
	UI         bool
	IP         string
	Port       int
	PrefixPath string
	Token      string
	Started    int64
}
