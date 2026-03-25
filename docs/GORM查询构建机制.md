# GORM 查询构建机制

## 核心概念

### ❌ `query.Where()` **不执行** SQL 查询！

`query.Where()` 只是**构建查询条件**，不会立即执行 SQL。

### ✅ 只有在调用**执行方法**时才执行 SQL

执行方法包括：
- `Find()` - 查询多条记录
- `First()` - 查询第一条记录
- `Count()` - 统计数量
- `Update()` - 更新记录
- `Delete()` - 删除记录
- `Create()` - 创建记录
- `Save()` - 保存记录

---

## 代码示例解析

### 示例 1：List 方法

```go
// internal/repository/user_repository.go
func (r *userRepository) List(...) ([]*model.User, int64, error) {
    var users []*model.User
    var total int64

    // 1. 创建查询构建器（不执行 SQL）
    query := r.db.Model(&model.User{})
    // 👆 此时还没有执行 SQL，只是创建了一个查询对象

    // 2. 添加 WHERE 条件（不执行 SQL）
    if keyword != "" {
        query = query.Where("username LIKE ?", "%"+keyword+"%")
        // 👆 只是添加条件到查询对象，不执行 SQL
    }

    if status != nil {
        query = query.Where("status = ?", *status)
        // 👆 继续添加条件，仍然不执行 SQL
    }

    // 3. 调用 Count() - 👈 这里才执行 SQL！
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    // SQL: SELECT COUNT(*) FROM users WHERE username LIKE ? AND status = ?

    // 4. 调用 Find() - 👈 这里才执行 SQL！
    if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
        return nil, 0, err
    }
    // SQL: SELECT * FROM users WHERE username LIKE ? AND status = ? LIMIT ? OFFSET ?
}
```

### 执行流程

```
1. query := r.db.Model(&model.User{})
   ↓
   创建查询对象（不执行 SQL）

2. query = query.Where("username LIKE ?", "%test%")
   ↓
   添加 WHERE 条件（不执行 SQL）

3. query = query.Where("status = ?", 1)
   ↓
   继续添加条件（不执行 SQL）

4. query.Count(&total)  ← 👈 这里执行 SQL！
   ↓
   SQL: SELECT COUNT(*) FROM users WHERE username LIKE '%test%' AND status = 1

5. query.Find(&users)  ← 👈 这里执行 SQL！
   ↓
   SQL: SELECT * FROM users WHERE username LIKE '%test%' AND status = 1 LIMIT 10 OFFSET 0
```

---

## GORM 的延迟执行机制

### 链式调用

GORM 使用**链式调用**和**延迟执行**：

```go
// 所有这些都是链式调用，不执行 SQL
query := r.db.Model(&model.User{})
    .Where("username LIKE ?", "%test%")  // 添加条件
    .Where("status = ?", 1)               // 继续添加条件
    .Order("created_at DESC")             // 添加排序
    .Limit(10)                            // 添加限制
    .Offset(0)                            // 添加偏移

// 👆 上面的代码都没有执行 SQL！

// 只有调用执行方法时才执行
err := query.Find(&users).Error  // 👈 这里才执行 SQL！
```

### 查询构建 vs SQL 执行

```go
// ========== 查询构建阶段（不执行 SQL）==========
query := r.db.Model(&model.User{})        // 1. 创建查询对象
query = query.Where("username = ?", "test") // 2. 添加条件
query = query.Where("status = ?", 1)       // 3. 继续添加条件
query = query.Order("created_at DESC")     // 4. 添加排序

// ========== SQL 执行阶段 ==========
err := query.Find(&users).Error  // 👈 执行 SQL: SELECT * FROM users WHERE ...
```

---

## 实际代码示例

### 示例 1：GetByUsername

