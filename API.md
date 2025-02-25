# 视频管理平台 API 文档

## 基础信息
- 基础路径: `/api/v1`
- 响应格式: 
  ```json
  {
    "code": 0,      // 0: 成功, 1: 失败
    "msg": "success", 
    "data": {}      // 响应数据
  }
  ```

## 视频管理接口

### 1. 上传视频
#### 请求
- **方法**: `POST`
- **路径**: `/videos`
- **Content-Type**: `multipart/form-data`
- **请求参数**:
  - `file`: 视频文件（必填，支持 mp4/mov/avi/wmv/flv/mkv）
  - `cover`: 封面图文件（可选，支持 jpg/jpeg/png，最大2MB）
  - `title`: 视频标题（必填）
  - `duration`: 视频时长（必填，单位：秒，支持小数点后1位）
  - `description`: 视频描述（可选）
  - `status`: 视频状态（可选，默认为 private）
    - `public`: 公开，所有人可见
    - `private`: 私有，仅作者可见
    - `draft`: 草稿，仅作者可见且不会出现在列表中
  - `tags`: 标签（可选，多个标签用逗号分隔）

#### 响应
##### 成功响应
````json
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
````

##### 失败响应
````json
{
  "code": 1,
  "msg": "文件太大",
  "data": null
}
````

### 2. 获取视频列表
#### 请求
- **方法**: `GET`
- **路径**: `/api/v1/videos`
- **参数**:
  - `page`: 页码（从1开始，默认1）
  - `pageSize`: 每页数量（默认10，最大50）
  - `userId`: 用户ID（可选，指定后只返回该用户的视频）
  - `status`: 视频状态（可选，多个状态用逗号分隔）
    - public: 公开视频
    - private: 私有视频（需要是视频作者）
    - draft: 草稿（需要是视频作者）
  - `sort`: 排序方式（可选，默认 "-created_at"）
    - created_at: 创建时间升序
    - -created_at: 创建时间降序
    - title: 标题升序
    - -title: 标题降序
  - `keyword`: 搜索关键词（可选，搜索标题和描述）

#### 响应
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "total": 100,
    "page": 1,
    "pageSize": 10,
    "items": [
      {
        "id": "string",
        "userId": "string",
        "title": "string",
        "description": "string",
        "coverUrl": "string",
        "duration": 180.5,
        "status": "public",
        "createdAt": "string",
        "updatedAt": "string"
      }
    ]
  }
}
```

#### 说明
1. 访问权限
   - 未登录用户只能看到 public 状态的视频
   - 登录用户可以看到：
     - 所有 public 状态的视频
     - 自己的 private 和 draft 状态的视频
   - 管理员可以看到所有视频

2. 列表过滤
   - 不指定 userId 时：
     - 默认只返回 public 状态的视频
     - 登录用户额外返回自己的 private 和 draft 视频
   - 指定 userId 时：
     - 只返回该用户的视频
     - 需要考虑访问权限（只能看到允许的状态）

3. 状态过滤
   - 可以指定多个状态：`status=public,private`
   - 权限验证：
     - private 和 draft 状态只对视频作者可见
     - 非视频作者指定这些状态会被忽略

4. 搜索说明
   - 关键词搜索范围：标题、描述
   - 支持模糊匹配
   - 不区分大小写
   - 多个关键词用空格分隔（与关系）

### 3. 获取视频详情
#### 请求
- **方法**: `GET`
- **路径**: `/videos/:id`

#### 响应
##### 成功响应
````json
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
````

##### 失败响应
````json
{
  "code": 1,
  "msg": "视频不存在",
  "data": null
}
````

### 4. 更新视频信息
#### 请求
- **方法**: `PUT`
- **路径**: `/videos/:id`
- **Content-Type**: `application/json`
- **请求体**:
````json
{
  "title": "新标题",
  "description": "新描述",
  "status": "public",
  "tags": ["标签1", "标签2"]
}
````

#### 响应
##### 成功响应
````json
{
  "code": 0,
  "msg": "success",
  "data": null
}
````

##### 失败响应
````json
{
  "code": 1,
  "msg": "视频更新失败",
  "data": null
}
````

### 5. 删除视频
#### 请求
- **方法**: `DELETE`
- **路径**: `/videos/:id`

#### 响应
##### 成功响应
````json
{
  "code": 0,
  "msg": "success",
  "data": null
}
````

##### 失败响应
````json
{
  "code": 1,
  "msg": "视频删除失败",
  "data": null
}
````

### 6. 批量操作视频
#### 请求
- **方法**: `POST`
- **路径**: `/videos/batch`
- **Content-Type**: `application/json`
- **请求体**:
````json
{
  "ids": ["视频ID1", "视频ID2"],
  "action": "update_status",  // delete/update_status
  "status": "private"         // 当action为update_status时需要，可选值：public/private/draft
}
````

#### 响应
##### 成功响应
````json
{
  "code": 0,
  "msg": "success",
  "data": {
    "successCount": 2,
    "failedCount": 0,
    "failedIds": []
  }
}
````

##### 失败响应
````json
{
  "code": 1,
  "msg": "批量操作失败",
  "data": null
}
````

### 7. 更新视频缩略图
#### 请求
- **方法**: `POST`
- **路径**: `/videos/:id/thumbnail`
- **Content-Type**: `multipart/form-data`
- **请求参数**:
  - `file`: 图片文件（必填，支持jpg/png/gif，最大2MB）

#### 响应
##### 成功响应
````json
{
  "code": 0,
  "msg": "success",
  "data": {
    "thumbnailUrl": "缩略图URL"
  }
}
````

##### 失败响应
````json
{
  "code": 1,
  "msg": "缩略图更新失败",
  "data": null
}
````

### 8. 获取视频统计信息
#### 请求
- **方法**: `GET`
- **路径**: `/videos/:id/stats`

#### 响应
##### 成功响应
````json
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
````

##### 失败响应
````json
{
  "code": 1,
  "msg": "视频统计信息获取失败",
  "data": null
}
````

### 9. 视频流式播放
#### 请求
- **方法**: `GET`
- **路径**: `/videos/:id/stream`
- **支持范围请求**（Range header）

#### 响应
- Content-Type: video/*
- Accept-Ranges: bytes
- 支持断点续传
- 直接返回视频流

## 标记相关接口

### 10. 添加标记
#### 请求
- **方法**: `POST`
- **路径**: `/marks/:userId/:id`
- **Content-Type**: `application/json`
- **请求体**:
```json
{
  "videoId": "test_video_id",
  "timestamp": 123.45,
  "content": "Test Mark"
}
```

#### 响应
##### 成功响应
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": "标记ID",
    "userId": "test_user_id",
    "videoId": "test_video_id",
    "timestamp": 123.45,
    "content": "Test Mark",
    "createdAt": "2024-01-20T10:00:00Z"
  }
}
```

