# dproxy

简易的网络代理工具，带有IP白名单限制管，带有简洁的UI管理界面，提供丰富的API接口，可方便的与各个系统集成，可编译为单文件运行

# 技术资源

## Backend

语言 Go 1.8.3

Web/API开发框架 Beego 1.9

ORM框架 Xorm 0.6.2

数据库 sqlite3

静态资源打包工具 statik

依赖管理工具 govendor

---

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