# GORM 如何确定查询哪张表

## 核心机制

GORM 通过**结构体类型**来确定表名。当你传入一个结构体变量时，GORM 会：

1. **识别结构体类型**：通过反射获取结构体的类型信息
2. **查找表名**：根据规则确定对应的数据库表名
3. **执行查询**：使用确定的表名构建 SQL 语句

---

## 表名确定规则

### 规则 1：TableName() 方法（优先级最高）

如果结构体实现了 `TableName()` 方法，GORM 会使用该方法返回的表名。

**示例**：

```go
// internal/domain/model/user.go
type User struct {
    ID       uint
    Username string
    // ...
}

// TableName 指定表名
func (User) TableName() string {
    return "users"  // 👈 GORM 会使用这个表名
}
```

**查询时**：

```go
var user model.User
r.db.Where("username = ?", username).First(&user)
// GORM 会查询 "users" 表
// SQL: SELECT * FROM users WHERE username = ? LIMIT 1
```

### 规则 2：自动复数化（默认规则）

如果没有定义 `TableName()` 方法，GORM 会将结构体名转换为复数形式：

```go
type User struct { ... }        // → users 表
type Product struct { ... }     // → products 表
type Category struct { ... }   // → categories 表
type Person struct { ... }      // → people 表（特殊规则）
```

---

## 代码示例解析

### 示例 1：GetByUsername 方法

```go
// internal/repository/user_repository.go
func (r *userRepository) GetByUsername(username string) (*model.User, error) {
    var user model.User  // 👈 声明 User 类型的变量
    
    // GORM 通过 &user 的类型（*model.User）确定：
    // 1. 结构体类型是 model.User
    // 2. 查找 TableName() 方法 → 返回 "users"
    // 3. 执行查询：SELECT * FROM users WHERE username = ? LIMIT 1
    err := r.db.Where("username = ?", username).First(&user).Error
    
    return &user, nil
}
```

**执行流程**：

```
1. var user model.User
   ↓
2. r.db.Where(...).First(&user)
   ↓
3. GORM 反射获取 user 的类型：model.User
   ↓
4. 查找 TableName() 方法：func (User) TableName() string { return "users" }
   ↓
5. 使用表名 "users"
   ↓
6. 生成 SQL: SELECT * FROM users WHERE username = ? LIMIT 1
```

### 示例 2：显式指定表名

```go
// 方式1: 使用 Model
r.db.Model(&model.User{}).Where("username = ?", username).First(&user)

// 方式2: 使用 Table
r.db.Table("users").Where("username = ?", username).First(&user)

// 方式3: 使用 TableName() 方法（当前项目使用的方式）
var user model.User
r.db.Where("username = ?", username).First(&user)
// 自动使用 User.TableName() 返回的 "users"
```

---

## 实际 SQL 生成

### 查看生成的 SQL

可以启用 GORM 的日志来查看实际生成的 SQL：

```go
// pkg/db/db.go
db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info), // 👈 会输出 SQL
})
```

**输出示例**：

```
[2025-12-18 12:00:00]  [1.23ms]  [rows:1]  SELECT * FROM `users` WHERE username = 'testuser' LIMIT 1
```

---

## 表名确定优先级

```
1. TableName() 方法（最高优先级）
   ↓ 如果没有
2. 显式指定（Table() 或 Model()）
   ↓ 如果没有
3. 自动复数化（User → users）
```

---

## 完整示例

### User 模型定义

```go
// internal/domain/model/user.go
package model

type User struct {
    ID       uint   `gorm:"primarykey"`
    Username string `gorm:"uniqueIndex"`
    Email    string
    // ...
}

// TableName 指定表名
func (User) TableName() string {
    return "users"  // 👈 明确指定表名为 "users"
}
```

### Repository 查询

```go
// internal/repository/user_repository.go
func (r *userRepository) GetByUsername(username string) (*model.User, error) {
    var user model.User  // 👈 User 类型
    
    // GORM 执行流程：
    // 1. 识别类型：model.User
    // 2. 查找 TableName()：返回 "users"
    // 3. 生成 SQL：SELECT * FROM users WHERE username = ? LIMIT 1
    err := r.db.Where("username = ?", username).First(&user).Error
    
    return &user, nil
}
```

### 生成的 SQL

```sql
SELECT * FROM `users` WHERE username = 'testuser' LIMIT 1
```

---

## 其他表名指定方式

### 方式 1：使用 Table() 方法

```go
r.db.Table("custom_table_name").Where("username = ?", username).First(&user)
```

### 方式 2：使用 Model() 方法

```go
r.db.Model(&model.User{}).Where("username = ?", username).First(&user)
// 会使用 User.TableName() 返回的表名
```

### 方式 3：全局表名规则

```go
// 在初始化时设置
db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
    NamingStrategy: schema.NamingStrategy{
        TablePrefix:   "t_",      // 表前缀
        SingularTable: true,      // 使用单数表名（User → user，而不是 users）
    },
})
```

---

## 常见问题

### Q1: 如果结构体名和表名不一致怎么办？

**A**: 实现 `TableName()` 方法：

```go
type User struct { ... }

func (User) TableName() string {
    return "t_user"  // 自定义表名
}
```

### Q2: 如何查看 GORM 生成的 SQL？

**A**: 启用日志：

```go
db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})
```

### Q3: 可以在运行时动态指定表名吗？

**A**: 可以，使用 `Table()` 方法：

```go
tableName := "users_" + time.Now().Format("200601")
r.db.Table(tableName).Where(...).First(&user)
```

### Q4: 为什么 First() 需要传入指针？

**A**: 
- GORM 需要通过指针修改结构体的值
- 通过指针类型确定结构体类型，从而确定表名

---

## 总结

### 关键点

1. **GORM 通过结构体类型确定表名**
   - `var user model.User` → GORM 识别 `User` 类型
   - `First(&user)` → 通过指针类型确定结构体类型

2. **表名确定优先级**
   - `TableName()` 方法（最高）
   - 显式指定（`Table()` 或 `Model()`）
   - 自动复数化（默认）

3. **当前项目的实现**
   ```go
   func (User) TableName() string {
       return "users"  // 明确指定表名
   }
   ```

4. **实际 SQL 生成**
   ```sql
   SELECT * FROM `users` WHERE username = ? LIMIT 1
   ```

### 代码流程

```
var user model.User
    ↓
r.db.Where(...).First(&user)
    ↓
GORM 反射获取类型：model.User
    ↓
查找 TableName() 方法
    ↓
返回表名："users"
    ↓
生成 SQL：SELECT * FROM users WHERE ...
    ↓
执行查询
```

现在你明白了：GORM 通过 `&user` 的类型（`*model.User`）确定结构体是 `User`，然后查找 `TableName()` 方法获取表名 "users"！










