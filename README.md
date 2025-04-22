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

## 测试架构与指南

### 测试架构概述

本项目采用分层测试策略，主要包括以下几个层次的测试：

1. **单元测试** - 测试独立的函数和方法
2. **集成测试** - 测试多个组件之间的交互
3. **处理器测试** - 测试HTTP请求处理逻辑
4. **端到端测试** - 测试整个API流程

### 测试文档

项目提供了详细的测试文档，帮助开发人员编写和维护测试：

- [测试指南](docs/testing.md) - 详细说明测试架构、编写方法和最佳实践
- [测试问题排查](docs/test-troubleshooting.md) - 常见测试问题和解决方案

### 单元测试框架

- 使用Go标准的`testing`包作为基础测试框架
- 使用`github.com/stretchr/testify/assert`简化断言
- 使用`github.com/stretchr/testify/mock`进行模拟

### 服务层测试

服务层测试主要测试业务逻辑，模拟数据库操作。以`service/user_test.go`为例：

```go
// 示例：测试用户登录功能
func TestUserLoginSuccess(t *testing.T) {
    // 准备测试数据
    userID := primitive.NewObjectID()
    user := &model.User{
        ID:        userID,
        Username:  "testuser",
        Password:  hashedPassword,
        Email:     "test@example.com",
        Status:    1,
    }
    
    // 验证预期结果
    assert.Equal(t, user.Username, "testuser")
    assert.Equal(t, user.Email, "test@example.com")
}
```

#### 模拟数据库

服务层测试中，使用`MockCollection`模拟MongoDB集合：

```go
// 模拟Collection接口
type MockCollection struct {
    mock.Mock
    *mongo.Collection
}

func (m *MockCollection) FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult {
    args := m.Called(ctx, filter)
    return args.Get(0).(*mongo.SingleResult)
}

// ... 其他数据库方法模拟 ...
```

### 处理器层测试

处理器测试使用Gin的测试工具和模拟服务层。以`handler/user_test.go`为例：

```go
// 测试获取用户详情
func TestGetUserProfile(t *testing.T) {
    c, w, mockService, handler := setupUserTest()

    // 模拟当前登录用户
    userId := primitive.NewObjectID().Hex()
    c.Set("userId", userId)
    c.Params = []gin.Param{{Key: "userId", Value: "me"}}
    c.Request = httptest.NewRequest("GET", "/", nil)  // 确保请求对象不为nil

    // ... 模拟服务响应与执行测试 ...

    // 验证响应
    assert.Equal(t, http.StatusOK, w.Code)
}
```

#### 模拟服务层

处理器测试中，使用`MockUserService`模拟服务层接口：

```go
// 创建UserService的Mock
type MockUserService struct {
    mock.Mock
}

func (m *MockUserService) GetUserProfile(ctx context.Context, id string) (*model.UserProfileResponse, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*model.UserProfileResponse), args.Error(1)
}

// ... 其他服务方法模拟 ...
```

### 测试的注意事项

1. **请求对象初始化**
   - 在Gin测试上下文中，确保`c.Request`不为nil
   - 使用`httptest.NewRequest`创建请求对象
   - 示例：`c.Request = httptest.NewRequest("GET", "/", nil)`

2. **模拟数据清理**
   - 测试后清理模拟的数据和连接
   - 使用`defer`确保资源释放
   - 示例：`defer cursor.Close(ctx)`

3. **避免测试间依赖**
   - 每个测试应该独立运行
   - 不要依赖测试的执行顺序
   - 使用`t.Parallel()`支持并行测试

4. **边界条件测试**
   - 测试空值、错误输入和边界情况
   - 验证错误处理逻辑
   - 测试权限验证失败情况

### 运行测试指南

1. **运行所有测试**
   ```bash
   go test ./...
   ```

2. **运行特定包的测试**
   ```bash
   go test ./internal/service
   go test ./internal/handler
   ```

3. **运行单个测试函数**
   ```bash
   go test ./internal/handler -run TestGetUserProfile
   ```

4. **带详细输出的测试**
   ```bash
   go test -v ./internal/handler
   ```

5. **强制不使用缓存的测试**
   ```bash
   go test -count=1 ./...
   ```

6. **生成测试覆盖率报告**
   ```bash
   go test -cover ./...
   go test -coverprofile=coverage.out ./...
   go tool cover -html=coverage.out
   ```

7. **测试超时设置**
   ```bash
   go test -timeout 30s ./...
   ```

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
  - ✅ 支持获取基本信息和播放地址
  - ✅ 支持获取统计数据
  - ✅ 支持获取当前用户的收藏状态

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

