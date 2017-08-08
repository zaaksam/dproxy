# dproxy

简易的网络代理工具，带有IP白名单限制管，带有简洁的UI管理界面，提供丰富的API接口，可方便的与各个系统集成，可编译为单文件运行

网络请求代理也有不少应用了，重复做一个的原因有很多，一方面是学习，一方面是技术栈的偏好，一方面想更好的跟其他系统集成，避免太独立，总之理由很多，大家将就吧。

`注意：目前仅实现部份场景下的应用，功能并未完全实现，开源出来交流学习`

# 已编译平台

[dproxy.v0.1.1-win64.zip](https://github.com/Zaaksam/dproxy/releases/download/v0.1.1/dproxy.v0.1.1-win64.zip)

[dproxy.v0.1.1-darwin64.zip](https://github.com/Zaaksam/dproxy/releases/download/v0.1.1/dproxy.v0.1.1-darwin64.zip)


# 运行参数

```
-debug 打开调试模式，beego设置为DEV模式，静态资源调用statik，关闭时，beego设置为PROD模式，静态资源使用web/static路径资源。默认：false

-autoopen 启动应用时自动打开浏览器进行访问。默认：true

-ip 应用监听IP地址。默认：127.0.0.1

-port 应用监听端口。默认：8080

-ui 是否打开WebUI管理界面服务。默认：true
```

# 界面预览

![](https://github.com/Zaaksam/dproxy/releases/download/v0.1.1/portmap.png)

![](https://github.com/Zaaksam/dproxy/releases/download/v0.1.1/whitelist.png)

![](https://github.com/Zaaksam/dproxy/releases/download/v0.1.1/log.png)

![](https://github.com/Zaaksam/dproxy/releases/download/v0.1.1/doc.png)

# 技术资源

## Backend

语言 Go 1.8.3

Web/API开发框架 Beego 1.9

ORM框架 Xorm 0.6.2

数据库 sqlite3

静态资源打包工具 statik

依赖管理工具 govendor

## Frontend

语言 Typescript 2.4.2

JS框架 Vue 2.4.2

路由 Vue-Router 2.7.0

UI框架 iView 2.0.0

网络请求 axios 0.16.2

工具库 lodash 4.17.4

日期时间库 moment 2.18.1

打包工具 webpack 3.4.1

依赖管理工具 yarn

# 注意

sqlite3使用了CGO，在不同平台编译时，请先确保执行了以下命令：

```go
go get -u github.com/mattn/go-sqlite3
```

# 更新日志

2017-08-08 v0.1.0

* 初始化项目

---

2017-08-08 v0.1.1

* 修正启动参数不起作用的bug
* 修正webui日志管理界面错误显示的按钮
* 优化web页面静态文件缓存策略
* 增加ui启动参数来决定是否开启WebUI管理服务