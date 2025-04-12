# 视频管理平台后端服务

## 项目简介
该项目是一个基于Go语言和Gin框架开发的视频管理平台后端服务，提供视频文件的基础管理功能，包括上传、获取和管理等功能。项目采用清晰的分层架构，遵循RESTful API设计规范，并实现了统一的错误处理和响应格式。

### 快速开始
1. 环境要求
   - Go 1.20+
   - MongoDB 4.4+
   - Redis 6.0+（可选，用于缓存）
   - FFmpeg（用于视频处理）

2. 安装依赖
```bash
# 克隆项目
git clone [项目地址]
cd video-platform

# 安装Go依赖
go mod download

# 安装FFmpeg（MacOS）
brew install ffmpeg

# 安装FFmpeg（Ubuntu/Debian）
sudo apt-get update
sudo apt-get install ffmpeg
```

3. 配置文件
```bash
# 复制配置文件模板
cp config/config.example.yaml config/config.yaml

# 修改配置文件
vim config/config.yaml
```

4. 运行项目
```bash
# 开发模式运行
go run cmd/main.go

# 或者编译后运行
go build -o video-platform cmd/main.go
./video-platform
```

### 开发指南
1. 项目结构说明
   - `api/`: API接口定义和文档
   - `cmd/`: 项目入口文件
   - `config/`: 配置文件和配置结构定义
   - `internal/`: 内部包
     - `handler/`: HTTP请求处理器
     - `middleware/`: 中间件
     - `model/`: 数据模型定义
     - `service/`: 业务逻辑层
   - `pkg/`: 可重用的工具包
   - `scripts/`: 工具脚本
   - `test/`: 测试文件

2. 开发规范
   - 代码风格遵循Go标准
   - 使用gofmt格式化代码
   - 必须添加单元测试
   - 遵循依赖注入原则
   - 使用接口进行解耦

3. 错误处理
   - 统一使用pkg/response包处理响应
   - 错误信息必须明确且用户友好
   - 日志记录必须包含足够上下文

4. 测试
   - 运行所有测试：`go test ./...`
   - 运行指定包测试：`go test ./internal/service`
   - 生成测试覆盖率报告：`go test -cover ./...`

## 主要功能
### 用户管理
- ✅ 用户注册与登录
  - ✅ 支持用户名、密码和邮箱注册
  - ✅ 密码加密存储
  - ✅ JWT令牌认证
- ✅ 用户认证
  - ✅ 基于JWT的身份验证
  - ✅ 接口权限控制
  - ✅ 用户会话管理

### 视频上传
- ✅ 视频上传（支持断点续传）
- ✅ 支持大文件上传（最大1GB）
- ✅ 支持主流视频格式（mp4、mov、avi等）
- ✅ 上传时可设置视频标题和描述信息

### 视频获取
- ✅ 视频流式播放
- ✅ 视频列表查询（支持分页、排序、搜索）
- ✅ 视频详情获取

### 视频管理
- ✅ 视频管理（列表、详情、修改、删除）
- ✅ 支持批量删除
- ✅ 支持修改视频基本信息

### 用户标记、注释和笔记功能
- ✅ 用户可以在视频播放过程中添加标记
  - ✅ 支持添加、更新、删除标记
  - ✅ 可以为标记添加时间戳和内容
- ✅ 用户可以为每个标记添加文字注释
  - ✅ 支持添加、更新、删除注释
  - ✅ 注释与标记关联
- ✅ 用户可以在视频播放过程中随时记录笔记
  - ✅ 支持添加、更新、删除笔记
  - ✅ 可以为笔记添加时间戳和内容
- ✅ 支持导出所有标记、注释和笔记
  - ✅ 导出为文本格式
  - ✅ 包含时间戳和内容信息

## 技术栈
- Go 1.20+
- Gin Web框架
- MongoDB数据库
- 文件存储系统