### 2024-05-01 视频级联删除功能优化
1. **问题与需求**
   - **问题分析**：删除视频时，相关的收藏、观看历史等关联数据没有被同步删除，导致数据库中留下大量"悬空"引用
   - **需求**：实现视频删除时的级联删除功能，确保数据一致性
   - **目标**：提高系统的数据完整性，避免垃圾数据积累

2. **功能实现**
   - **事务支持**：使用MongoDB事务确保所有删除操作的原子性
   - **级联删除内容**：
     - 视频记录本身
     - 相关收藏记录（favorites集合）
     - 观看历史记录（watch_history集合）
     - 视频评论（comments集合）
     - 视频标记（marks集合）
     - 视频注释（annotations集合）
   - **文件系统操作**：
     - 删除视频文件
     - 删除缩略图文件
     - 错误处理优化，即使文件删除失败也不影响主流程

3. **技术亮点**
   - **事务隔离**：使用MongoDB事务保证数据删除的一致性
   - **容错设计**：文件系统操作与数据库操作分离，避免文件系统错误影响数据库事务
   - **日志记录**：使用结构化日志记录删除过程，便于问题排查
   - **细粒度错误处理**：针对每个删除步骤提供详细的错误信息

4. **测试策略**
   - 创建`TestDeleteVideoCascade`测试用例框架
   - 由于涉及事务和文件系统操作，需要在集成测试环境中完整测试
   - 测试步骤包括：
     - 预创建测试数据（视频、收藏、历史等）
     - 调用删除方法
     - 验证所有关联记录和文件是否正确删除

5. **代码示例**
   ```go
   // 在事务中执行所有删除操作
   _, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
       // 1. 删除相关收藏记录
       _, err := database.GetCollection("favorites").DeleteMany(
           sessCtx,
           bson.M{"video_id": id},
       )
       // ... 其他删除操作
       return nil, nil
   })
   ```

6. **技术挑战与解决方案**
   - **挑战**：MongoDB事务需要副本集支持
   - **解决方案**：确保生产环境使用MongoDB副本集配置
   - **挑战**：文件删除失败不应影响主流程
   - **解决方案**：将文件删除操作从事务中分离，仅记录错误不中断流程

### 2024-04-21 用户收藏状态功能实现
1. 视频详情功能增强
   - 新增返回用户收藏状态
   - 优化查询性能
   - 添加单元测试
   - 更新API文档

2. 用户服务改进
   - 新增`CheckFavoriteStatus`方法
   - 优化数据库查询
   - 增加测试用例

3. 测试架构优化
   - 修复依赖注入问题
   - 优化模拟对象
   - 提高测试覆盖率

4. API文档更新
   - 添加视频收藏状态说明
   - 完善响应示例
   - 更新接口参数说明

### 2024-04-28 收藏与点赞同步功能实现
1. **问题与需求分析**
   - **发现问题**：视频的收藏数据和点赞数据不同步，造成统计不准确
   - **原因**：收藏行为没有影响点赞计数
   - **目标**：实现收藏与点赞计数同步，提高用户体验和数据一致性

2. **功能实现**
   - **收藏时增加点赞数**：用户添加收藏时自动增加视频点赞计数
   - **取消收藏时减少点赞数**：用户取消收藏时自动减少视频点赞计数
   - **事务支持**：使用MongoDB事务确保数据一致性
   - **安全保障**：添加条件判断，确保点赞数不会小于0

3. **技术实现**
   - **MongoDB事务**：使用事务包装收藏和点赞更新操作
   - **原子性操作**：使用`$inc`操作符增减点赞数
   - **条件更新**：使用条件查询确保点赞数不会变为负数
   - **错误处理**：详细的错误报告和回滚机制

4. **技术亮点**
   ```go
   // 添加收藏时增加点赞数
   _, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
       // 1. 添加到收藏表
       _, err := collection.InsertOne(sessCtx, favorite)
       if err != nil {
           return nil, err
       }

       // 2. 更新视频的likes统计
       update := bson.M{
           "$inc": bson.M{
               "stats.likes": 1,
           },
       }
       _, err = videoCollection.UpdateOne(sessCtx, bson.M{"_id": objectID}, update)
       return nil, err
   })
   
   // 取消收藏时减少点赞数，同时确保不会小于0
   filter := bson.M{
       "_id":         objectID,
       "stats.likes": bson.M{"$gt": 0}, // 确保likes大于0
   }
   
   update := bson.M{
       "$inc": bson.M{
           "stats.likes": -1,
       },
   }
   
   _, err = videoCollection.UpdateOne(sessCtx, filter, update)
   ```

