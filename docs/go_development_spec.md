# Go 开发规范手册

> **作者**：Seven Wu  
> **版本**：v1.0.0  
> **更新日期**：2026-07-02  
> **适用范围**：所有 Go 语言后端项目

---

## 目录

1. [项目结构规范](#1-项目结构规范)
2. [命名规范](#2-命名规范)
3. [注释规范](#3-注释规范)
4. [包（Package）规范](#4-包package规范)
5. [函数与方法规范](#5-函数与方法规范)
6. [控制器（Controller）规范](#6-控制器controller规范)
7. [Service 层规范](#7-service-层规范)
8. [Repository / DAO 层规范](#8-repository--dao-层规范)
9. [Model / Entity 规范](#9-model--entity-规范)
10. [错误处理规范](#10-错误处理规范)
11. [日志规范](#11-日志规范)
12. [配置规范](#12-配置规范)
13. [完整示例代码](#13-完整示例代码)

---

## 1. 项目结构规范

标准目录布局如下，基于 [golang-standards/project-layout](https://github.com/golang-standards/project-layout)：

```
myproject/
├── cmd/                    # 程序入口（main 包）
│   └── server/
│       └── main.go
├── internal/               # 项目内部代码（外部无法导入）
│   ├── controller/         # 控制器层：处理 HTTP 请求/响应
│   ├── service/            # 业务逻辑层
│   ├── repository/         # 数据访问层（DAO）
│   ├── model/              # 数据模型（实体、DTO、VO）
│   ├── middleware/         # 中间件
│   └── router/             # 路由注册
├── pkg/                    # 可被外部复用的公共库
│   ├── logger/             # 日志封装
│   ├── response/           # 统一响应封装
│   └── errors/             # 自定义错误
├── config/                 # 配置文件（yaml/json）
├── docs/                   # 项目文档、Swagger
├── scripts/                # 脚本文件
├── go.mod
├── go.sum
└── README.md
```

---

## 2. 命名规范

### 2.1 文件命名

| 类型 | 规则 | 示例 |
|------|------|------|
| 普通文件 | 小写 + 下划线 | `user_service.go` |
| 测试文件 | 文件名 + `_test` | `user_service_test.go` |
| 接口文件 | 功能名 + `_interface` | `user_repository_interface.go` |

### 2.2 变量命名

```go
// ✅ 正确：驼峰命名，首字母小写（包内私有）
userID := 1
userName := "Seven Wu"

// ✅ 正确：首字母大写（包外公开）
UserName := "Seven Wu"

// ❌ 错误：下划线命名（非常量）
user_name := "Seven Wu"
```

### 2.3 常量命名

```go
// ✅ 全大写 + 下划线，表示常量
const (
    MAX_RETRY_COUNT = 3
    DEFAULT_TIMEOUT = 30
)

// ✅ 使用 iota 的枚举常量
type UserStatus int

const (
    UserStatusActive   UserStatus = iota + 1 // 激活
    UserStatusInactive                        // 禁用
    UserStatusDeleted                         // 已删除
)
```

### 2.4 接口命名

```go
// ✅ 单方法接口：方法名 + "er"
type Reader interface {
    Read(p []byte) (n int, err error)
}

// ✅ 多方法接口：功能名 + "Service" / "Repository"
type UserService interface {
    GetByID(ctx context.Context, id int64) (*model.User, error)
    Create(ctx context.Context, req *model.CreateUserReq) error
}
```

---

## 3. 注释规范

> 注释是代码的"说明书"，所有公开的包、类型、函数、方法必须有注释。

### 3.1 文件头注释

每个 `.go` 文件顶部必须包含文件头注释，格式如下：

```go
// Package controller 处理用户相关的 HTTP 请求。
//
// 包含用户注册、登录、信息查询等控制器方法。
//
// Author: Seven Wu
// Date: 2026-07-02
// Version: 1.0.0
package controller
```

### 3.2 函数/方法注释

遵循 Go 官方 `godoc` 规范：**注释第一行必须以函数名开头**。

```go
// GetUserByID 根据用户 ID 查询用户信息。
//
// 参数:
//   - ctx: 请求上下文，用于超时和取消控制
//   - id:  用户唯一标识符，必须大于 0
//
// 返回:
//   - *model.User: 用户实体对象
//   - error: 查询失败时返回错误信息
//
// Author: Seven Wu
func (s *userService) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
    // ...
}
```

### 3.3 结构体注释

```go
// User 表示系统中的用户实体。
//
// Author: Seven Wu
type User struct {
    ID        int64     `json:"id" gorm:"primaryKey"`     // 用户唯一 ID
    Username  string    `json:"username" gorm:"not null"` // 用户名（唯一）
    Email     string    `json:"email"`                    // 用户邮箱
    Status    int       `json:"status"`                   // 状态：1=激活 2=禁用
    CreatedAt time.Time `json:"created_at"`               // 创建时间
    UpdatedAt time.Time `json:"updated_at"`               // 更新时间
}
```

### 3.4 行内注释

```go
func processOrder(order *Order) error {
    // 检查库存是否充足，不足则提前返回错误
    if order.Stock < order.Quantity {
        return errors.New("库存不足")
    }

    total := order.Price * float64(order.Quantity) // 计算订单总金额
    order.Total = total

    return nil
}
```

### 3.5 TODO / FIXME 注释

```go
// TODO(Seven Wu): 后续需要增加缓存层，减少数据库查询压力
// FIXME(Seven Wu): 当并发量过高时此处存在竞态条件，需加锁处理
// NOTE: 此逻辑与产品文档 PRD-2026-007 第 3.2 节对应
```

---

## 4. 包（Package）规范

```go
// ✅ 正确：包名全小写，简短，无下划线
package service
package repository
package middleware

// ❌ 错误
package userService   // 不用驼峰
package user_service  // 不用下划线
```

---

## 5. 函数与方法规范

### 5.1 函数参数规范

- 参数不超过 **5 个**，超出时使用结构体封装
- 上下文 `context.Context` 必须作为**第一个参数**

```go
// ✅ 正确：context 作为第一个参数
func CreateUser(ctx context.Context, req *CreateUserReq) (*User, error) {}

// ❌ 错误：context 不在第一位
func CreateUser(req *CreateUserReq, ctx context.Context) (*User, error) {}
```

### 5.2 函数返回值规范

- 错误值 `error` 必须作为**最后一个返回值**
- 避免使用命名返回值（除非能显著提升可读性）

```go
// ✅ 推荐
func GetUser(ctx context.Context, id int64) (*User, error) {}

// ⚠️ 仅在必要时使用命名返回值
func divide(a, b float64) (result float64, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("除法发生 panic: %v", r)
        }
    }()
    result = a / b
    return
}
```

### 5.3 函数职责单一

每个函数只做一件事，行数建议不超过 **80 行**。

---

## 6. 控制器（Controller）规范

控制器负责：接收 HTTP 请求 → 参数校验 → 调用 Service → 返回响应。

**不允许**在 Controller 层写业务逻辑。

### 6.1 Controller 文件示例

```go
// Package controller 处理用户相关的 HTTP 请求。
//
// 职责：参数绑定与校验、调用 Service 层、返回统一响应。
// 不包含任何业务逻辑。
//
// Author: Seven Wu
// Date: 2026-07-02
package controller

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "myproject/internal/model"
    "myproject/internal/service"
    "myproject/pkg/response"
)

// UserController 用户控制器，负责处理用户相关的 HTTP 接口。
//
// Author: Seven Wu
type UserController struct {
    userSvc service.UserService // 用户业务逻辑层
}

// NewUserController 创建并返回 UserController 实例。
//
// 参数:
//   - userSvc: 用户业务逻辑层接口
//
// Author: Seven Wu
func NewUserController(userSvc service.UserService) *UserController {
    return &UserController{userSvc: userSvc}
}

// GetUser 获取用户详情接口。
//
// @Summary     获取用户详情
// @Description 根据用户 ID 查询用户信息
// @Tags        用户管理
// @Accept      json
// @Produce     json
// @Param       id   path     int           true "用户 ID"
// @Success     200  {object} response.R{data=model.UserVO}
// @Failure     400  {object} response.R
// @Failure     500  {object} response.R
// @Router      /api/v1/users/{id} [get]
//
// Author: Seven Wu
func (c *UserController) GetUser(ctx *gin.Context) {
    // 解析路径参数
    var req model.GetUserReq
    if err := ctx.ShouldBindUri(&req); err != nil {
        response.Fail(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
        return
    }

    // 调用 Service 层获取用户信息
    user, err := c.userSvc.GetByID(ctx.Request.Context(), req.ID)
    if err != nil {
        response.Fail(ctx, http.StatusInternalServerError, err.Error())
        return
    }

    response.Success(ctx, user)
}

// CreateUser 创建用户接口。
//
// @Summary     创建用户
// @Description 注册新用户账号
// @Tags        用户管理
// @Accept      json
// @Produce     json
// @Param       body body     model.CreateUserReq true "创建用户请求体"
// @Success     200  {object} response.R{data=model.UserVO}
// @Failure     400  {object} response.R
// @Failure     500  {object} response.R
// @Router      /api/v1/users [post]
//
// Author: Seven Wu
func (c *UserController) CreateUser(ctx *gin.Context) {
    // 绑定并校验请求体参数
    var req model.CreateUserReq
    if err := ctx.ShouldBindJSON(&req); err != nil {
        response.Fail(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
        return
    }

    // 调用 Service 层创建用户
    user, err := c.userSvc.Create(ctx.Request.Context(), &req)
    if err != nil {
        response.Fail(ctx, http.StatusInternalServerError, err.Error())
        return
    }

    response.Success(ctx, user)
}
```

---

## 7. Service 层规范

Service 层负责：业务逻辑处理、调用 Repository、事务管理。

### 7.1 接口定义

```go
// Package service 提供用户相关的业务逻辑处理。
//
// Author: Seven Wu
// Date: 2026-07-02
package service

import (
    "context"

    "myproject/internal/model"
)

// UserService 用户业务逻辑层接口。
//
// 所有用户相关业务操作均通过此接口定义，便于单元测试时 Mock。
//
// Author: Seven Wu
type UserService interface {
    // GetByID 根据 ID 获取用户信息。
    GetByID(ctx context.Context, id int64) (*model.UserVO, error)

    // Create 创建新用户。
    Create(ctx context.Context, req *model.CreateUserReq) (*model.UserVO, error)

    // Update 更新用户信息。
    Update(ctx context.Context, id int64, req *model.UpdateUserReq) error

    // Delete 删除用户（软删除）。
    Delete(ctx context.Context, id int64) error
}
```

### 7.2 接口实现

```go
// userServiceImpl 是 UserService 接口的实现。
//
// Author: Seven Wu
type userServiceImpl struct {
    userRepo repository.UserRepository // 用户数据访问层
}

// NewUserService 创建并返回 UserService 实例。
//
// 参数:
//   - userRepo: 用户数据访问层接口
//
// Author: Seven Wu
func NewUserService(userRepo repository.UserRepository) UserService {
    return &userServiceImpl{userRepo: userRepo}
}

// GetByID 根据用户 ID 查询用户信息。
//
// 业务规则：
//  1. ID 必须大于 0
//  2. 用户不存在时返回自定义 NotFound 错误
//
// Author: Seven Wu
func (s *userServiceImpl) GetByID(ctx context.Context, id int64) (*model.UserVO, error) {
    if id <= 0 {
        return nil, errors.New("用户 ID 无效")
    }

    user, err := s.userRepo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("查询用户失败: %w", err)
    }
    if user == nil {
        return nil, errors.New("用户不存在")
    }

    // 将 Entity 转换为 VO（视图对象）
    return model.ToUserVO(user), nil
}
```

---

## 8. Repository / DAO 层规范

Repository 层负责：数据库 CRUD 操作，不包含业务逻辑。

```go
// Package repository 提供用户数据的持久化操作。
//
// Author: Seven Wu
// Date: 2026-07-02
package repository

import (
    "context"

    "gorm.io/gorm"

    "myproject/internal/model"
)

// UserRepository 用户数据访问层接口。
//
// Author: Seven Wu
type UserRepository interface {
    // FindByID 根据主键查询用户。
    FindByID(ctx context.Context, id int64) (*model.User, error)

    // FindByEmail 根据邮箱查询用户。
    FindByEmail(ctx context.Context, email string) (*model.User, error)

    // Create 插入新用户记录。
    Create(ctx context.Context, user *model.User) error

    // Update 更新用户信息。
    Update(ctx context.Context, user *model.User) error

    // Delete 软删除用户记录。
    Delete(ctx context.Context, id int64) error
}

// userRepository 是 UserRepository 接口的 GORM 实现。
//
// Author: Seven Wu
type userRepository struct {
    db *gorm.DB // 数据库连接
}

// NewUserRepository 创建并返回 UserRepository 实例。
//
// 参数:
//   - db: GORM 数据库连接实例
//
// Author: Seven Wu
func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

// FindByID 根据主键查询用户记录。
//
// 返回:
//   - *model.User: 用户实体，不存在时返回 nil
//   - error: 数据库错误
//
// Author: Seven Wu
func (r *userRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
    var user model.User
    result := r.db.WithContext(ctx).First(&user, id)
    if result.Error != nil {
        // 记录不存在不视为错误，返回 nil
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, fmt.Errorf("FindByID 查询失败: %w", result.Error)
    }
    return &user, nil
}
```

---

## 9. Model / Entity 规范

```go
// Package model 定义项目中所有的数据模型。
//
// 包含三类模型：
//   - Entity: 与数据库表对应的实体
//   - Request (Req): 接收前端请求的 DTO
//   - View Object (VO): 返回给前端的视图对象
//
// Author: Seven Wu
// Date: 2026-07-02
package model

import "time"

// User 用户实体，对应数据库 users 表。
//
// Author: Seven Wu
type User struct {
    ID        int64     `gorm:"primaryKey;autoIncrement"` // 用户 ID
    Username  string    `gorm:"size:50;not null;unique"`  // 用户名
    Email     string    `gorm:"size:100;not null;unique"` // 邮箱
    Password  string    `gorm:"size:200;not null"`        // 加密后的密码
    Status    int8      `gorm:"default:1"`                // 状态：1=正常 2=禁用
    CreatedAt time.Time `gorm:"autoCreateTime"`           // 创建时间
    UpdatedAt time.Time `gorm:"autoUpdateTime"`           // 更新时间
}

// TableName 返回用户实体对应的数据库表名。
//
// Author: Seven Wu
func (User) TableName() string {
    return "users"
}

// CreateUserReq 创建用户的请求参数。
//
// Author: Seven Wu
type CreateUserReq struct {
    Username string `json:"username" binding:"required,min=3,max=50"` // 用户名，3-50个字符
    Email    string `json:"email" binding:"required,email"`            // 邮箱地址
    Password string `json:"password" binding:"required,min=8"`         // 密码，最少8位
}

// GetUserReq 获取用户详情的请求参数。
//
// Author: Seven Wu
type GetUserReq struct {
    ID int64 `uri:"id" binding:"required,gt=0"` // 用户 ID，必须大于 0
}

// UserVO 用户视图对象，返回给前端的用户信息。
//
// 不包含敏感字段（如密码）。
//
// Author: Seven Wu
type UserVO struct {
    ID        int64  `json:"id"`         // 用户 ID
    Username  string `json:"username"`   // 用户名
    Email     string `json:"email"`      // 邮箱
    Status    int8   `json:"status"`     // 状态
    CreatedAt string `json:"created_at"` // 创建时间（格式化字符串）
}

// ToUserVO 将 User Entity 转换为 UserVO。
//
// Author: Seven Wu
func ToUserVO(u *User) *UserVO {
    return &UserVO{
        ID:        u.ID,
        Username:  u.Username,
        Email:     u.Email,
        Status:    u.Status,
        CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
    }
}
```

---

## 10. 错误处理规范

### 10.1 基本原则

```go
// ✅ 错误必须处理，不允许用 _ 忽略
result, err := doSomething()
if err != nil {
    return fmt.Errorf("doSomething 执行失败: %w", err) // 用 %w 包裹，保留调用链
}

// ❌ 禁止忽略错误
result, _ := doSomething()
```

### 10.2 自定义错误类型

```go
// Package errors 提供项目统一的错误类型定义。
//
// Author: Seven Wu
package errors

import "fmt"

// AppError 业务错误类型，包含错误码和错误消息。
//
// Author: Seven Wu
type AppError struct {
    Code    int    // 业务错误码
    Message string // 错误描述
    Err     error  // 原始错误（可选）
}

// Error 实现 error 接口。
//
// Author: Seven Wu
func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
    }
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 支持 errors.Is / errors.As 调用链追踪。
//
// Author: Seven Wu
func (e *AppError) Unwrap() error {
    return e.Err
}

// 预定义常用错误码
var (
    ErrNotFound     = &AppError{Code: 10001, Message: "资源不存在"}
    ErrUnauthorized = &AppError{Code: 10002, Message: "未授权访问"}
    ErrForbidden    = &AppError{Code: 10003, Message: "无权限操作"}
    ErrBadRequest   = &AppError{Code: 10004, Message: "请求参数错误"}
)
```

---

## 11. 日志规范

```go
// ✅ 使用结构化日志，包含关键上下文字段
logger.Info("用户登录成功",
    zap.Int64("user_id", user.ID),
    zap.String("username", user.Username),
    zap.String("ip", ctx.ClientIP()),
)

// ✅ 错误日志必须附带 error 对象
logger.Error("查询用户失败",
    zap.Int64("user_id", id),
    zap.Error(err),
)

// ❌ 禁止使用 fmt.Println 打印日志
fmt.Println("用户登录")
```

日志级别使用规范：

| 级别 | 使用场景 |
|------|---------|
| `Debug` | 开发调试，生产环境关闭 |
| `Info` | 关键业务节点，如登录、下单 |
| `Warn` | 可预期的异常，如参数校验失败 |
| `Error` | 程序错误，需要告警处理 |
| `Fatal` | 系统无法启动，程序退出 |

---

## 12. 配置规范

```go
// Config 项目全局配置结构体。
//
// 对应 config/config.yaml 文件中的配置项。
//
// Author: Seven Wu
type Config struct {
    App      AppConfig      `yaml:"app"`      // 应用配置
    Database DatabaseConfig `yaml:"database"` // 数据库配置
    Redis    RedisConfig    `yaml:"redis"`    // Redis 配置
    Log      LogConfig      `yaml:"log"`      // 日志配置
}

// AppConfig 应用基础配置。
//
// Author: Seven Wu
type AppConfig struct {
    Name string `yaml:"name"` // 应用名称
    Port int    `yaml:"port"` // 监听端口
    Mode string `yaml:"mode"` // 运行模式：debug / release
}
```

---

## 13. 完整示例代码

> 以下是一个完整的 **用户查询接口** 完整调用链示例，从路由到数据库。

```
HTTP 请求
  │
  ▼
Router (路由注册)
  │
  ▼
Middleware (鉴权、日志、限流)
  │
  ▼
Controller (参数绑定与校验)
  │
  ▼
Service (业务逻辑)
  │
  ▼
Repository (数据库操作)
  │
  ▼
Database (MySQL / PostgreSQL)
```

### router/router.go

```go
// Package router 注册项目所有 HTTP 路由。
//
// Author: Seven Wu
// Date: 2026-07-02
package router

import (
    "github.com/gin-gonic/gin"

    "myproject/internal/controller"
    "myproject/internal/middleware"
)

// Setup 初始化并返回 Gin 路由引擎。
//
// Author: Seven Wu
func Setup(userCtrl *controller.UserController) *gin.Engine {
    r := gin.New()

    // 注册全局中间件
    r.Use(middleware.Logger())   // 请求日志
    r.Use(middleware.Recovery()) // panic 恢复

    // API v1 路由组
    v1 := r.Group("/api/v1")
    {
        users := v1.Group("/users")
        users.Use(middleware.Auth()) // 需要鉴权
        {
            users.GET("/:id", userCtrl.GetUser)       // 获取用户详情
            users.POST("", userCtrl.CreateUser)        // 创建用户
        }
    }

    return r
}
```

### cmd/server/main.go

```go
// Package main 是程序入口，负责初始化依赖并启动 HTTP 服务。
//
// Author: Seven Wu
// Date: 2026-07-02
package main

import (
    "fmt"
    "log"

    "myproject/internal/controller"
    "myproject/internal/repository"
    "myproject/internal/router"
    "myproject/internal/service"
)

// main 程序入口函数。
//
// 启动顺序：加载配置 → 初始化数据库 → 依赖注入 → 启动 HTTP 服务器。
//
// Author: Seven Wu
func main() {
    // 1. 加载配置
    cfg := config.Load("config/config.yaml")

    // 2. 初始化数据库
    db := database.Init(cfg.Database)

    // 3. 依赖注入（手动 DI，也可使用 wire）
    userRepo := repository.NewUserRepository(db)
    userSvc  := service.NewUserService(userRepo)
    userCtrl := controller.NewUserController(userSvc)

    // 4. 注册路由
    r := router.Setup(userCtrl)

    // 5. 启动 HTTP 服务
    addr := fmt.Sprintf(":%d", cfg.App.Port)
    log.Printf("服务启动成功，监听地址: %s", addr)
    if err := r.Run(addr); err != nil {
        log.Fatalf("服务启动失败: %v", err)
    }
}
```

---

## 附录：快速注释模板

复制到你的编辑器代码片段（Snippet）中使用：

### 文件头模板

```go
// Package xxx 此包的功能简述。
//
// 详细描述（可选）。
//
// Author: Seven Wu
// Date: 2026-07-02
package xxx
```

### 函数/方法模板

```go
// FuncName 函数功能简述（首行以函数名开头）。
//
// 参数:
//   - ctx: 请求上下文
//   - param: 参数说明
//
// 返回:
//   - *Type: 返回值说明
//   - error: 错误信息
//
// Author: Seven Wu
func FuncName(ctx context.Context, param string) (*Type, error) {
}
```

### 结构体模板

```go
// StructName 结构体功能描述。
//
// Author: Seven Wu
type StructName struct {
    Field1 string // 字段说明
    Field2 int    // 字段说明
}
```

---

*文档维护：Seven Wu | 如有疑问请提 Issue 或联系作者*
