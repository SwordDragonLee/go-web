---
name: create-crud
description: 快速创建完整的三层架构 CRUD 功能（Model、DTO、Repository、Service、Handler、Router）
---

# 创建 CRUD 功能

你是一个 Go-Gin 项目的代码生成助手。用户会提供一个模块名称（如 "product"、"order"），你需要按照项目的三层架构规范，自动生成完整的 CRUD 功能代码。

## 项目架构规范

### 三层架构
```
Handler (控制器层) → Service (业务逻辑层) → Repository (数据访问层)
```

### 文件位置
- Model: `internal/domain/model/{module}.go`
- DTO: `internal/domain/dto/{module}.go`
- Repository: `internal/repository/{module}_repository.go`
- Service: `internal/service/{module}_service.go`
- Handler: `internal/api/handler/{module}_handler.go`
- Router: `internal/api/router/{module}.go`

### 统一响应格式
```go
// 成功响应
dto.Success(data)
dto.SuccessWithMessage("消息", data)

// 错误响应
dto.Error(400, "错误信息")
```

## 生成规则

### 1. 模块名称处理
- 用户输入的模块名称可能是英文或中文
- 如果是中文，转换为对应的英文（如："用户" → "user"，"商品" → "product"）
- 使用 camelCase 命名函数和变量
- 使用 PascalCase 命名类型和文件名

### 2. 数据库模型（Model）
- 创建基础字段：ID、创建时间、更新时间
- 根据模块名称添加合理的业务字段
- 使用 GORM 标签
- 示例字段类型：
  - 字符串：`gorm:"size:255;not null"`
  - 数字：`gorm:"not null"`
  - 时间：`gorm:"type:datetime"`
  - 布尔：`gorm:"type:tinyint(1);default:1"`

### 3. DTO 结构
需要创建以下 DTO：
- `Create{Module}Request` - 创建请求
- `Update{Module}Request` - 更新请求
- `{Module}Response` - 响应结构
- `List{Module}Request` - 列表查询请求（包含分页）
- `{Module}ListResponse` - 列表响应

### 4. Repository 接口和实现
定义接口：
```go
type {Module}Repository interface {
    Create(ctx context.Context, req *dto.Create{Module}Request) error
    Update(ctx context.Context, id uint, req *dto.Update{Module}Request) error
    Delete(ctx context.Context, id uint) error
    GetByID(ctx context.Context, id uint) (*model.{Module}, error)
    List(ctx context.Context, req *dto.List{Module}Request) ([]*model.{Module}, int64, error)
}
```

### 5. Service 接口和实现
- 包含业务逻辑验证
- 调用 Repository 进行数据操作
- 返回 DTO 格式的数据

### 6. Handler 实现
注册路由：
```go
POST   /api/{module}        - 创建
GET    /api/{module}/:id    - 获取详情
PUT    /api/{module}/:id    - 更新
DELETE /api/{module}/:id    - 删除
GET    /api/{module}        - 列表查询
```

### 7. Router 注册
在 `internal/api/router/router.go` 中添加路由组

## 代码模板要点

### Model 示例
```go
package model

import (
    "time"
    "gorm.io/gorm"
)

type {Module} struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    Name      string         `gorm:"size:100;not null" json:"name"`
    // 其他业务字段...
    Status    int            `gorm:"type:tinyint(1);default:1" json:"status"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func ({Module}) TableName() string {
    return "{module}s"
}
```

### Repository 示例
```go
func (r *{Module}Repository) Create(ctx context.Context, req *dto.Create{Module}Request) error {
    {module} := &model.{Module}{
        Name: req.Name,
        // 映射其他字段...
    }
    return r.db.WithContext(ctx).Create({module}).Error
}
```

### Service 示例
```go
func (s *{Module}Service) Create(ctx context.Context, req *dto.Create{Module}Request) error {
    // 业务逻辑验证
    if err := s.repo.Create(ctx, req); err != nil {
        return err
    }
    return nil
}
```

### Handler 示例
```go
func (h *{Module}Handler) Create(c *gin.Context) {
    var req dto.Create{Module}Request
    if err := c.ShouldBindJSON(&req); err != nil {
        dto.Error(400, "参数错误: "+err.Error())
        return
    }

    if err := h.service.Create(c.Request.Context(), &req); err != nil {
        dto.Error(500, "创建失败")
        return
    }

    dto.SuccessWithMessage("创建成功", nil)
}
```

## 执行步骤

当用户提供模块名称时：

1. **确认信息**：告诉用户将要创建的模块名称和功能
2. **创建 Model**：生成数据库模型文件
3. **创建 DTO**：生成请求/响应 DTO 文件
4. **创建 Repository**：生成数据访问层
5. **创建 Service**：生成业务逻辑层
6. **创建 Handler**：生成控制器层
7. **创建 Router**：生成路由注册
8. **更新 main.go**：提示用户在 main.go 中添加依赖注入
9. **更新数据库迁移**：提示用户在 main.go 的 AutoMigrate 中添加新模型

## 注意事项

- 严格遵循项目的三层架构规范
- 使用项目的统一响应格式
- 添加适当的错误处理
- 添加必要的注释（中文）
- 代码格式化（使用 gofmt）
- 不要覆盖已存在的文件（如果文件已存在，提示用户）

## 开始生成

用户会输入模块名称，例如：
- "product" 或 "商品"
- "order" 或 "订单"
- "article" 或 "文章"

立即开始生成完整的 CRUD 代码！
