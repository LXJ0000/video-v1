# 视频管理平台 API 文档

## 基础信息
- 基础路径: `/api/v1`
- 认证方式: Bearer Token
- 响应格式: JSON

## 通用响应格式
```json
{
    "code": 0,       // 0表示成功，非0表示失败
    "msg": "string", // 响应消息
    "data": {}       // 响应数据
}
```

## 认证相关

### 用户注册
- 请求方式: `POST`
- 路径: `/users/register`
- 请求体:
```json
{
    "username": "string",  // 用户名
    "password": "string",  // 密码
    "email": "string"     // 邮箱
}
```
- 响应示例:
```json
{
    "code": 0,
    "msg": "注册成功",
    "data": {
        "id": "string",
        "username": "string",
        "email": "string",
        "createdAt": "2024-02-26T10:00:00Z"
    }
}
```

### 用户登录
- 请求方式: `POST`
- 路径: `/users/login`
- 请求体:
```json
{
    "username": "string", // 用户名
    "password": "string"  // 密码
}
```
- 响应示例:
```json
{
    "code": 0,
    "msg": "登录成功",
    "data": {
        "token": "string",
        "user": {
            "id": "string",
            "username": "string",
            "email": "string"
        }
    }
}
```

## 视频相关

### 获取公开视频列表
- 请求方式: `GET`
- 路径: `/videos/public`
- 查询参数:
  - `page`: 页码，默认1
  - `pageSize`: 每页数量，默认10
  - `keyword`: 搜索关键词（可选）
  - `sortBy`: 排序字段（可选，支持：created_at/views/likes）
  - `sortOrder`: 排序方向（可选，asc/desc）
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "total": 100,
        "page": 1,
        "pageSize": 10,
        "items": [{
            "id": "string",
            "title": "string",
            "description": "string",
            "coverUrl": "string",
            "duration": 180.5,
            "status": "public",
            "views": 1000,
            "likes": 100,
            "createdAt": "2024-02-26T10:00:00Z"
        }]
    }
}
```

### 获取视频列表（需要认证）
- 请求方式: `GET`
- 路径: `/videos`
- 查询参数:
  - `page`: 页码，默认1
  - `pageSize`: 每页数量，默认10
  - `status`: 视频状态（可选，支持：public/private/draft）
  - `keyword`: 搜索关键词（可选）
  - `sortBy`: 排序字段（可选）
  - `sortOrder`: 排序方向（可选）
- 响应格式同上

### 上传视频
- 请求方式: `POST`
- 路径: `/videos`
- Content-Type: `multipart/form-data`
- 请求参数:
  - `file`: 视频文件
  - `cover`: 封面图（可选）
  - `title`: 标题
  - `description`: 描述（可选）
  - `status`: 状态（public/private/draft）
  - `duration`: 时长（秒）
  - `tags`: 标签（可选，逗号分隔）

### 获取视频详情
- 请求方式: `GET`
- 路径: `/videos/:id`
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "id": "string",
        "title": "string",
        "description": "string",
        "coverUrl": "string",
        "duration": 180.5,
        "status": "public",
        "views": 1000,
        "likes": 100,
        "createdAt": "2024-02-26T10:00:00Z",
        "updatedAt": "2024-02-26T10:00:00Z"
    }
}
```

### 更新视频信息
- 请求方式: `PUT`
- 路径: `/videos/:id`
- 请求体:
```json
{
    "title": "string",
    "description": "string",
    "status": "string",
    "tags": ["string"]
}
```

### 删除视频
- 请求方式: `DELETE`
- 路径: `/videos/:id`

### 更新视频缩略图
- 请求方式: `POST`
- 路径: `/videos/:id/thumbnail`
- Content-Type: `multipart/form-data`
- 请求参数:
  - `file`: 图片文件（支持jpg/jpeg/png，最大2MB）

### 获取视频统计信息
- 请求方式: `GET`
- 路径: `/videos/:id/stats`
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "views": 1000,
        "likes": 100,
        "comments": 50
    }
}
```

## 标记相关

### 添加标记
- 请求方式: `POST`
- 路径: `/marks`
- 请求体:
```json
{
    "videoId": "string",
    "timestamp": 123.45,
    "content": "string"
}
```

### 获取标记列表
- 请求方式: `GET`
- 路径: `/marks`
- 查询参数:
  - `videoId`: 视频ID

### 更新标记
- 请求方式: `PUT`
- 路径: `/marks/:markId`
- 请求体:
```json
{
    "content": "string"
}
```

### 删除标记
- 请求方式: `DELETE`
- 路径: `/marks/:markId`

### 添加注释
- 请求方式: `POST`
- 路径: `/marks/:markId/annotations`
- 请求体:
```json
{
    "content": "string"
}
```

### 获取注释列表
- 请求方式: `GET`
- 路径: `/marks/:markId/annotations`

### 更新注释
- 请求方式: `PUT`
- 路径: `/marks/annotations/:annotationId`
- 请求体:
```json
{
    "content": "string"
}
```

### 删除注释
- 请求方式: `DELETE`
- 路径: `/marks/annotations/:annotationId`

## 笔记相关

### 添加笔记
- 请求方式: `POST`
- 路径: `/notes`
- 请求体:
```json
{
    "videoId": "string",
    "timestamp": 123.45,
    "content": "string"
}
```

### 获取笔记列表
- 请求方式: `GET`
- 路径: `/notes`
- 查询参数:
  - `videoId`: 视频ID

### 更新笔记
- 请求方式: `PUT`
- 路径: `/notes/:noteId`
- 请求体:
```json
{
    "content": "string"
}
```

### 删除笔记
- 请求方式: `DELETE`
- 路径: `/notes/:noteId`

## 导出相关

### 导出标记和笔记
- 请求方式: `GET`
- 路径: `/videos/export`
- 查询参数:
  - `videoId`: 视频ID
- 响应: 文件下载（CSV格式） 