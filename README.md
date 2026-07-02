# README

## run

```bash
go mod tidy: 清理未使用的依赖并添加缺失的依赖。
go mod download: 将依赖包下载到本地缓存。
go get <package>: 添加或升级依赖包。

# 运行
go run main.go
# 构建
go build .
./demo-api
```

## ent自动生成实体

```bash
# 新建实体
go run -mod=mod entgo.io/ent/cmd/ent new User
# 本项目被下面elk 指令替换
# go generate ./ent
```

## elk自动生成CRUD接口代码

```bash
# 首次新建 ent/entc.go
//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate ./schema
# 替换为
//go:generate go run -mod=mod entc.go

# 首次
go get github.com/mailru/easyjson \
       github.com/masseelch/render \
       github.com/go-chi/chi/v5 \
       go.uber.org/zap

# 日常生成代码
go generate ./...
```

## CRUD

```bash
# 列表
curl http://localhost:8080/users

# 创建
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"age": 21, "name": "dlice"}'

# 查询
curl http://localhost:8080/users/1

# 更新（注意用 PATCH）
curl -X PATCH http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"age": 26, "name": "bob"}'

# 删除
curl -X DELETE http://localhost:8080/users/1
```