## 项目特点
- 统一的响应格式
- 完善的错误处理
- 规范的代码注释
- 合理的项目分层

## 开发日志

### 2024-03-16 标记和笔记功能实现
1. 标记功能实现
   - 添加标记模型和服务
   - 实现标记的增删改查API
   - 支持按视频ID获取标记列表
   - 添加标记权限验证

2. 注释功能实现
   - 添加注释模型
   - 实现注释与标记的关联
   - 支持注释的增删改查API
   - 添加注释权限验证

3. 笔记功能实现
   - 添加笔记模型和服务
   - 实现笔记的增删改查API
   - 支持按视频ID获取笔记列表
   - 添加笔记权限验证

4. 导出功能
   - 支持导出视频相关的所有标记、注释和笔记
   - 支持多种导出格式
   - 添加导出权限验证

5. 用户认证完善
   - 完善用户认证机制
   - 添加路由权限控制
   - 实现用户注册和登录API

### 2024-02-21 功能完善
1. 视频时长功能
   - 添加视频时长字段（duration）
   - 上传时必须提供视频时长
   - 修改相关测试用例
   - 更新 API 文档

2. 测试改进
   - 完善测试数据清理机制
   - 修复测试用例中的数据验证
   - 统一响应格式验证
   - 添加更多边界测试

3. 文档更新
   - 更新 API 文档，添加时长字段说明
   - 完善响应示例
   - 添加字段验证说明
   - 更新开发日志

4. 经验总结
   - 上传接口需要严格的参数验证
   - 测试数据要及时清理
   - 保持文档的同步更新
   - 统一的响应格式很重要

### 2024-04-12 用户注册功能优化
1. 问题发现与分析
   - 发现用户注册时的表单验证问题
   - 错误日志：`Field validation for 'Username' failed on the 'min' tag`
   - 原因：用户名长度不满足最小要求（3个字符）
   - 当前错误提示不够友好，只返回"无效的请求参数"

2. 验证规则完善
   - 用户名验证规则
     - 必填字段
     - 长度限制：3-32个字符
     - 允许字符：字母、数字、下划线
     - 不允许纯数字
   - 密码验证规则
     - 必填字段
     - 长度限制：6-32个字符
     - 必须包含：字母和数字
     - 建议包含特殊字符
   - 邮箱验证规则
     - 必填字段
     - 必须是有效的邮箱格式
     - 不允许临时邮箱域名

3. 代码优化
   - 完善验证规则定义
   - 添加自定义验证错误信息
   - 实现验证错误的详细返回

4. 下一步计划
   - 实现前端表单验证
   - 添加密码强度检查
   - 优化错误提示信息
   - 添加用户名保留字列表

5. 开发经验总结
   - 表单验证应该前后端同时进行
   - 错误提示要清晰明确
   - 安全性和用户体验要平衡
   - 日志记录要包含关键信息

## 项目结构
```
.
├── api/            # API接口定义
├── cmd/            # 主程序入口
│   └── main.go     # 主程序入口文件
├── config/         # 配置文件
│   ├── config.go   # 配置结构定义
│   └── config_test.go # 配置测试文件
├── internal/       # 内部包
│   ├── handler/    # HTTP处理器
│   │   ├── middleware.go # 中间件
│   │   ├── routes.go    # 路由定义
│   │   ├── user.go      # 用户处理器
│   │   ├── video.go     # 视频处理器
│   │   └── mark.go      # 标记处理器（包含标记、注释和笔记）
│   ├── middleware/ # 中间件
│   │   └── auth.go      # 认证中间件
│   ├── model/     # 数据模型
│   │   ├── user.go      # 用户模型
│   │   ├── video.go     # 视频模型
│   │   └── mark.go      # 标记模型（包含标记、注释和笔记）
│   └── service/   # 业务逻辑
│       ├── user.go      # 用户服务
│       ├── video.go     # 视频服务
│       └── mark.go      # 标记服务（包含标记、注释和笔记）
├── pkg/           # 可重用的包
│   ├── database/  # 数据库相关
│   │   ├── mongodb.go   # MongoDB连接
│   │   └── mongodb_test.go # MongoDB测试
│   └── response/  # 响应处理
│       └── response.go  # 统一响应格式
├── scripts/       # 脚本文件
│   └── clean_db.go      # 数据库清理工具
└── test/          # 测试文件
```

