# CORS 配置说明

## 命令解释

```bash
go get github.com/gin-contrib/cors
```

### 这个命令的作用

1. **下载依赖包**

   - 从 GitHub 下载 `gin-contrib/cors` 包
   - 这是一个 Gin 框架的 CORS（跨域资源共享）中间件

2. **添加到 go.mod**

   - 自动将依赖添加到 `go.mod` 文件
   - 更新 `go.sum` 文件（依赖校验和）

3. **安装到本地**
   - 将包下载到 `$GOPATH/pkg/mod` 目录
   - 供项目使用

### 为什么需要这个包？

**CORS 问题**：

- 前端运行在 `http://localhost:9000`
- 后端运行在 `http://localhost:8080`
- 浏览器阻止跨域请求（不同端口视为不同源）

**解决方案**：

- 使用 `gin-contrib/cors` 中间件
- 在响应头中添加 `Access-Control-Allow-Origin`
- 允许前端跨域访问

---

## 完整解决方案

### 步骤 1：安装 CORS 包

```bash
go get github.com/gin-contrib/cors
```

### 步骤 2：在路由中配置 CORS

修改 `internal/api/router/router.go`：

```go
package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/SwordDragonLee/go-web/internal/api/handler"
	"github.com/SwordDragonLee/go-web/internal/api/middleware"
)

// SetupRouter 设置路由
func SetupRouter(userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()

	// 配置 CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:9000"}, // 前端地址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // 允许携带凭证（Cookie、Authorization 等）
		MaxAge:           12 * time.Hour,
	}))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// ... 其他路由
}
```

### 步骤 3：添加 time 包导入

```go
import (
	"time"
	// ... 其他导入
)
```

---

## CORS 配置选项说明

### AllowOrigins（允许的源）

```go
AllowOrigins: []string{"http://localhost:9000"}
```

- 指定允许跨域的前端地址
- 可以添加多个地址
- 生产环境建议使用具体域名

**示例**：

```go
// 开发环境
AllowOrigins: []string{"http://localhost:9000", "http://localhost:3000"}

// 生产环境
AllowOrigins: []string{"https://yourdomain.com"}
```

### AllowMethods（允许的 HTTP 方法）

```go
AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
```

- 指定允许的 HTTP 请求方法
- `OPTIONS` 是预检请求必需的

### AllowHeaders（允许的请求头）

```go
AllowHeaders: []string{"Origin", "Content-Type", "Authorization"}
```

- 指定允许的请求头
- `Authorization` 用于携带 Token
- `Content-Type` 用于 JSON 请求

### AllowCredentials（允许凭证）

```go
AllowCredentials: true
```

- 允许携带 Cookie、Authorization 等凭证
- 如果设置为 `true`，`AllowOrigins` 不能使用 `"*"`

### MaxAge（预检请求缓存时间）

```go
MaxAge: 12 * time.Hour
```

- 预检请求（OPTIONS）的缓存时间
- 减少重复的预检请求

---

## 简化配置（开发环境）

如果只是开发环境，可以使用简化配置：

```go
import "github.com/gin-contrib/cors"

func SetupRouter(userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()

	// 简化配置（仅开发环境）
	r.Use(cors.Default()) // 允许所有源，仅用于开发！

	// ... 其他路由
}
```

**注意**：`cors.Default()` 允许所有源，**仅用于开发环境**，生产环境必须配置具体域名！

---

## 完整示例代码

```go
package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/SwordDragonLee/go-web/internal/api/handler"
	"github.com/SwordDragonLee/go-web/internal/api/middleware"
)

// SetupRouter 设置路由
func SetupRouter(userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()

	// 配置 CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:9000"}, // 前端地址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1路由组
	v1 := r.Group("/api/v1")
	{
		// 认证相关路由（无需认证）
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
			auth.POST("/send-code", userHandler.SendVerificationCode)
			auth.POST("/verify-email", userHandler.VerifyEmail)
			auth.POST("/reset-password", userHandler.ResetPassword)
		}

		// 用户相关路由（需要认证）
		user := v1.Group("/user")
		user.Use(middleware.AuthMiddleware())
		{
			user.GET("/profile", userHandler.GetProfile)
			user.PUT("/profile", userHandler.UpdateProfile)
			user.PUT("/password", userHandler.ChangePassword)
			user.GET("/login-logs", userHandler.GetLoginLogs)
			user.POST("/refresh-token", userHandler.RefreshToken)
		}

		// 管理员路由（需要管理员权限）
		admin := v1.Group("/users")
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.AdminMiddleware())
		{
			admin.GET("", userHandler.GetUserList)
			admin.PUT("/:id/status", userHandler.UpdateUserStatus)
			admin.DELETE("/:id", userHandler.DeleteUser)
		}
	}

	return r
}
```

---

## 验证 CORS 配置

配置完成后，重新启动服务器，前端请求应该可以正常发送。

**检查响应头**：

```bash
curl -X OPTIONS http://localhost:8080/api/v1/auth/login \
  -H "Origin: http://localhost:9000" \
  -H "Access-Control-Request-Method: POST" \
  -v
```

应该看到响应头包含：

```
Access-Control-Allow-Origin: http://localhost:9000
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Origin, Content-Type, Authorization
```

---

## 常见问题

### Q1: 为什么需要 OPTIONS 方法？

**A**: 浏览器在发送跨域请求前，会先发送一个 OPTIONS 预检请求，检查服务器是否允许跨域。

### Q2: 生产环境如何配置？

**A**:

```go
AllowOrigins: []string{"https://yourdomain.com"} // 使用具体域名
// 不要使用 "*" 或 cors.Default()
```

### Q3: 多个前端地址怎么办？

**A**:

```go
AllowOrigins: []string{
	"http://localhost:9000",
	"http://localhost:3000",
	"https://yourdomain.com",
}
```

### Q4: 为什么设置了 CORS 还是报错？

**A**: 检查：

1. 服务器是否重启
2. `AllowOrigins` 是否包含前端地址
3. `AllowHeaders` 是否包含 `Authorization`
4. `AllowCredentials` 是否正确设置

---

## 总结

1. **安装包**：`go get github.com/gin-contrib/cors`
2. **配置中间件**：在路由中添加 CORS 配置
3. **设置允许的源**：`AllowOrigins: []string{"http://localhost:9000"}`
4. **允许凭证**：`AllowCredentials: true`（用于 Token）
5. **重启服务器**：使配置生效

配置完成后，前端就可以正常跨域访问后端 API 了！