```go
func (r *userRepository) GetByUsername(username string) (*model.User, error) {
    var user model.User
    
    // 链式调用，构建查询
    err := r.db.Where("username = ?", username).First(&user).Error
    //     │              │                    │
    //     │              │                    └─ 执行方法（执行 SQL）
    //     │              └─ 添加条件（不执行 SQL）
    //     └─ 创建查询对象（不执行 SQL）
    
    return &user, nil
}
```

**执行时机**：
- `r.db` - 不执行 SQL
- `.Where(...)` - 不执行 SQL，只是添加条件
- `.First(&user)` - **执行 SQL**！

**生成的 SQL**：
```sql
SELECT * FROM `users` WHERE username = 'testuser' LIMIT 1
```

### 示例 2：List 方法（复杂查询）

```go
func (r *userRepository) List(...) ([]*model.User, int64, error) {
    var users []*model.User
    var total int64

    // ========== 查询构建（不执行 SQL）==========
    query := r.db.Model(&model.User{})  // 1. 创建查询对象

    if keyword != "" {
        query = query.Where("username LIKE ?", "%"+keyword+"%")  // 2. 添加条件
    }

    if status != nil {
        query = query.Where("status = ?", *status)  // 3. 继续添加条件
    }

    // ========== SQL 执行 ==========
    // 第一次执行 SQL：统计总数
    if err := query.Count(&total).Error; err != nil {  // 👈 执行 SQL
        return nil, 0, err
    }
    // SQL: SELECT COUNT(*) FROM users WHERE username LIKE ? AND status = ?

    // 第二次执行 SQL：查询数据
    offset := (page - 1) * pageSize
    if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {  // 👈 执行 SQL
        return nil, 0, err
    }
    // SQL: SELECT * FROM users WHERE username LIKE ? AND status = ? LIMIT ? OFFSET ?
}
```

---

## 执行方法列表

### 查询方法（执行 SQL）

| 方法 | 说明 | SQL 示例 |
|------|------|----------|
| `Find(&result)` | 查询多条记录 | `SELECT * FROM ...` |
| `First(&result)` | 查询第一条 | `SELECT * FROM ... LIMIT 1` |
| `Last(&result)` | 查询最后一条 | `SELECT * FROM ... ORDER BY id DESC LIMIT 1` |
| `Take(&result)` | 随机取一条 | `SELECT * FROM ... LIMIT 1` |
| `Count(&count)` | 统计数量 | `SELECT COUNT(*) FROM ...` |
| `Pluck("field", &result)` | 查询单个字段 | `SELECT field FROM ...` |

### 修改方法（执行 SQL）

| 方法 | 说明 | SQL 示例 |
|------|------|----------|
| `Create(&model)` | 创建记录 | `INSERT INTO ...` |
| `Update("field", value)` | 更新字段 | `UPDATE ... SET field = ?` |
| `Updates(&model)` | 更新多条字段 | `UPDATE ... SET field1 = ?, field2 = ?` |
| `Delete(&model)` | 删除记录 | `DELETE FROM ...` |
| `Save(&model)` | 保存（创建或更新） | `INSERT INTO ...` 或 `UPDATE ...` |

---

## 查询构建方法（不执行 SQL）

这些方法只是**构建查询**，不执行 SQL：

- `Where()` - 添加 WHERE 条件
- `Or()` - 添加 OR 条件
- `Not()` - 添加 NOT 条件
- `Select()` - 指定查询字段
- `Order()` - 添加排序
- `Group()` - 添加分组
- `Having()` - 添加 HAVING 条件
- `Limit()` - 添加限制
- `Offset()` - 添加偏移
- `Join()` - 添加 JOIN
- `Preload()` - 预加载关联

---

## 完整示例

### 构建复杂查询