##### 失败响应
```json
{
  "code": 1,
  "msg": "添加标记失败",
  "data": null
}
```

### 11. 获取标记列表
#### 请求
- **方法**: `GET`
- **路径**: `/marks/:userId/:id`

#### 响应
##### 成功响应
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "id": "标记ID",
      "userId": "test_user_id",
      "videoId": "test_video_id",
      "timestamp": 123.45,
      "content": "Test Mark",
      "annotations": [
        {
          "id": "注释ID",
          "userId": "test_user_id",
          "markId": "标记ID",
          "content": "Test Annotation",
          "createdAt": "2024-01-20T10:00:00Z",
          "updatedAt": "2024-01-20T10:00:00Z"
        }
      ],
      "createdAt": "2024-01-20T10:00:00Z",
      "updatedAt": "2024-01-20T10:00:00Z"
    }
  ]
}
```

##### 失败响应
```json
{
  "code": 1,
  "msg": "获取标记失败",
  "data": null
}
```

### 12. 更新标记
#### 请求
- **方法**: `PUT`
- **路径**: `/marks/:userId/:id/:markId`
- **Content-Type**: `application/json`
- **请求体**:
```json
{
  "content": "Updated Mark",
  "timestamp": 124.5
}
```

#### 响应
##### 成功响应
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": "标记ID",
    "userId": "test_user_id",
    "videoId": "test_video_id",
    "timestamp": 124.5,
    "content": "Updated Mark",
    "createdAt": "2024-01-20T10:00:00Z"
  }
}
```

##### 失败响应
```json
{
  "code": 1,
  "msg": "更新标记失败",
  "data": null
}
```

### 13. 删除标记
#### 请求
- **方法**: `DELETE`
- **路径**: `/marks/:userId/:id/:markId`

#### 响应
##### 成功响应
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

##### 失败响应
```json
{
  "code": 1,
  "msg": "删除标记失败",
  "data": null
}
```

### 14. 添加注释
#### 请求
- **方法**: `POST`
- **路径**: `/marks/:userId/:id/annotations/:markId`
- **Content-Type**: `application/json`
- **请求体**:
```json
{
  "content": "Test Annotation"
}
```

