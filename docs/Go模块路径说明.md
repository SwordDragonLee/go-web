# Go 模块路径说明

## 为什么使用 `github.com/SwordDragonLee/go-web`？

### 核心原因：Go 模块系统

在 Go 1.11+ 中，引入了**模块（Module）**系统。每个项目都需要定义一个**模块名**，这个模块名就是项目的"根路径"。

### 模块名定义

查看 `go.mod` 文件的第一行：

```go
module github.com/SwordDragonLee/go-web
```

这行代码定义了：

- **模块名**：`github.com/SwordDragonLee/go-web`
- 这是整个项目的**唯一标识符**
- 所有内部包的导入路径都**必须**基于这个模块名

---

## 导入路径规则

### 内部包导入格式

```
模块名 + 包路径（相对于项目根目录）
```

**示例**：

```
模块名: github.com/SwordDragonLee/go-web
包路径: internal/service

完整导入路径: github.com/SwordDragonLee/go-web/internal/service
```

### 实际代码示例

```go
// internal/service/user_service.go
package service

import (
    // ✅ 正确：使用模块名 + 相对路径
    "github.com/SwordDragonLee/go-web/internal/domain/dto"
    "github.com/SwordDragonLee/go-web/internal/repository"

    // ❌ 错误：不能使用相对路径
    // "../domain/dto"  // 编译错误！
    // "./domain/dto"   // 编译错误！
)
```

---

## 为什么使用 GitHub 路径格式？

### 1. 约定俗成的标准

虽然代码在本地，但使用 GitHub 路径格式是 Go 社区的**标准做法**：

- ✅ 如果将来要发布到 GitHub，路径已经匹配
- ✅ 其他项目可以轻松导入你的包
- ✅ 符合 Go 官方推荐的最佳实践

### 2. 唯一性保证

GitHub 路径格式确保模块名的**全局唯一性**：

```
github.com/用户名/项目名
```

这样可以避免不同开发者的模块名冲突。

### 3. 便于依赖管理

如果项目被其他项目引用，Go 工具链可以通过这个路径：

- 自动下载依赖
- 管理版本
- 解析依赖关系

---

## 可以修改模块名吗？

### 可以！但需要注意

你可以将模块名改为任何你想要的名称：

```go
// go.mod
module go-web                    // 简单名称
// 或
module myproject                 // 自定义名称
// 或
module company.com/go-web        // 公司域名
```

**修改步骤**：

1. **修改 go.mod**

   ```go
   module go-web  // 改为新名称
   ```

2. **修改所有导入语句**

   ```go
   // 将所有文件中的导入路径改为：
   import "go-web/internal/service"
   import "go-web/internal/repository"
   ```

3. **运行命令更新**
   ```bash
   go mod tidy
   ```

---

## 不同模块名示例

### 示例 1：简单名称

```go
// go.mod
module go-web

// 导入
import "go-web/internal/service"
```

**优点**：

- 简洁
- 输入方便

**缺点**：

- 如果发布到 GitHub，需要修改所有导入路径
- 可能与其他项目冲突

### 示例 2：GitHub 路径（当前）

```go
// go.mod
module github.com/SwordDragonLee/go-web

// 导入
import "github.com/SwordDragonLee/go-web/internal/service"
```

**优点**：

- 全局唯一
- 便于发布和分享
- 符合 Go 社区标准

**缺点**：

- 路径较长

### 示例 3：公司域名

```go
// go.mod
module company.com/go-web

// 导入
import "company.com/go-web/internal/service"
```

**优点**：

- 专业
- 适合企业内部项目

---

## 实际工作原理

### Go 工具链如何解析导入路径？

```
导入路径: github.com/SwordDragonLee/go-web/internal/service
         │                              │
         └── 模块名                      └── 包路径

Go 工具链会：
1. 检查是否是标准库（如 "fmt", "net/http"）
   → 不是，继续

2. 检查是否是第三方依赖（go.mod 的 require）
   → 不是，继续

3. 检查是否是当前模块的内部包
   → 是！模块名匹配 github.com/SwordDragonLee/go-web
   → 查找项目根目录下的 internal/service 目录
   → 找到！使用本地代码
```

### 本地 vs 远程

**本地代码**（当前项目）：

```go
import "github.com/SwordDragonLee/go-web/internal/service"
// Go 会在本地项目目录查找
```

**远程依赖**（其他项目）：

```go
import "github.com/gin-gonic/gin"
// Go 会从 GitHub 下载（如果不在本地缓存）
```

---

## 常见问题

### Q1: 为什么不能使用相对路径？

**A**: Go 模块系统**不支持**相对路径导入，这是设计决策：

```go
// ❌ 不支持
import "../domain/dto"
import "./utils"

// ✅ 必须使用完整路径
import "github.com/SwordDragonLee/go-web/internal/domain/dto"
```

**原因**：

- 确保导入路径的唯一性和明确性
- 便于代码重构和移动
- 支持跨模块引用

### Q2: 模块名必须和 GitHub 仓库名一致吗？

**A**: **不一定**！模块名可以是任何字符串。但如果：

- 要发布到 GitHub：建议保持一致
- 只是本地项目：可以随意命名

### Q3: 修改模块名后需要做什么？

**A**:

1. 修改 `go.mod` 中的模块名
2. 修改所有 `.go` 文件中的导入路径
3. 运行 `go mod tidy` 更新依赖

### Q4: 可以使用中文或特殊字符吗？

**A**: 可以，但不推荐。建议使用：

- 小写字母
- 数字
- 连字符（-）
- 点号（.）用于域名

---

## 总结

### 关键点

1. **模块名是必需的**

   - 每个 Go 项目都必须定义模块名
   - 在 `go.mod` 的第一行

2. **导入路径基于模块名**

   - 所有内部包导入 = 模块名 + 包路径
   - 不能使用相对路径

3. **GitHub 路径是标准做法**

   - 即使代码在本地，也推荐使用
   - 便于将来发布和分享

4. **可以自定义模块名**
   - 根据项目需求选择
   - 修改后需要更新所有导入路径

### 当前项目的配置

```go
// go.mod
module github.com/SwordDragonLee/go-web

// 所有导入都基于这个模块名
import "github.com/SwordDragonLee/go-web/internal/service"
```

这是**标准且推荐**的配置方式！✅
