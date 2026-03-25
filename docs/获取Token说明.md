# 如何获取当前登录用户的 Token

## 当前实现

在认证中间件中，Token 从请求头提取并解析，但**原始 Token 字符串没有存储到 Context 中**。

### 中间件中的 Token 处理

```go
// internal/api/middleware/auth.go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从请求头获取 Token
        authHeader := c.GetHeader("Authorization")
        // "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

        // 2. 提取 Token 字符串
        parts := strings.SplitN(authHeader, " ", 2)
        token := parts[1]  // 👈 Token 在这里，但没有存储到 Context

        // 3. 解析 Token，获取用户信息
        claims, err := utils.ParseToken(token)

        // 4. 存储用户信息到 Context（但没有存储 Token）
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)

        c.Next()
    }
}
```

---

## 解决方案

### 方案 1：在中间件中存储 Token（推荐）

修改认证中间件，将 Token 也存储到 Context：

```go
// internal/api/middleware/auth.go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, dto.Error(...))
            c.Abort()
            return
        }

        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, dto.Error(...))
            c.Abort()
            return
        }

        token := parts[1]
        if token == "" {
            c.JSON(http.StatusUnauthorized, dto.Error(...))
            c.Abort()
            return
        }

        // 解析Token
        claims, err := utils.ParseToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, dto.Error(...))
            c.Abort()
            return
        }

        // 存储用户信息和 Token
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)
        c.Set("token", token)  // 👈 添加这一行，存储 Token

        c.Next()
    }
}
```

**在 Handler 中使用**：

```go
func (h *UserHandler) SomeHandler(c *gin.Context) {
    // 获取 Token
    token, exists := c.Get("token")
    if !exists {
        c.JSON(http.StatusUnauthorized, dto.Error(...))
        return
    }

    tokenString := token.(string)
    // 使用 tokenString
}
```

---

### 方案 2：创建辅助函数提取 Token

创建一个工具函数，从请求头中提取 Token：

```go
// internal/utils/token.go
package utils

import (
    "strings"
    "github.com/gin-gonic/gin"
)

// GetTokenFromContext 从 Context 中获取 Token
func GetTokenFromContext(c *gin.Context) (string, error) {
    // 先尝试从 Context 获取（如果中间件存储了）
    if token, exists := c.Get("token"); exists {
        return token.(string), nil
    }

    // 否则从请求头提取
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        return "", errors.New("缺少Authorization头")
    }

    parts := strings.SplitN(authHeader, " ", 2)
    if len(parts) != 2 || parts[0] != "Bearer" {
        return "", errors.New("Authorization格式错误")
    }

    return parts[1], nil
}
```

**在 Handler 中使用**：

```go
import "github.com/SwordDragonLee/go-web/internal/utils"

func (h *UserHandler) SomeHandler(c *gin.Context) {
    token, err := utils.GetTokenFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, dto.Error(...))
        return
    }

    // 使用 token
    fmt.Println("当前 Token:", token)
}
```

---

### 方案 3：直接从请求头提取（简单但不推荐）

在 Handler 中直接从请求头提取：

```go
func (h *UserHandler) SomeHandler(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, dto.Error(...))
        return
    }

    parts := strings.SplitN(authHeader, " ", 2)
    if len(parts) != 2 || parts[0] != "Bearer" {
        c.JSON(http.StatusUnauthorized, dto.Error(...))
        return
    }

    token := parts[1]
    // 使用 token
}
```

**缺点**：代码重复，每个 Handler 都要写一遍。

---

## 推荐实现

### 步骤 1：修改中间件存储 Token

```go
// internal/api/middleware/auth.go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // ... 前面的代码 ...

        token := parts[1]

        // 解析Token
        claims, err := utils.ParseToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, dto.Error(...))
            c.Abort()
            return
        }

        // 存储用户信息和 Token
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)
        c.Set("token", token)  // 👈 存储 Token

        c.Next()
    }
}
```

### 步骤 2：创建辅助函数（可选）

```go
// internal/utils/context.go
package utils

import (
    "errors"
    "github.com/gin-gonic/gin"
)

// GetTokenFromContext 从 Context 获取 Token
func GetTokenFromContext(c *gin.Context) (string, error) {
    token, exists := c.Get("token")
    if !exists {
        return "", errors.New("Token不存在")
    }

    tokenString, ok := token.(string)
    if !ok {
        return "", errors.New("Token类型错误")
    }

    return tokenString, nil
}

// GetUserIDFromContext 从 Context 获取用户ID
func GetUserIDFromContext(c *gin.Context) (uint, error) {
    userID, exists := c.Get("user_id")
    if !exists {
        return 0, errors.New("用户ID不存在")
    }

    id, ok := userID.(uint)
    if !ok {
        return 0, errors.New("用户ID类型错误")
    }

    return id, nil
}
```

### 步骤 3：在 Handler 中使用

