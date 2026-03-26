# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 构建和运行命令

```bash
# 安装依赖
go mod download

# 开发模式运行
go run cmd/api/main.go

# 生产环境编译
go build -o bin/api cmd/api/main.go

# 运行所有测试
go test ./...

# 运行测试并查看覆盖率
go test -cover ./...

# 运行指定包的测试
go test ./internal/service/...

# 使用热重载工具（开发环境推荐）
go install github.com/cosmtrek/air@latest
air
```

## 架构概述

本项目使用 Gin 框架，采用经典的**三层架构**：

```
Handler (控制器层) → Service (业务逻辑层) → Repository (数据访问层)
```

### 各层职责

- **Router** (`internal/api/router/`): 定义 API 端点和中间件链
- **Handler** (`internal/api/handler/`): HTTP 请求/响应处理，参数绑定，调用 Service
- **Service** (`internal/service/`): 业务逻辑，数据验证，协调 Repository 调用
- **Repository** (`internal/repository/`): 使用 GORM 进行数据库 CRUD 操作
- **Domain** (`internal/domain/`):
  - `model/`: 数据库模型（GORM 结构体）
  - `dto/`: API 请求/响应的数据传输对象
  - `constant/`: 常量和错误码

### 依赖流向

```
Router → Handler → Service → Repository → Database
           ↓         ↓          ↓
         DTO      Model      GORM
```

**规则：**
- Handler 不直接调用 Repository
- Service 不直接操作数据库
- Repository 只处理数据访问，不包含业务逻辑

## 项目结构

```
go-web/
├── cmd/api/main.go          # 应用入口，依赖注入
├── configs/config.yaml      # 配置文件
├── internal/
│   ├── api/
│   │   ├── router/          # 路由定义
│   │   ├── handler/         # HTTP 处理器
│   │   └── middleware/      # 认证中间件等
│   ├── service/             # 业务逻辑层
│   ├── repository/          # 数据访问层
│   ├── domain/
│   │   ├── model/           # 数据库模型
│   │   ├── dto/             # 请求/响应 DTO
│   │   └── constant/        # 常量定义
│   ├── config/              # 配置加载（viper）
│   └── utils/               # 工具函数（JWT、密码加密等）
├── pkg/
│   ├── db/                  # 数据库连接
│   └── logger/              # Zap 日志配置
└── docs/                    # 文档
```

## 开发约定

### 新增功能步骤

添加新功能（如"商品"模块）时，按以下顺序进行：

1. **Model** (`internal/domain/model/product.go`): 定义数据库结构体
2. **DTO** (`internal/domain/dto/product.go`): 定义请求/响应结构
3. **Repository** (`internal/repository/product_repository.go`): 定义接口并实现数据访问
4. **Service** (`internal/service/product_service.go`): 定义接口并实现业务逻辑
5. **Handler** (`internal/api/handler/product_handler.go`): 实现 HTTP 处理器
6. **Router** (`internal/api/router/product.go`): 注册路由
7. **注入依赖** 在 `cmd/api/main.go` 中：初始化 repository → service → handler，传递给 router

### 统一响应格式

所有 API 使用统一的响应结构：

```go
// 成功响应
dto.Success(data)                    // {"code": 200, "message": "success", "data": ...}
dto.SuccessWithMessage("消息", data)  // {"code": 200, "message": "消息", "data": ...}

// 错误响应
dto.Error(400, "错误信息")            // {"code": 400, "message": "错误信息", "data": null}
```

### 认证机制

- 使用 JWT Token，相关代码在 `internal/utils/jwt.go`
- 受保护的路由使用 `middleware.AuthMiddleware()` 中间件
- 在 Handler 中获取当前用户 ID：`userID := utils.GetUserID(c)`

### 配置

配置文件：`configs/config.yaml`（通过 viper 加载）

主要配置项：
- `server`: 端口（默认 8081）、模式（debug/release）、超时时间
- `database`: MySQL 连接（支持 DSN 或单独字段配置，单独字段配置是当前使用的）
- `jwt`: 密钥和过期时间（7天）

**注意**：数据库配置默认使用单独字段（host、port、user 等），DSN 方式被注释掉。

### 数据库

- 使用 GORM 连接 MySQL
- 启动时自动迁移表结构（在 `main.go` 中配置）
- 表根据 Model 结构体自动创建
- 默认数据库名：`dragon`

### CORS 配置

- 前端地址硬编码在 `internal/api/router/router.go` 中：`http://localhost:9000`
- 如需修改，编辑 `router.go` 中的 `AllowOrigins` 配置

## 可用 Skills

项目提供以下自定义开发辅助命令：

- `/create-crud`: 快速创建完整的三层架构 CRUD 功能（推荐使用）
  - 用法：`/create-crud 模块名`（如：`/create-crud product` 或 `/create-crud 商品`）
  - 自动生成：Model、DTO、Repository、Service、Handler、Router
  - 包含完整的增删改查接口
- `/create-api`: 创建新的 API 接口，遵循三层架构
- `/full-crud`: 生成完整的 CRUD 功能（模型、DTO、Repository、Service、Handler、路由）
- `/add-config`: 在系统设置页面新增配置项

## 文档

详细文档位于 `docs/` 目录：
- `架构说明.md`: 架构详解
- `运行指南.md`: 运行说明
- `调试指南.md`: 调试指南
- `请求参数说明.md`: 请求参数处理说明
