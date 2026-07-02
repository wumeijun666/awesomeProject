# Repository Guidelines

## 项目结构与模块组织

本仓库是 Go 模块 `awesomeProject`。入口程序位于 `main.go`。Ent 数据模型定义放在 `ent/schema/`，例如 `ent/schema/todo.go`；`ent/` 下其余大部分文件由 Ent 生成，包含客户端、查询构造器、迁移和实体代码。修改数据模型时，应编辑 schema 后重新生成代码，不要直接修改带有 `Code generated ... DO NOT EDIT.` 标记的文件。当前仓库没有独立的资源目录或测试目录。

## 构建、测试与本地开发

在仓库根目录使用 PowerShell 执行：

```powershell
go run .
go build ./...
go test ./...
go vet ./...
go generate ./ent
```

`go run .` 启动示例程序；`go build ./...` 编译全部包；`go test ./...` 运行现有及后续新增测试；`go vet ./...` 检查常见错误；修改 `ent/schema/` 后运行 `go generate ./ent` 更新生成代码。依赖发生变化时执行 `go mod tidy`，并同时提交 `go.mod` 和 `go.sum`。

## 编码风格与命名约定

提交前对修改过的 Go 文件运行 `gofmt -w <文件>`，使用 Go 标准制表符缩进和 import 分组。包名使用简短小写单词；导出标识符使用 `PascalCase`，非导出标识符使用 `camelCase`。保持函数职责单一，并为导出的类型和函数添加以标识符开头的 GoDoc 注释。

## 测试规范

使用标准库 `testing`。测试文件与被测代码同目录，命名为 `*_test.go`；测试函数命名为 `TestXxx`。对多个输入优先采用表驱动测试。新增行为或修复缺陷时应补充覆盖正常路径和失败路径的测试，并确保 `go test ./...` 通过。

## 提交与拉取请求

当前工作目录未包含 Git 历史，无法归纳既有约定。建议使用简洁、祈使语气的提交标题，可采用 `feat: add todo fields`、`fix: handle empty input` 等格式。拉取请求应说明变更目的、主要实现和验证命令；关联相关 issue。若 schema 有变化，应同时提交 schema 与重新生成的 Ent 文件，并明确迁移影响。
