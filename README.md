# dproxy

简易的网络代理工具，带有IP白名单限制管，带有简洁的UI管理界面，提供丰富的API接口，可方便的与各个系统集成，可编译为单文件运行

`注意：目前仅实现部份场景下的应用，功能并未完全实现，开源出来交流学习`

# 已编译平台

[dproxy.v0.2.2-win64.zip](http://git.oschina.net/zaaksam/dproxy/attach_files/download?i=91560&u=http%3A%2F%2Ffiles.git.oschina.net%2Fgroup1%2FM00%2F01%2FBA%2FPaAvDFmaOt-AOq8tAHxzCWkFB80496.zip%3Ftoken%3D7b7777eaa9643222e9814158318b0fa7%26ts%3D1503280320%26attname%3Ddproxy.v0.2.2-win64.zip)

[dproxy.v0.2.2-darwin64.zip](http://git.oschina.net/zaaksam/dproxy/attach_files/download?i=91561&u=http%3A%2F%2Ffiles.git.oschina.net%2Fgroup1%2FM00%2F01%2FBA%2FPaAvDFmaPC6AKqepAIa8E6YTJdY650.zip%3Ftoken%3Dad90c6254e454abad363ab2be9bf0ea2%26ts%3D1503280320%26attname%3Ddproxy.v0.2.2-darwin64.zip)


# 运行参数

```
-debug 打开调试模式，beego设置为DEV模式，静态资源调用statik，关闭时，beego设置为PROD模式，静态资源使用web/static路径资源。默认：false

-as 启动应用时自动启动所有端口映射任务。默认：false

-ao 启动应用时自动打开浏览器进行访问。默认：true

-ip 应用监听IP地址。默认：127.0.0.1

-port 应用监听端口。默认：8080

-ui 是否打开WebUI管理界面服务。默认：true

-prefix WebUI的路径前缀，默认空

-token API请求时的校验令牌，非空时API请求的URL须带上此参数，如：`/?token=abc`
```

# 界面预览

![](http://git.oschina.net/zaaksam/dproxy/attach_files/download?i=89952&u=http%3A%2F%2Ffiles.git.oschina.net%2Fgroup1%2FM00%2F01%2FA6%2FPaAvDFmJdUiANn77AAEfwxkjjsI886.png%3Ftoken%3D49dc2684cf356f1a51f522042e86245a%26ts%3D1502533216%26attname%3Dportmap.png)

![](http://git.oschina.net/zaaksam/dproxy/attach_files/download?i=89953&u=http%3A%2F%2Ffiles.git.oschina.net%2Fgroup1%2FM00%2F01%2FA6%2FPaAvDFmJdViAOjEfAAC0dKxZmo4230.png%3Ftoken%3D848e73a93bac218decca2378fb08a0bb%26ts%3D1502533216%26attname%3Dwhitelist.png)

![](http://git.oschina.net/zaaksam/dproxy/attach_files/download?i=89954&u=http%3A%2F%2Ffiles.git.oschina.net%2Fgroup1%2FM00%2F01%2FA6%2FPaAvDFmJdWOABwLKAAEDyxfwF6c412.png%3Ftoken%3Daf7eb35b608885c4b448261dd7610bb7%26ts%3D1502533216%26attname%3Dlog.png)

![](http://git.oschina.net/zaaksam/dproxy/attach_files/download?i=89955&u=http%3A%2F%2Ffiles.git.oschina.net%2Fgroup1%2FM00%2F01%2FA6%2FPaAvDFmJdW-AUgTEAAC0zMqe6dI445.png%3Ftoken%3Daf6f3fd0c2065e82efac837617f55bd6%26ts%3D1502533216%26attname%3Ddoc.png)

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