5. **测试优化**
   - 添加事务模拟测试
   - 验证点赞数增减逻辑
   - 边界条件测试（如点赞数为0时的行为）
   - 错误处理和恢复测试

6. **API接口优化**
   - `/videos/:videoId/favorite` POST：添加收藏，同时增加点赞
   - `/videos/:videoId/favorite` DELETE：取消收藏，同时减少点赞
   - 视频详情接口返回`isFavorite`字段，标识用户收藏状态

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
        "video": {
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
        },
        "isFavorite": true  // 当用户已登录时，返回用户是否已收藏此视频
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

## 系统架构图

### 前端架构图
```
@startuml
package "前端应用" {
  [Vue 3 应用] as VueApp
  [Pinia 状态管理] as Pinia
  [Vue Router] as Router
  [Element Plus UI] as ElementUI
  [Video.js 播放器] as VideoJS
  
  VueApp --> Pinia : 使用
  VueApp --> Router : 使用
  VueApp --> ElementUI : 使用
  VueApp --> VideoJS : 使用
  
  package "视频模块" as VideoModule {
    [视频播放组件] as VideoPlayer
    [视频列表组件] as VideoList
    [视频上传组件] as VideoUpload
    
    VideoPlayer --> VideoJS : 基于
  }
  
  package "标记笔记模块" as MarkModule {
    [标记编辑器] as MarkEditor
    [标记列表] as MarkList
    [笔记编辑器] as NoteEditor
    [导出工具] as ExportTool
  }
  
  package "用户模块" as UserModule {
    [用户认证] as UserAuth
    [个人资料] as UserProfile
    [收藏管理] as Favorites
  }
  
  VueApp --> VideoModule : 包含
  VueApp --> MarkModule : 包含
  VueApp --> UserModule : 包含
}

package "API请求" {
  [Axios客户端] as Axios
  [请求拦截器] as RequestInterceptor
  [响应拦截器] as ResponseInterceptor
  
  Axios --> RequestInterceptor : 使用
  Axios --> ResponseInterceptor : 使用
}

VueApp --> Axios : HTTP请求

cloud "后端服务" {
  [RESTful API] as API
}

Axios --> API : 发送请求
@enduml
```

### 数据流图
```
@startuml
actor 用户
participant "前端应用" as Frontend
participant "API网关" as API
participant "用户服务" as UserService
participant "视频服务" as VideoService
participant "标记服务" as MarkService
database "MongoDB" as DB
collections "Redis缓存" as Cache

== 用户登录流程 ==
用户 -> Frontend: 输入用户名和密码
Frontend -> API: 发送登录请求
API -> UserService: 验证用户凭据
UserService -> DB: 查询用户信息
DB -> UserService: 返回用户数据
UserService -> UserService: 验证密码
UserService -> API: 生成JWT令牌
API -> Frontend: 返回令牌和用户信息
Frontend -> Frontend: 存储令牌

== 视频播放流程 ==
用户 -> Frontend: 选择视频
Frontend -> API: 请求视频详情
API -> VideoService: 获取视频信息
VideoService -> Cache: 查询缓存
Cache -> VideoService: 返回缓存结果(未命中)
VideoService -> DB: 查询视频数据
DB -> VideoService: 返回视频信息
VideoService -> Cache: 更新缓存
VideoService -> API: 返回视频详情
API -> Frontend: 返回视频详情
Frontend -> API: 请求视频流
API -> VideoService: 流式传输视频
VideoService -> Frontend: 返回视频流
Frontend -> Frontend: 播放视频

== 添加标记流程 ==
用户 -> Frontend: 在视频上创建标记
Frontend -> API: 发送标记数据
API -> MarkService: 创建标记
MarkService -> DB: 存储标记数据
DB -> MarkService: 确认存储成功
MarkService -> API: 返回标记信息
API -> Frontend: 返回标记详情
Frontend -> Frontend: 更新标记列表
@enduml
```

