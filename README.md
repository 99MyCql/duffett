# Duffett

## Introduction

高频交易系统后端。

前端见：[duffett_frontend](https://github.com/99MyCql/duffett_frontend)

## Quickstart

### Prerequisites

语言版本：Go1.14

下载依赖：

```cmd
go mod download
```

根目录下创建配置文件 `conf.yaml` ，内容如下：

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
go run scripts\migrateDB.go
```

### Run

生成 swagger 文档：

```cmd
swag init
```

运行：

```cmd
go run main.go
```

or

```cmd
go build
.\duffett
```