```go
// internal/api/handler/user_handler.go
import "github.com/SwordDragonLee/go-web/internal/utils"

func (h *UserHandler) SomeHandler(c *gin.Context) {
    // 方式1: 使用辅助函数
    token, err := utils.GetTokenFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, dto.Error(...))
        return
    }

    // 方式2: 直接从 Context 获取
    token, exists := c.Get("token")
    if !exists {
        c.JSON(http.StatusUnauthorized, dto.Error(...))
        return
    }
    tokenString := token.(string)

    // 使用 tokenString
    fmt.Println("当前 Token:", tokenString)
}
```

---

## 使用场景示例

### 场景 1：刷新 Token

```go
func (h *UserHandler) RefreshToken(c *gin.Context) {
    // 获取当前 Token
    oldToken, err := utils.GetTokenFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, dto.Error(...))
        return
    }

    // 验证旧 Token（可选）
    claims, err := utils.ParseToken(oldToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, dto.Error(...))
        return
    }

    // 生成新 Token
    newToken, err := h.userService.RefreshToken(claims.UserID, claims.Username, claims.Role)
    // ...
}
```

### 场景 2：记录操作日志（包含 Token 信息）

```go
func (h *UserHandler) SomeOperation(c *gin.Context) {
    token, _ := utils.GetTokenFromContext(c)

    // 记录操作日志，包含 Token 信息
    logger.Info("用户操作",
        zap.String("token", token[:20]+"..."), // 只记录前20个字符
        zap.Uint("user_id", userID),
    )
}
```

### 场景 3：Token 黑名单（登出功能）

```go
func (h *UserHandler) Logout(c *gin.Context) {
    token, err := utils.GetTokenFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, dto.Error(...))
        return
    }

    // 将 Token 加入黑名单（需要实现黑名单存储）
    blacklist.Add(token)

    c.JSON(http.StatusOK, dto.SuccessWithMessage("登出成功", nil))
}
```

---

## 完整示例代码

### 修改后的中间件

```go
// internal/api/middleware/auth.go
package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/SwordDragonLee/go-web/internal/domain/constant"
    "github.com/SwordDragonLee/go-web/internal/domain/dto"
    "github.com/SwordDragonLee/go-web/internal/utils"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, dto.Error(constant.CodeUnauthorized, "缺少Authorization头"))
            c.Abort()
            return
        }

        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, dto.Error(constant.CodeUnauthorized, "Authorization格式错误"))
            c.Abort()
            return
        }

        token := parts[1]
        if token == "" {
            c.JSON(http.StatusUnauthorized, dto.Error(constant.CodeUnauthorized, "Token为空"))
            c.Abort()
            return
        }

        // 解析Token
        claims, err := utils.ParseToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, dto.Error(constant.CodeUnauthorized, "Token无效或已过期"))
            c.Abort()
            return
        }

        // 存储用户信息和 Token
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)
        c.Set("token", token)  // 👈 存储 Token

        c.Next()
    }
}
```

### 使用示例

```go
// internal/api/handler/user_handler.go
func (h *UserHandler) GetCurrentToken(c *gin.Context) {
    // 获取 Token
    token, exists := c.Get("token")
    if !exists {
        c.JSON(http.StatusUnauthorized, dto.Error(constant.CodeUnauthorized, "Token不存在"))
        return
    }

    tokenString := token.(string)

    // 返回 Token 信息（注意：实际项目中不应该返回完整 Token）
    c.JSON(http.StatusOK, dto.Success(gin.H{
        "token_length": len(tokenString),
        "token_preview": tokenString[:20] + "...", // 只显示前20个字符
    }))
}
```

---

## 安全注意事项

### ⚠️ 不要返回完整 Token

```go
// ❌ 不安全：返回完整 Token
c.JSON(http.StatusOK, dto.Success(gin.H{
    "token": tokenString,  // 不要这样做！
}))44444

// ✅ 安全：只返回必要信息
c.JSON(http.StatusOK, dto.Success(gin.H{
    "token_length": len(tokenString),
    "expires_in": expiresIn,
}))
```

### ⚠️ Token 应该存储在 Context 中

- 不要在日志中记录完整 Token
- 不要在错误消息中返回 Token
- Token 应该只在需要时临时使用

---

## 总结

### 获取 Token 的方式

1. **从 Context 获取**（推荐）

   ```go
   token, _ := c.Get("token")
   ```

2. **使用辅助函数**

   ```go
   token, err := utils.GetTokenFromContext(c)
   ```

3. **从请求头提取**（不推荐，代码重复）

### 推荐步骤

1. ✅ 修改 `AuthMiddleware()`，添加 `c.Set("token", token)`
2. ✅ 创建辅助函数 `GetTokenFromContext()`
3. ✅ 在 Handler 中使用辅助函数获取 Token

这样就能方便地获取当前登录用户的 Token 了！