### 后端处理流程图
```
@startuml
participant "客户端" as Client
participant "Gin路由" as Gin
participant "认证中间件" as Auth
participant "处理器" as Handler
participant "服务层" as Service
participant "MongoDB" as DB
participant "Redis缓存" as Cache
participant "文件系统" as FS

== 视频上传流程 ==
Client -> Gin: 发送视频上传请求(POST /videos)
Gin -> Auth: 验证请求token
Auth -> Gin: 验证通过，添加userId到上下文
Gin -> Handler: 调用上传处理器
Handler -> Handler: 检查请求参数和文件类型
Handler -> Service: 调用视频服务上传方法
Service -> FS: 保存视频文件
FS -> Service: 返回文件路径
Service -> DB: 创建视频记录
DB -> Service: 返回新创建的视频ID
Service -> Handler: 返回上传结果
Handler -> Client: 返回成功响应(视频详情)

== 视频流播放流程 ==
Client -> Gin: 请求视频流(GET /videos/:id/stream)
Gin -> Handler: 调用流处理器
Handler -> Service: 获取视频信息
Service -> DB: 查询视频记录
DB -> Service: 返回视频数据
Service -> Handler: 返回视频信息
Handler -> FS: 打开视频文件
Handler -> Handler: 处理Range请求头
Handler -> Client: 流式返回视频数据(206 Partial Content)

== 收藏视频流程 ==
Client -> Gin: 添加收藏请求(POST /videos/:id/favorite)
Gin -> Auth: 验证请求token
Auth -> Gin: 验证通过，添加userId到上下文
Gin -> Handler: 调用收藏处理器
Handler -> Service: 调用用户服务添加收藏方法
Service -> DB: 开始数据库事务
DB -> Service: 返回事务会话
Service -> DB: 添加收藏记录
Service -> DB: 更新视频点赞计数
DB -> Service: 提交事务
Service -> Handler: 返回操作结果
Handler -> Client: 返回成功响应
@enduml
```

## 技术实现重点

### MongoDB事务支持实现

为确保数据一致性，项目中的关键操作(如添加/取消收藏同步更新点赞数)使用MongoDB事务：

```go
// 使用事务包装收藏和点赞操作的示例
func (s *UserService) AddToFavorites(ctx context.Context, videoID, userID string) error {
    // 1. 开始事务会话
    session, err := database.GetClient().StartSession()
    if err != nil {
        return err
    }
    defer session.EndSession(ctx)

    // 2. 在事务内执行多个操作
    _, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
        // 2.1 验证视频存在性
        video, err := s.getVideoByID(sessCtx, videoID)
        if err != nil {
            return nil, err
        }

        // 2.2 检查是否已收藏
        count, err := database.GetCollection("favorites").CountDocuments(
            sessCtx,
            bson.M{"user_id": userID, "video_id": videoID},
        )
        if err != nil {
            return nil, err
        }
        if count > 0 {
            return nil, errors.New("已经收藏过该视频")
        }

        // 2.3 添加收藏记录
        favorite := model.Favorite{
            ID:        primitive.NewObjectID(),
            UserID:    userID,
            VideoID:   videoID,
            CreatedAt: time.Now(),
        }
        _, err = database.GetCollection("favorites").InsertOne(sessCtx, favorite)
        if err != nil {
            return nil, err
        }

        // 2.4 增加视频点赞计数
        _, err = database.GetCollection("videos").UpdateOne(
            sessCtx,
            bson.M{"_id": video.ID},
            bson.M{"$inc": bson.M{"stats.likes": 1}},
        )
        return nil, err
    })

    return err
}
```

### 视频流式传输与断点续传实现

视频流式传输使用HTTP Range请求实现，支持断点续传和播放进度控制：

```go
// 视频流式播放的核心实现
func (h *VideoHandler) Stream(c *gin.Context) {
    videoID := c.Param("videoId")
    
    // 获取视频信息
    video, err := h.videoService.GetByID(c, videoID)
    if err != nil {
        response.Fail(c, "视频不存在")
        return
    }
    
    // 获取文件路径
    filePath := filepath.Join(config.GlobalConfig.Storage.UploadDir, video.FileName)
    
    // 打开文件
    file, err := os.Open(filePath)
    if err != nil {
        response.Fail(c, "视频文件不存在")
        return
    }
    defer file.Close()
    
    // 获取文件信息
    fileInfo, err := file.Stat()
    if err != nil {
        response.Fail(c, "无法获取文件信息")
        return
    }
    
    // 设置内容类型
    c.Header("Content-Type", fmt.Sprintf("video/%s", video.Format))
    c.Header("Accept-Ranges", "bytes")
    
    // 处理Range请求
    rangeHeader := c.Request.Header.Get("Range")
    if rangeHeader != "" {
        // 解析Range头
        ranges, err := parseRange(rangeHeader, fileInfo.Size())
        if err != nil {
            c.Status(http.StatusRequestedRangeNotSatisfiable)
            return
        }
        
        // 只处理第一个范围请求
        if len(ranges) > 0 {
            start, end := ranges[0][0], ranges[0][1]
            
            // 设置部分内容响应头
            c.Status(http.StatusPartialContent)
            c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileInfo.Size()))
            c.Header("Content-Length", fmt.Sprintf("%d", end-start+1))
            
            // 定位到指定位置
            file.Seek(start, io.SeekStart)
            
            // 限制读取长度
            io.CopyN(c.Writer, file, end-start+1)
            return
        }
    }
    
    // 如果不是Range请求，发送整个文件
    c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
    io.Copy(c.Writer, file)
}
```

