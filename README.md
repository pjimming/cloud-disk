# Jimmy-CloudDisk

> 基于 go-zero，xorm 实现的轻量级云盘系统。

需要用到的命令

```text
go mod init cloud-disk

# 安装 xorm
go get xorm.io/xorm

# 1. 把项目所依赖的包添加到go.mod文件中
# 2. 去掉go.mod文件中项目不需要的依赖包。
go mod tidy

# 安装 go-zero
go get -u github.com/zeromicro/go-zero

# 安装 Goctl
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 创建 API 服务
goctl api new core

# 启动服务
go run core.go -f etc/core-api.yaml

# 使用api文件生成代码
goctl api go -api core.api -dir . -style go_zero

# 安装 Email 库
go get github.com/jordan-wright/email
```

[电子邮件库](https://github.com/jordan-wright/email)