## API文档
### 统一响应格式
```json
{
    "code": 0,      // 0:成功，1:失败
    "msg": "success", // 响应信息
    "data": {}      // 响应数据
}
```

### 视频管理接口

#### 1. 上传视频
- 请求方法: `POST`
- 路径: `/api/v1/videos`
- Content-Type: `multipart/form-data`
- 参数:
  - `file`: 视频文件（必填，支持 mp4/mov/avi/wmv/flv/mkv）
  - `cover`: 封面图文件（可选，支持 jpg/jpeg/png，最大2MB）
  - `title`: 视频标题（必填）
  - `duration`: 视频时长（必填，单位：秒，支持小数点后1位）
  - `description`: 视频描述（可选）
  - `status`: 视频状态（可选，默认为 private）
  - `tags`: 标签（可选，多个标签用逗号分隔）

#### 响应
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "id": "视频ID",
        "title": "视频标题",
        "description": "视频描述",
        "fileName": "存储的文件名",
        "fileSize": 1024,
        "format": "mp4",
        "status": "private",
        "duration": 180.5,
        "coverUrl": "封面图URL",
        "thumbnailUrl": "缩略图URL",
        "createdAt": "2024-01-20T10:00:00Z",
        "updatedAt": "2024-01-20T10:00:00Z"
    }
}
```

### 2. 获取视频列表
- 请求方法: `GET`
- 路径: `/api/v1/videos`
- 参数:
  - `page`: 页码（默认1）
  - `pageSize`: 每页数量（默认10，最大50）
  - `keyword`: 关键词搜索，匹配标题和描述
  - `status`: 视频状态筛选
  - `startDate`: 开始日期（格式：YYYY-MM-DD）
  - `endDate`: 结束日期（格式：YYYY-MM-DD）
  - `tags`: 标签筛选，多个标签用逗号分隔
  - `sortBy`: 排序字段（created_at/views/likes/file_size）
  - `sortOrder`: 排序方向（asc/desc）

#### 响应
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "total": 100,
        "items": [{
            "id": "视频ID",
            "title": "视频标题",
            "description": "视频描述",
            "fileSize": 1024,
            "format": "mp4",
            "status": "public",
            "thumbnailUrl": "缩略图URL",
            "tags": ["标签1", "标签2"],
            "stats": {
                "views": 100,
                "likes": 50,
                "comments": 20,
                "shares": 10
            },
            "createdAt": "2024-01-20T10:00:00Z"
        }]
    }
}
```

### 3. 获取视频详情
- 请求方法: `GET`
- 路径: `/api/v1/videos/:id`

#### 响应
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "id": "视频ID",
        "title": "视频标题",
        "description": "视频描述",
        "fileSize": 1024,
        "format": "mp4",
        "status": "public",
        "tags": ["标签1", "标签2"],
        "thumbnailUrl": "缩略图URL",
        "stats": {
            "views": 100,
            "likes": 50,
            "comments": 20,
            "shares": 10
        },
        "createdAt": "2024-01-20T10:00:00Z",
        "updatedAt": "2024-01-20T10:00:00Z"
    }
}
```

### 4. 更新视频信息
- 请求方法: `PUT`
- 路径: `/api/v1/videos/:id`
- Content-Type: `application/json`
- 请求体:
```json
{
    "title": "新标题",
    "description": "新描述",
    "status": "public",
    "tags": ["标签1", "标签2"]
}
```

#### 响应
```json
{
    "code": 0,
    "msg": "success",
    "data": null
}
```

### 5. 删除视频
- 请求方法: `DELETE`
- 路径: `/api/v1/videos/:id`

#### 响应
```json
{
    "code": 0,
    "msg": "success",
    "data": null
}
```

### 6. 批量操作视频
- 请求方法: `POST`
- 路径: `/api/v1/videos/batch`
- Content-Type: `application/json`
- 请求体:
```json
{
    "ids": ["视频ID1", "视频ID2"],
    "action": "update_status",  // delete/update_status
    "status": "private"         // 当action为update_status时需要，可选值：public/private/draft
}
```

#### 响应
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "successCount": 2,
        "failedCount": 0,
        "failedIds": []
    }
}
```

