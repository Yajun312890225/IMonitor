<h1 align="center">iMonitor</h1>

<div align="center">基于Gin + Casbin + Grom 的IM后台管理系统<br/></div>

## 特征
- 遵循 `RESTful API` 设计规范 & 基于接口的编程规范
- 基于 `GIN` 框架，提供了丰富的中间件支持
- 基于 `Gorilla/Websocket` 实现实时日志
- 基于 `Casbin` 的 RBAC 访问控制模型
- 基于 `Gorm` 的数据库存储
- 基于 `Logrus` 实现了日志输出
- 基于 `Swaggo` 自动生成 `Swagger` 文档
- 基于 `go mod` 的依赖管理

## 依赖工具

```
go get -u github.com/swaggo/swag/cmd/swag
```
- [swag](https://github.com/swaggo/swag) -- Automatically generate RESTful API documentation with Swagger 2.0 for Go.

## 依赖框架

- [Gin](https://gin-gonic.com/) -- The fastest full-featured web framework for Go.
- [Gorilla/Websocket](https://pkg.go.dev/github.com/gorilla/websocket) -- Gorilla WebSocket is a Go implementation of the WebSocket protocol.
- [GORM](http://gorm.io/) -- The fantastic ORM library for Golang
- [Casbin](https://casbin.org/) -- An authorization library that supports access control models like ACL, RBAC, ABAC in Golang


## 使用 `go mod` 进行依赖管理
```
#Go版本需要大于1.13 
go env -w GO111MODULE=on 
#如果没有科学上网
go env -w GOPROXY=https://goproxy.cn,direct
#go mod 命令
go mod download    下载依赖的module到本地cache（默认为$GOPATH/pkg/mod目录）
go mod edit        编辑go.mod文件
go mod graph       打印模块依赖图
go mod init        初始化当前文件夹, 创建go.mod文件
go mod tidy        增加缺少的module，删除无用的module
go mod vendor      将依赖复制到vendor下
go mod verify      校验依赖
go mod why         解释为什么需要依赖
```
## 生成`swagger`文档

```
# 直接运行
make
# 直接编译linux
make build
# 基于Makefile
make swagger
# 使用swag命令
swag init
```