### 模块化设计与依赖注入

项目采用接口驱动设计，服务层通过接口实现依赖注入，便于模块解耦和单元测试：

```go
// 定义服务接口
type VideoService interface {
    Upload(ctx context.Context, videoFile *multipart.FileHeader, coverFile *multipart.FileHeader, info model.Video) (*model.Video, error)
    GetList(ctx context.Context, query model.VideoQuery) (*model.VideoList, error)
    GetByID(ctx context.Context, id string) (*model.Video, error)
    Update(ctx context.Context, id string, video model.Video) error
    Delete(ctx context.Context, id string) error
    // ... 其他方法
}

// 处理器依赖服务接口，而非具体实现
type VideoHandler struct {
    videoService VideoService
    userService  UserService
}

// 通过构造函数注入依赖
func NewVideoHandler(videoService VideoService, userService UserService) *VideoHandler {
    return &VideoHandler{
        videoService: videoService,
        userService:  userService,
    }
}

// 路由初始化中使用依赖注入
func InitRoutes(r *gin.Engine) {
    // 创建服务实例
    userService := service.NewUserService()
    videoService := service.NewVideoService()
    
    // 创建处理器实例，注入服务依赖
    userHandler := handler.NewUserHandler(userService)
    videoHandler := handler.NewVideoHandler(videoService, userService)
    
    // 配置路由
    api := r.Group("/api/v1")
    // ... 路由配置
}
```

### Redis缓存优化

为提高系统性能，项目对热点数据（如视频详情、用户信息）实现了Redis缓存：

```go
// 缓存服务接口定义
type CacheService interface {
    Get(ctx context.Context, key string, value interface{}) error
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
    Delete(ctx context.Context, key string) error
}

// 带缓存的视频详情获取实现
func (s *videoService) GetByID(ctx context.Context, id string) (*model.Video, error) {
    var video model.Video
    
    // 生成缓存键
    cacheKey := fmt.Sprintf("video:%s", id)
    
    // 尝试从缓存获取
    err := s.cache.Get(ctx, cacheKey, &video)
    if err == nil {
        // 缓存命中，记录指标并返回
        metrics.CacheHitCounter.Inc()
        return &video, nil
    }
    
    // 缓存未命中，从数据库查询
    metrics.CacheMissCounter.Inc()
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }
    
    err = database.GetCollection("videos").FindOne(
        ctx,
        bson.M{"_id": objectID},
    ).Decode(&video)
    if err != nil {
        return nil, err
    }
    
    // 写入缓存，过期时间30分钟
    _ = s.cache.Set(ctx, cacheKey, video, 30*time.Minute)
    
    return &video, nil
}
```

## 项目开发经验总结

### 开发过程中的关键决策

1. **MongoDB vs 关系型数据库**：选择MongoDB的主要考量点是：
   - 文档模型适合存储视频元数据和标记等嵌套结构数据
   - Schema灵活性，便于快速迭代和需求变更
   - 原生支持JSON格式，减少数据转换开销

2. **Gin框架选择**：使用Gin框架是基于性能和开发效率的权衡：
   - 高性能HTTP路由
   - 中间件机制设计合理
   - 强大的参数验证和绑定功能
   - 社区活跃度高，文档完善

3. **视频文件存储策略**：
   - 当前使用本地文件系统存储，便于开发和测试
   - 预留对象存储接口，后续可无缝迁移到S3兼容对象存储
   - 区分元数据和文件存储，为后续微服务拆分做准备

### 后续发展规划

1. **微服务架构演进**：
   - 将用户服务、视频服务、标记服务拆分为独立微服务
   - 引入服务注册发现组件
   - 实现API网关层，统一认证和请求分发

2. **性能优化方向**：
   - 实现内容分发网络(CDN)集成
   - 添加视频转码服务，支持自适应码率
   - 优化数据访问模式，减少数据库负载

3. **功能扩展计划**：
   - 实现视频评论功能
   - 添加社交分享功能
   - 开发专注于学习场景的智能推荐系统