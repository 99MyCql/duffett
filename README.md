# Duffett

## Introduction

综合训练课设，高频交易分析系统后端。

前端见：[duffett_frontend](https://github.com/99MyCql/duffett_frontend)

## Quickstart

### Prerequisites

语言版本：Go1.14

安装 goimports 工具（执行用户上传的代码时所用）：

```
go get -u golang.org/x/tools/cmd/goimports
```

安装 swag 生成工具：

```
go get -u github.com/swaggo/swag/cmd/swag
```

生成 swagger 文档：

```cmd
cd duffett
swag init
```

在根目录下创建配置文件 `conf.yaml` ，内容如下：

```yaml
# 运行时的 IP 地址和端口号
addr: 0.0.0.0:8080
# MySQL 数据库 url
mysqlUrl: username:password@tcp(ip:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
# 日志文件路径（为空则输出到控制台，本地运行为空即可）
logPath:
# jwt 密钥
jwtSecret: somethingyoulike
# Tushare 社区（https://waditu.com/）获取数据所需的 token
tushareToken: xxxxxx
```

自动生成数据库表：

```cmd
cd duffett
go run scripts\migrateDB.go
```

### Run

运行：

```cmd
cd duffett
go run main.go
```

or

```cmd
cd duffett
go build
.\duffett
```

注：需在代码目录下运行程序，因为执行用户上传的代码时，需调用本项目代码。
