# Duffett

## Quick Start

根目录下创建配置文件 `config.yaml` ，内容如下：

```yaml
# 运行时的 IP 地址和端口号
addr: 0.0.0.0:8080
# MySQL 数据库 url
mysqlUrl: username:password@tcp(ip:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
# 日志文件路径（为空则输出到控制台，本地运行为空即可）
logPath: duffett.log
# jwt 密钥
jwtSecret: somethingyoulike
```

安装 `swag` 工具：

```cmd
go get -u github.com/swaggo/swag/cmd/swag
```

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
