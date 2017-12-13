# dproxy

[![Go Report Card](https://goreportcard.com/badge/github.com/zaaksam/dproxy)](https://goreportcard.com/report/github.com/zaaksam/dproxy)

简易的网络代理工具，带有IP白名单限制管，带有简洁的UI管理界面，提供丰富的API接口，可方便的与各个系统集成，可编译为单文件运行

`注意：目前仅实现部份场景下的应用，功能并未完全实现，开源出来交流学习`


# 运行参数

```
-debug 打开调试模式，beego设置为DEV模式，静态资源调用statik，关闭时，beego设置为PROD模式，静态资源使用web/static路径资源。默认：false

-mode 运行模式：server：API模式，不提供WebUI界面，自动运行所有端口映射任务；web：Web模式；默认：web

-ip 应用监听IP地址。默认：127.0.0.1

-port 应用监听端口。默认：8080

-prefix WebUI的路径前缀，默认空

-token API请求时的校验令牌，非空时API请求的URL须带上此参数，如：`/?token=abc`
```

# 界面预览

![](https://static.oschina.net/uploads/img/201712/13181727_jyz0.png)

![](https://static.oschina.net/uploads/img/201712/13181739_ULte.png)

![](https://static.oschina.net/uploads/img/201712/13181748_RkMe.png)

![](https://static.oschina.net/uploads/img/201712/13181825_v0Lo.png)

# 技术资源

Backend Go + beego

Frontend Typescript + Vue + iView

# 注意

sqlite3使用了CGO，在不同平台编译时，请先确保执行了以下命令：

```go
go get -u github.com/mattn/go-sqlite3
```

# 更新日志

2017-12-13 v0.4.0

* 更新go依赖项
* 更新web依赖项，优化编译配置文件
* 增加App（webview，默认）独立运行模式

2017-09-12 v0.3.3

* 重写声明文件
* 拆分webpack为dev、prod环境
* 前端重构，并使用async/await
* Log增加删除功能
* 升级相关依赖项

---

2017-08-26 v0.3.2

* 优化代理逻辑
* 采用go1.9编译

---

2017-08-23 v0.3.1

* 修改启动参数，精简为 mode 的设置形式
* 增加App模式(webview)，跨平台有兼容问题，入口暂时屏蔽
* 白名单现在可以设置和修改过期时间了
* 优化请求代理的错误处理

---

2017-08-21 v0.2.2

* 修正端口映射源资料修改无效的bug
* web类message提示持续时间改为5秒

---

2017-08-19 v0.2.1

* 白名单列表API增加 isExpired 参数

---

2017-08-12 v0.2.0

* 增加WebUI的前缀路径命令行参数：-prefix
* 增加API调用令牌校验命令行参数：-token
* 增加端口映射自动启动命令行参数：-as
* -autoopen 命令行参数简写为 -ao
* WebUI管理功能放开更多操作空间
* 整体设计倾向为API+后台服务为主，WebUI为辅

---

2017-08-08 v0.1.1

* 修正启动参数不起作用的bug
* 修正webui日志管理界面错误显示的按钮
* 优化web页面静态文件缓存策略
* 增加ui启动参数来决定是否开启WebUI管理服务

---

2017-08-08 v0.1.0

* 初始化项目