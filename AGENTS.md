# Repository Guidelines

## 项目结构与模块组织

本仓库是 Go 模块 `awesomeProject`。入口程序位于 `main.go`。Ent 数据模型定义放在 `ent/schema/`，例如 `ent/schema/todo.go`；`ent/` 下其余大部分文件由 Ent 生成，包含客户端、查询构造器、迁移和实体代码。项目文档及开发规范位于 `docs/`。修改数据模型时，应编辑 schema 后重新生成代码，不要直接修改带有 `Code generated ... DO NOT EDIT.` 标记的文件。当前仓库没有独立的资源目录或测试目录。

## 开发规范（必须遵循）

开始分析、设计或修改 Go 代码前，必须先完整阅读 `docs/go_development_spec.md`。所有新增及修改的手写 Go 代码必须遵循其中的项目结构、命名、注释、分层、错误处理、日志和配置规范；自动生成代码及第三方代码不适用。开发完成后，应按该文档逐项自查，并执行本指南列出的格式化、测试和静态检查命令。若规范与用户在当前任务中的明确要求冲突，以用户要求为准，并在交付说明中指出差异。

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

现有历史采用 Conventional Commits 风格。提交标题应简洁并使用祈使语气，例如 `feat: add todo fields`、`fix: handle empty input`。拉取请求应说明变更目的、主要实现和验证命令，并关联相关 issue。若 schema 有变化，应同时提交 schema 与重新生成的 Ent 文件，并明确迁移影响。