#### 响应
##### 成功响应
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": "注释ID",
    "userId": "test_user_id",
    "markId": "标记ID",
    "content": "Test Annotation",
    "createdAt": "2024-01-20T10:00:00Z",
    "updatedAt": "2024-01-20T10:00:00Z"
  }
}
```

##### 失败响应
```json
{
  "code": 1,
  "msg": "添加注释失败",
  "data": null
}
```

### 15. 获取注释
#### 请求
- **方法**: `GET`
- **路径**: `/marks/:userId/:id/annotations/:markId`

#### 响应
##### 成功响应
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "id": "注释ID",
      "userId": "test_user_id",
      "markId": "标记ID",
      "content": "Test Annotation",
      "createdAt": "2024-01-20T10:00:00Z",
      "updatedAt": "2024-01-20T10:00:00Z"
    }
  ]
}
```

##### 失败响应
```json
{
  "code": 1,
  "msg": "获取注释失败",
  "data": null
}
```

### 16. 更新注释
#### 请求
- **方法**: `PUT`
- **路径**: `/marks/:userId/:id/annotations/:annotationId`
- **Content-Type**: `application/json`
- **请求体**:
```json
{
  "content": "Updated Annotation"
}
```

#### 响应
##### 成功响应
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": "注释ID",
    "userId": "test_user_id",
    "markId": "标记ID",
    "content": "Updated Annotation",
    "createdAt": "2024-01-20T10:00:00Z",
    "updatedAt": "2024-01-20T10:30:00Z"
  }
}
```

##### 失败响应
```json
{
  "code": 1,
  "msg": "更新注释失败",
  "data": null
}
```

### 17. 删除注释
#### 请求
- **方法**: `DELETE`
- **路径**: `/marks/:userId/:id/annotations/:annotationId`

#### 响应
##### 成功响应
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

##### 失败响应
```json
{
  "code": 1,
  "msg": "删除注释失败",
  "data": null
}
```

### 18. 添加笔记
#### 请求
- **方法**: `POST`
- **路径**: `/notes/:userId/:id`
- **Content-Type**: `application/json`
- **请求体**:
```json
{
  "videoId": "test_video_id",
  "timestamp": 123.45,
  "content": "Test Note"
}
```

#### 响应
##### 成功响应
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": "笔记ID",
    "userId": "test_user_id",
    "videoId": "test_video_id",
    "timestamp": 123.45,
    "content": "Test Note",
    "createdAt": "2024-01-20T10:00:00Z",
    "updatedAt": "2024-01-20T10:00:00Z"
  }
}
```

##### 失败响应
```json
{
  "code": 1,
  "msg": "添加笔记失败",
  "data": null
}
```

### 19. 获取笔记列表
#### 请求
- **方法**: `GET`
- **路径**: `/notes/:userId/:id`

#### 响应
##### 成功响应
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "id": "笔记ID",
      "userId": "test_user_id",
      "videoId": "test_video_id",
      "timestamp": 123.45,
      "content": "Test Note",
      "createdAt": "2024-01-20T10:00:00Z",
      "updatedAt": "2024-01-20T10:00:00Z"
    }
  ]
}
```

##### 失败响应
```json
{
  "code": 1,
  "msg": "获取笔记失败",
  "data": null
}
```

### 20. 更新笔记
#### 请求
- **方法**: `PUT`
- **路径**: `/notes/:userId/:id/:noteId`
- **Content-Type**: `application/json`
- **请求体**:
```json
{
  "content": "Updated Note",
  "timestamp": 124.5
}
```

#### 响应
##### 成功响应
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": "笔记ID",
    "userId": "test_user_id",
    "videoId": "test_video_id",
    "timestamp": 124.5,
    "content": "Updated Note",
    "createdAt": "2024-01-20T10:00:00Z",
    "updatedAt": "2024-01-20T10:30:00Z"
  }
}
```

##### 失败响应
```json
{
  "code": 1,
  "msg": "更新笔记失败",
  "data": null
}
```

### 21. 删除笔记
#### 请求
- **方法**: `DELETE`
- **路径**: `/notes/:userId/:id/:noteId`

#### 响应
##### 成功响应
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

##### 失败响应
```json
{
  "code": 1,
  "msg": "删除笔记失败",
  "data": null
}
```

### 22. 导出标记、注释和笔记
#### 请求
- **方法**: `GET`
- **路径**: `/videos/export/:userId/:id`