### 7. 更新视频缩略图
- 请求方法: `POST`
- 路径: `/api/v1/videos/:id/thumbnail`
- Content-Type: `multipart/form-data`
- 参数:
  - `file`: 图片文件（必填，支持jpg/png/gif，最大2MB）

#### 响应
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "thumbnailUrl": "缩略图URL"
    }
}
```

### 8. 获取视频统计信息
- 请求方法: `GET`
- 路径: `/api/v1/videos/:id/stats`

#### 响应
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "views": 100,
        "likes": 50,
        "comments": 20,
        "shares": 10
    }
}
```

### 9. 视频流式播放
- 请求方法: `GET`
- 路径: `/api/v1/videos/:id/stream`
- 支持范围请求（Range header）

#### 响应
- Content-Type: video/*
- Accept-Ranges: bytes
- 支持断点续传
- 直接返回视频流

### 用户管理接口

#### 1. 用户注册
- 请求方法: `POST`
- 路径: `/api/v1/users/register`
- Content-Type: `application/json`
- 请求体:
```json
{
    "username": "user123",
    "password": "password123",
    "email": "user@example.com"
}
```

#### 响应
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "id": "用户ID",
        "username": "user123",
        "email": "user@example.com",
        "status": 1,
        "createdAt": "2024-03-16T12:00:00Z",
        "updatedAt": "2024-03-16T12:00:00Z"
    }
}
```

#### 2. 用户登录
- 请求方法: `POST`
- 路径: `/api/v1/users/login`
- Content-Type: `application/json`
- 请求体:
```json
{
    "username": "user123",
    "password": "password123"
}
```

#### 响应
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "user": {
            "id": "用户ID",
            "username": "user123",
            "email": "user@example.com",
            "status": 1,
            "createdAt": "2024-03-16T12:00:00Z",
            "updatedAt": "2024-03-16T12:00:00Z"
        },
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
}
```

### 标记、注释和笔记管理接口

#### 10. 添加标记
- 请求方法: `POST`
- 路径: `/api/v1/marks`
- Content-Type: `application/json`
- 请求体:
```json
{
    "videoId": "视频ID",
    "timestamp": 125.5,
    "content": "这是一个标记"
}
```

#### 响应
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "id": "标记ID",
        "userId": "用户ID",
        "videoId": "视频ID",
        "timestamp": 125.5,
        "content": "这是一个标记",
        "annotations": [],
        "createdAt": "2024-03-16T10:00:00Z",
        "updatedAt": "2024-03-16T10:00:00Z"
    }
}
```

#### 11. 获取标记列表
- 请求方法: `GET`
```

## 常见问题
1. 运行时遇到"MongoDB连接失败"
   ```
   解决方案：
   1. 检查MongoDB服务是否启动
   2. 验证配置文件中的连接字符串
   3. 确认网络连接是否正常
   ```

2. 上传视频失败
   ```
   解决方案：
   1. 检查文件大小是否超过限制（默认1GB）
   2. 确认文件格式是否支持
   3. 验证存储目录权限
   ```

3. JWT token验证失败
   ```
   解决方案：
   1. 检查token是否过期
   2. 验证token格式是否正确
   3. 确认密钥配置是否正确
   ```

## 贡献指南
1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交Pull Request

## 版本历史
- v0.1.0 (2024-02-21)
  - 基础功能实现
  - 视频上传和管理
  - 用户认证系统
- v0.2.0 (2024-03-16)
  - 添加标记和笔记功能
  - 完善权限控制
  - 优化错误处理
- v0.2.1 (2024-04-12)
  - 优化用户注册流程
  - 完善表单验证
  - 改进错误提示

## 维护者
- 小江 (@xiaojiang)

## 开源协议
本项目采用 MIT 协议 - 详见 [LICENSE](LICENSE) 文件

## 技术实现细节

### 1. 架构设计
```
+------------------+
|     客户端        |
+--------+---------+
         |
         | HTTP/HTTPS
         |
+--------+---------+
|   Nginx 反向代理   |
+--------+---------+
         |
         |
+--------+---------+
|    Gin Web框架    |
+--------+---------+
         |
+---+----+----+----+
    |         |
+---+---+ +---+----+
|Handler| |中间件   |
+---+---+ +---+----+
    |         |
+---+---+ +---+----+
|Service| |工具包   |
+---+---+ +---+----+
    |         |
+---+---+ +---+----+
|Model  | |数据库   |
+-------+ +--------+
```

### 2. 核心模块说明

#### 2.1 用户认证模块
- JWT认证实现
  ```go
  // JWT密钥配置
  var jwtKey = []byte(config.GlobalConfig.JWT.Secret)
  
  // JWT Claims结构
  type Claims struct {
      UserID   string
      Username string
      jwt.StandardClaims
  }
  
  // Token生成逻辑
  func GenerateToken(userID, username string) (string, error) {
      expirationTime := time.Now().Add(24 * time.Hour)
      claims := &Claims{
          UserID:   userID,
          Username: username,
          StandardClaims: jwt.StandardClaims{
              ExpiresAt: expirationTime.Unix(),
          },
      }
      token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
      return token.SignedString(jwtKey)
  }
  ```

#### 2.2 视频处理模块
- 文件上传
  ```go
  // 支持的视频格式
  var validVideoFormats = map[string]bool{
      ".mp4": true,
      ".mov": true,
      ".avi": true,
      ".wmv": true,
      ".flv": true,
      ".mkv": true,
  }
  
  // 分片上传配置
  const (
      maxFileSize   = 1 << 30    // 1GB
      chunkSize     = 1 << 22    // 4MB
      maxRetries    = 3
  )
  ```

- 视频流式播放
  ```go
  // 支持范围请求的视频流处理
  func (h *VideoHandler) Stream(c *gin.Context) {
      // 支持Range header
      rangeHeader := c.GetHeader("Range")
      // 实现206 Partial Content响应
      // 支持断点续传
  }
  ```

#### 2.3 数据库设计
- MongoDB集合设计
  ```go
  // 用户集合
  type User struct {
      ID        primitive.ObjectID `bson:"_id,omitempty"`
      Username  string            `bson:"username"`
      Password  string            `bson:"password"`
      Email     string            `bson:"email"`
      Status    int              `bson:"status"`
      CreatedAt time.Time         `bson:"created_at"`
      UpdatedAt time.Time         `bson:"updated_at"`
  }
  
  // 视频集合
  type Video struct {
      ID           primitive.ObjectID `bson:"_id,omitempty"`
      UserID       string            `bson:"user_id"`
      Title        string            `bson:"title"`
      Description  string            `bson:"description"`
      FileName     string            `bson:"file_name"`
      FileSize     int64             `bson:"file_size"`
      Format       string            `bson:"format"`
      Duration     float64           `bson:"duration"`
      Status       string            `bson:"status"`
      Tags         []string          `bson:"tags"`
      Stats        VideoStats        `bson:"stats"`
      CreatedAt    time.Time         `bson:"created_at"`
      UpdatedAt    time.Time         `bson:"updated_at"`
  }
  ```

#### 2.4 缓存策略
- Redis缓存实现
  ```go
  // 缓存配置
  type CacheConfig struct {
      Enabled        bool
      ExpireTime     time.Duration
      CleanupInterval time.Duration
  }
  
  // 缓存键设计
  const (
      VideoDetailKey = "video:detail:%s"    // 视频详情缓存
      VideoListKey   = "video:list:%s"      // 视频列表缓存
      UserProfileKey = "user:profile:%s"    // 用户信息缓存
  )
  ```

### 3. 性能优化

#### 3.1 数据库优化
- MongoDB索引设计
  ```js
  // 视频集合索引
  db.videos.createIndex({ "user_id": 1 })
  db.videos.createIndex({ "status": 1, "created_at": -1 })
  db.videos.createIndex({ "tags": 1 })
  db.videos.createIndex({ 
      "title": "text", 
      "description": "text" 
  })
  
  // 用户集合索引
  db.users.createIndex({ "username": 1 }, { unique: true })
  db.users.createIndex({ "email": 1 }, { unique: true })
  ```

#### 3.2 并发处理
- 视频处理并发控制
  ```go
  // 并发上传控制
  var (
      uploadSemaphore = make(chan struct{}, 10)  // 最多10个并发上传
      processQueue    = make(chan *VideoTask, 100) // 处理队列
  )
  
  // 异步处理器
  func videoProcessor() {
      for task := range processQueue {
          // 处理视频转码、生成缩略图等
      }
  }
  ```

#### 3.3 错误处理
- 统一错误处理
  ```go
  // 错误码定义
  const (
      ErrCodeSuccess        = 0
      ErrCodeParamInvalid   = 400
      ErrCodeUnauthorized   = 401
      ErrCodeForbidden      = 403
      ErrCodeNotFound       = 404
      ErrCodeServerError    = 500
  )
  
  // 错误响应结构
  type ErrorResponse struct {
      Code    int    `json:"code"`
      Message string `json:"msg"`
      Details string `json:"details,omitempty"`
  }
  ```

### 4. 安全措施

#### 4.1 密码加密
```go
// 密码加密
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// 密码验证
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

#### 4.2 输入验证
```go
// 用户注册请求验证
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=32,username"`
    Password string `json:"password" binding:"required,min=6,max=32,password"`
    Email    string `json:"email" binding:"required,email"`
}

// 自定义验证器
func usernameValidator(fl validator.FieldLevel) bool {
    username := fl.Field().String()
    match, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
    return match
}
```

#### 4.3 文件上传安全
```go
// 文件类型验证
func validateFileType(file *multipart.FileHeader) bool {
    // 检查文件MIME类型
    buffer := make([]byte, 512)
    f, _ := file.Open()
    defer f.Close()
    n, _ := f.Read(buffer)
    fileType := http.DetectContentType(buffer[:n])
    
    return allowedTypes[fileType]
}
```

### 5. 监控和日志

#### 5.1 日志配置
```go
// 日志配置
type LogConfig struct {
    Level      string    // 日志级别
    Filename   string    // 日志文件
    MaxSize    int      // 单个文件最大尺寸
    MaxBackups int      // 最大保留文件数
    MaxAge     int      // 最大保留天数
    Compress   bool     // 是否压缩
}

// 结构化日志示例
log.Info("视频上传成功",
    "videoId", video.ID,
    "userId", user.ID,
    "fileSize", video.FileSize,
    "duration", video.Duration,
)
```

#### 5.2 性能监控
```go
// Prometheus指标
var (
    uploadLatency = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "video_upload_latency_seconds",
            Help: "视频上传耗时分布",
        },
        []string{"status"},
    )
    
    activeConnections = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_connections",
            Help: "当前活动连接数",
        },
    )
)
```