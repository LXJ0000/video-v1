# 视频管理平台后端服务

## 项目简介
该项目是一个基于Go语言和Gin框架开发的视频管理平台后端服务，提供视频文件的基础管理功能，包括上传、获取和管理等功能。项目采用清晰的分层架构，遵循RESTful API设计规范，并实现了统一的错误处理和响应格式。

## 主要功能
### 视频上传
- 视频上传（支持断点续传）
- 支持大文件上传（最大1GB）
- 支持主流视频格式（mp4、mov、avi等）
- 上传时可设置视频标题和描述信息

### 视频获取
- 视频流式播放
- 视频列表查询（支持分页、排序、搜索）
- 视频详情获取

### 视频管理
- 视频管理（列表、详情、修改、删除）
- 支持批量删除
- 支持修改视频基本信息

### 用户标记、注释和笔记功能
- 用户可以在视频播放过程中添加标记
  - 支持添加、更新、删除标记
  - 可以为标记添加时间戳和内容
- 用户可以为每个标记添加文字注释
  - 支持添加、更新、删除注释
  - 注释与标记关联
- 用户可以在视频播放过程中随时记录笔记
  - 支持添加、更新、删除笔记
  - 可以为笔记添加时间戳和内容
- 支持导出所有标记、注释和笔记
  - 导出为文本格式
  - 包含时间戳和内容信息

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
│   │   └── video.go     # 视频处理器
│   ├── model/     # 数据模型
│   │   └── video.go     # 视频模型
│   └── service/   # 业务逻辑
│       ├── video.go     # 视频服务
│       └── video_test.go # 视频服务测试
├── pkg/           # 可重用的包
│   ├── database/  # 数据库相关
│   │   ├── mongodb.go   # MongoDB连接
│   │   └── mongodb_test.go # MongoDB测试
│   └── response/  # 响应处理
│       └── response.go  # 统一响应格式
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

## 错误码说明
- 400: 请求参数错误
- 401: 未授权
- 403: 权限不足
- 404: 资源不存在
- 413: 文件太大
- 415: 不支持的文件格式
- 500: 服务器内部错误

## 注意事项
1. 文件上传限制
   - 视频文件最大 1GB
   - 封面图最大 2MB（1280x720 推荐）
   - 缩略图由系统自动生成
   - 支持的视频格式: mp4, mov, avi, wmv, flv, mkv
   - 支持的图片格式: jpg, jpeg, png, gif

2. 封面图和缩略图
   - 封面图（Cover）
     - 用户上传的视频展示封面
     - 支持自定义设计，可包含文字和logo
     - 用于视频列表和详情页面展示
     - 建议尺寸：1280x720

   - 缩略图（Thumbnail）
     - 系统自动从视频中生成
     - 用于视频预览和进度条预览
     - 自动截取视频帧
     - 尺寸：320x180

3. 视频状态
   - public: 公开
   - private: 私有
   - draft: 草稿

4. 权限要求
   - 需要登录才能进行上传、修改、删除操作
   - 只有视频作者可以修改和删除视频
   - 批量操作会验证每个视频的权限

## 视频状态说明
视频有三种状态：
1. `public`（公开）
   - 所有用户可见
   - 会出现在视频列表中
   - 可以被搜索到

2. `private`（私有）
   - 仅作者可见
   - 对作者显示在视频列表中
   - 其他用户无法访问
   - 新上传的视频默认为此状态

3. `draft`（草稿）
   - 仅作者可见
   - 对作者显示在视频列表中
   - 其他用户无法访问
   - 通常用于未完成编辑的视频

状态转换规则：
- 任何状态都可以互相转换
- 只有视频作者可以更改状态
- 批量操作时会验证每个视频的权限

## 工具说明

### 数据库清理工具
位于 `scripts/clean_db.go`，用于清理数据库中的历史数据。

#### 编译
```bash
go build -o clean_db scripts/clean_db.go
```

#### 功能
- 支持清理多种类型的数据（视频、标记、注释、笔记）
- 支持按时间清理数据
- 支持只清理测试数据
- 支持试运行模式
- 支持单独确认每个集合的删除操作

#### 使用方法
```bash
# 查看帮助信息
./clean_db -h

# 清理选项：
  -all         清理所有数据
  -days n      清理n天前的数据（默认30天）
  -dry-run     试运行模式，不实际删除数据
  -test-only   只清理测试数据
  -type string 要清理的数据类型：
               all: 所有数据（默认）
               videos: 视频数据
               marks: 标记数据
               annotations: 注释数据
               notes: 笔记数据

# 使用示例：
# 清理30天前的所有数据
./clean_db -days 30

# 清理所有测试数据
./clean_db -test-only

# 试运行模式查看将要删除的数据
./clean_db -dry-run

# 只清理视频数据
./clean_db -type videos
```

#### 注意事项
- 执行清理操作前建议先使用 `-dry-run` 选项预览要删除的数据
- 清理操作不可恢复，请谨慎使用
- 建议在低峰期执行清理操作
- 每个集合的删除操作都需要单独确认 