#### 响应
##### 成功响应
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "marks": [
      {
        "id": "标记ID",
        "userId": "test_user_id",
        "videoId": "test_video_id",
        "timestamp": 123.45,
        "content": "Test Mark",
        "annotations": [
          {
            "id": "注释ID",
            "userId": "test_user_id",
            "markId": "标记ID",
            "content": "Test Annotation",
            "createdAt": "2024-01-20T10:00:00Z",
            "updatedAt": "2024-01-20T10:00:00Z"
          }
        ],
        "createdAt": "2024-01-20T10:00:00Z",
        "updatedAt": "2024-01-20T10:00:00Z"
      }
    ],
    "notes": [
      {
        "id": "笔记ID",
        "userId": "test_user_id",
        "videoId": "test_video_id",
        "timestamp": 123.45,
        "content": "Test Note",
        "createdAt": "2024-01-20T10:00:00Z",
        "updatedAt": "2024-01-20T10:00:00Z"
      }
    ]
  }
}
```

##### 失败响应
```json
{
  "code": 1,
  "msg": "导出失败",
  "data": null
}
```

## 数据管理

### 数据清理
平台提供数据库清理工具，用于维护和清理历史数据。

#### 清理范围
- 视频数据（videos）
- 标记数据（marks）
- 注释数据（annotations）
- 笔记数据（notes）

#### 清理策略
1. 时间策略
   - 可以清理指定天数前的数据
   - 默认清理30天前的数据

2. 数据类型
   - 可以选择清理特定类型的数据
   - 支持清理所有类型数据

3. 测试数据
   - 支持只清理测试数据（user_id 以 "test_" 开头）

#### 使用建议
1. 定期清理
   - 建议每月清理一次30天前的数据
   - 可以通过 cron 任务自动执行

2. 数据备份
   - 清理前确保重要数据已备份
   - 可以使用 MongoDB 的备份工具

3. 执行时间
   - 建议在系统负载较低时执行
   - 可以分批次清理不同类型的数据

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

## 用户系统 API

### 用户注册
#### 请求
- **方法**: `POST`
- **路径**: `/api/v1/users/register`
- **Content-Type**: `application/json`
- **请求参数**:
  ```json
  {
    "username": "string",  // 用户名，3-32个字符
    "password": "string",  // 密码，6-32个字符
    "email": "string"     // 邮箱地址
  }
  ```

#### 响应
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": "string",
    "username": "string",
    "email": "string",
    "createdAt": "string",
    "updatedAt": "string"
  }
}
```

### 用户登录
#### 请求
- **方法**: `POST`
- **路径**: `/api/v1/users/login`
- **Content-Type**: `application/json`
- **请求参数**:
  ```json
  {
    "username": "string",  // 用户名
    "password": "string"   // 密码
  }
  ```

#### 响应
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "user": {
      "id": "string",
      "username": "string",
      "email": "string",
      "createdAt": "string",
      "updatedAt": "string"
    },
    "token": "string"  // JWT token
  }
}
```

### 认证说明
1. Token 格式
   - 使用 JWT (JSON Web Token)
   - 有效期：24小时
   - 格式：`Bearer <token>`

2. 认证方式
   - 在请求头中添加 Authorization
   - 示例：`Authorization: Bearer eyJhbGciOiJIUzI1NiIs...`

3. Token 使用
   - 所有需要认证的 API 都需要在请求头中携带 token
   - token 过期需要重新登录获取
   - 无效的 token 会返回 401 错误

4. 用户状态
   - 正常（status: 1）：可以正常使用所有功能
   - 禁用（status: 2）：无法登录和使用功能

5. 安全建议
   - 密码要求：至少6位，包含字母和数字
   - 建议使用 HTTPS 传输
   - 定期更换密码
   - 不要在客户端明文存储 token

### 错误处理
1. 注册相关错误
   - 用户名已存在：`{"code": 1, "msg": "用户名已存在"}`
   - 邮箱已注册：`{"code": 1, "msg": "邮箱已注册"}`
   - 密码不符合要求：`{"code": 1, "msg": "密码不符合要求"}`

2. 登录相关错误
   - 用户不存在：`{"code": 1, "msg": "用户不存在"}`
   - 密码错误：`{"code": 1, "msg": "密码错误"}`
   - 账号被禁用：`{"code": 1, "msg": "账号已被禁用"}`

3. 认证相关错误
   - token 无效：`{"code": 1, "msg": "无效的token"}`
   - token 过期：`{"code": 1, "msg": "token已过期"}`
   - 未授权访问：`{"code": 1, "msg": "未授权"}` 