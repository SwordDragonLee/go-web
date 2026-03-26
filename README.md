# Go-Web 后端项目

基于 Gin 框架的 RESTful API 后端服务，采用经典三层架构设计。

## ✨ 特性

- 🏗️ **三层架构** - Handler → Service → Repository 清晰分层
- 🔐 **JWT 认证** - 完整的用户认证授权机制
- 📝 **RBAC 权限** - 基于角色的访问控制
- 🛡️ **中间件支持** - 认证、CORS、日志等
- 📊 **自动迁移** - GORM 自动同步数据库表结构
- 🔥 **热重载** - 开发环境支持 Air 热更新

## 🏗️ 项目架构

```
┌─────────────┐      ┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│   Router    │ ───▶ │   Handler   │ ───▶ │   Service   │ ───▶ │ Repository  │
│  (路由层)    │      │  (控制器层)  │      │  (业务逻辑层) │      │ (数据访问层) │
└─────────────┘      └─────────────┘      └─────────────┘      └─────────────┘
                                                           │
                                                           ▼
                                                    ┌─────────────┐
                                                    │  Database   │
                                                    │  (MySQL)    │
                                                    └─────────────┘
```

### 目录结构

```
go-web/
├── cmd/api/main.go           # 应用入口
├── configs/config.yaml       # 配置文件
├── internal/
│   ├── api/                  # API 层
│   │   ├── handler/          # HTTP 处理器
│   │   ├── router/           # 路由定义
│   │   └── middleware/       # 中间件
│   ├── service/              # 业务逻辑层
│   ├── repository/           # 数据访问层
│   ├── domain/               # 领域模型
│   │   ├── model/            # 数据库模型
│   │   ├── dto/              # 数据传输对象
│   │   └── constant/         # 常量定义
│   ├── config/               # 配置加载
│   └── utils/                # 工具函数
├── pkg/                      # 公共包
│   ├── db/                   # 数据库连接
│   └── logger/               # 日志配置
└── docs/                     # 文档
```

## 🚀 快速开始

### 环境要求

- Go 1.21+
- MySQL 8.0+

### 安装步骤

1. **克隆仓库**
```bash
git clone https://github.com/SwordDragonLee/go-web.git
cd go-web
```

2. **安装依赖**
```bash
go mod download
```

3. **配置数据库**

编辑 `configs/config.yaml`：
```yaml
database:
  host: localhost
  port: 3306
  user: root
  password: your_password
  database: dragon
```

4. **创建数据库**
```sql
CREATE DATABASE dragon CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

5. **运行项目**
```bash
# 开发模式
go run cmd/api/main.go

# 或使用 Air 热重载（推荐）
air
```

6. **访问服务**

API 地址：http://localhost:8081

## 📦 可用功能

### 用户管理
- 用户注册
- 用户登录（JWT Token）
- 发送验证码
- 修改密码
- 查看用户信息

### 角色权限
- 角色管理
- 权限管理
- 用户角色分配
- 角色权限分配

## 🔧 开发工具

### 运行测试
```bash
# 运行所有测试
go test ./...

# 查看测试覆盖率
go test -cover ./...

# 运行指定包测试
go test ./internal/service/...
```

### 热重载开发
```bash
# 安装 Air
go install github.com/cosmtrek/air@latest

# 运行
air
```

### 生产构建
```bash
go build -o bin/api cmd/api/main.go
```

## 📖 API 文档

### 统一响应格式

**成功响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": { ... }
}
```

**错误响应：**
```json
{
  "code": 400,
  "message": "错误信息",
  "data": null
}
```

### 主要接口

- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `GET /api/user/info` - 获取用户信息（需认证）
- `POST /api/user/send-code` - 发送验证码
- `POST /api/user/reset-password` - 重置密码

详细文档请查看 `docs/` 目录。

## 🔐 认证机制

使用 JWT Token 进行用户认证：

1. 登录成功后返回 Token
2. 后续请求在 Header 中携带：
```bash
Authorization: Bearer <token>
```

Token 默认有效期为 7 天。

## 🛠️ 技术栈

- **框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **日志**: Zap
- **配置**: Viper
- **认证**: JWT

## 📄 License

MIT License

## 👤 作者

SwordDragonLee

---

**注意**: 本项目仅供学习和参考使用，请勿直接用于生产环境。