```go
func (r *userRepository) ComplexQuery(keyword string, status int, page, pageSize int) ([]*model.User, int64, error) {
    var users []*model.User
    var total int64

    // ========== 查询构建阶段（不执行 SQL）==========
    query := r.db.Model(&model.User{})  // 1. 创建查询对象
    
    // 2. 添加多个条件（都不执行 SQL）
    query = query.Where("username LIKE ?", "%"+keyword+"%")
    query = query.Where("status = ?", status)
    query = query.Where("deleted_at IS NULL")  // 软删除过滤
    query = query.Order("created_at DESC")
    query = query.Select("id, username, email, created_at")  // 只查询指定字段

    // ========== SQL 执行阶段 ==========
    // 执行 SQL 1：统计总数
    if err := query.Count(&total).Error; err != nil {  // 👈 执行 SQL
        return nil, 0, err
    }
    // 生成的 SQL:
    // SELECT COUNT(*) FROM users 
    // WHERE username LIKE '%test%' 
    //   AND status = 1 
    //   AND deleted_at IS NULL

    // 执行 SQL 2：查询数据
    offset := (page - 1) * pageSize
    if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {  // 👈 执行 SQL
        return nil, 0, err
    }
    // 生成的 SQL:
    // SELECT id, username, email, created_at FROM users 
    // WHERE username LIKE '%test%' 
    //   AND status = 1 
    //   AND deleted_at IS NULL 
    // ORDER BY created_at DESC 
    // LIMIT 10 OFFSET 0

    return users, total, nil
}
```

---

## 验证：查看实际执行的 SQL

### 启用 GORM 日志

```go
// pkg/db/db.go
db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),  // 👈 启用日志
})
```

### 日志输出示例

```go
// 代码
query := r.db.Model(&model.User{}).Where("username = ?", "test").First(&user)

// 日志输出（只有 First() 执行时才输出）
[2025-12-18 12:00:00]  [1.23ms]  [rows:1]  
SELECT * FROM `users` WHERE username = 'test' LIMIT 1
```

**注意**：日志只会在**执行方法**调用时输出，`Where()` 调用时不会输出。

---

## 常见错误理解

### ❌ 错误理解 1：Where() 执行 SQL

```go
// 错误理解
query.Where("username = ?", "test")  // 认为这里执行了 SQL

// 正确理解
query.Where("username = ?", "test")  // 只是添加条件，不执行 SQL
query.Find(&users)  // 这里才执行 SQL
```

### ❌ 错误理解 2：每次链式调用都执行 SQL

```go
// 错误理解：认为每次调用都执行一次 SQL
query.Where(...)  // SQL 1
query.Where(...)  // SQL 2
query.Find(...)   // SQL 3

// 正确理解：只有 Find() 执行一次 SQL
query.Where(...)  // 构建条件
query.Where(...)  // 继续构建条件
query.Find(...)   // 执行 SQL（包含所有条件）
```

---

## 总结

### 关键点

1. **`query.Where()` 不执行 SQL**
   - 只是构建查询条件
   - 返回查询对象，可以继续链式调用

2. **只有执行方法才执行 SQL**
   - `Find()`, `First()`, `Count()`, `Update()`, `Delete()` 等
   - 调用执行方法时，GORM 才会生成并执行 SQL

3. **GORM 使用延迟执行**
   - 先构建查询（链式调用）
   - 最后执行（调用执行方法）

4. **可以多次构建，一次执行**
   ```go
   query.Where(...).Where(...).Order(...).Limit(...).Find(&result)
   // 所有构建方法都不执行 SQL，只有 Find() 执行一次 SQL
   ```

### 代码流程

```
查询构建阶段（不执行 SQL）:
query := r.db.Model(&model.User{})
query = query.Where("username = ?", "test")
query = query.Where("status = ?", 1)
query = query.Order("created_at DESC")

SQL 执行阶段:
query.Find(&users)  // 👈 这里才执行 SQL！
```

### 实际 SQL

```sql
SELECT * FROM `users` 
WHERE username = 'test' 
  AND status = 1 
ORDER BY created_at DESC
```

**结论**：`query.Where()` **不执行 SQL**，只是构建查询条件。只有调用执行方法（如 `Find()`, `First()`, `Count()`）时才会执行 SQL！








