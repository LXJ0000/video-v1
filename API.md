# 视频管理平台 API 文档

## 基础信息
- 基础路径: `/api/v1`
- 认证方式: Bearer Token
- 响应格式: JSON
- 服务器地址: `https://api.videoplatform.com`

## 通用响应格式
```json
{
    "code": 0,       // 0表示成功，非0表示失败
    "msg": "string", // 响应消息
    "data": {}       // 响应数据
}
```

## 错误码说明
| 错误码 | 说明 |
|-------|------|
| 0 | 成功 |
| 400 | 参数错误 |
| 401 | 未认证或认证失败 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 认证相关

### 用户注册
- 请求方式: `POST`
- 路径: `/users/register`
- Content-Type: `application/json`
- 请求体:
```json
{
    "username": "string",  // 用户名，3-32个字符，字母、数字、下划线
    "password": "string",  // 密码，6-32个字符，必须包含字母和数字
    "email": "string"      // 邮箱，有效的邮箱格式
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
- 错误情况:
  - 400: 参数不合法
  - 409: 用户名或邮箱已存在

### 用户登录
- 请求方式: `POST`
- 路径: `/users/login`
- Content-Type: `application/json`
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
        "token": "string", // JWT令牌，请在后续请求的Authorization头中使用
        "user": {
            "id": "string",
            "username": "string",
            "email": "string"
        }
    }
}
```
- 错误情况:
  - 400: 参数不合法
  - 401: 用户名或密码错误
  - 403: 账号被禁用

### 发送短信验证码
- 请求方式: `POST`
- 路径: `/users/sms-code`
- Content-Type: `application/json`
- 请求体:
```json
{
    "phone": "string"  // 手机号码，例如：13800138000
}
```
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "message": "验证码已发送"
    }
}
```
- 错误情况:
  - 400: 参数不合法或手机号格式错误
  - 429: 请求过于频繁，请稍后再试

### 短信验证码登录
- 请求方式: `POST`
- 路径: `/users/login/sms`
- Content-Type: `application/json`
- 请求体:
```json
{
    "phone": "string",  // 手机号码
    "code": "string"    // 短信验证码
}
```
- 响应示例:
```json
{
    "code": 0,
    "msg": "登录成功",
    "data": {
        "token": "string", // JWT令牌
        "user": {
            "id": "string",
            "username": "string",
            "phone": "string",
            "email": "string",
            "status": 1
        }
    }
}
```
- 错误情况:
  - 400: 参数不合法或手机号格式错误
  - 401: 验证码错误或已过期
  - 403: 账号被禁用

### 刷新令牌
- 请求方式: `POST`
- 路径: `/users/refresh-token`
- 请求头: `Authorization: Bearer {token}`
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "token": "string",
        "expiresAt": "2024-03-26T10:00:00Z"
    }
}
```

### 获取当前用户信息
- 请求方式: `GET`
- 路径: `/users/me`
- 请求头: `Authorization: Bearer {token}`
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "id": "string",
        "username": "string",
        "email": "string",
        "status": 1,
        "profile": {
            "nickname": "string",
            "avatar": "string",
            "bio": "string"
        },
        "stats": {
            "videos": 10,
            "followers": 20,
            "following": 30
        },
        "createdAt": "2024-02-26T10:00:00Z",
        "updatedAt": "2024-02-26T10:00:00Z"
    }
}
```

### 获取用户个人资料
- 请求方式: `GET`
- 路径: `/users/:userId/profile`
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "userId": "string",
        "username": "string",
        "nickname": "string",
        "avatar": "string",
        "bio": "string",
        "stats": {
            "videos": 10,
            "followers": 20,
            "following": 30
        },
        "createdAt": "2024-02-26T10:00:00Z"
    }
}
```

### 更新用户个人资料
- 请求方式: `PUT`
- 路径: `/users/:userId/profile`
- 请求头: `Authorization: Bearer {token}`
- Content-Type: `application/json`
- 请求体:
```json
{
    "nickname": "string",  // 可选
    "avatar": "string",    // 可选，头像URL或base64
    "bio": "string"        // 可选，个人简介
}
```
- 响应示例:
```json
{
    "code": 0,
    "msg": "更新成功",
    "data": null
}
```
- 错误情况:
  - 400: 参数不合法
  - 403: 无权操作
  - 404: 用户不存在

### 获取观看历史
- 请求方式: `GET`
- 路径: `/users/:userId/watch-history`
- 请求头: `Authorization: Bearer {token}`
- 查询参数:
  - `page`: 页码，默认1
  - `pageSize`: 每页数量，默认20
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "total": 50,
        "page": 1,
        "pageSize": 20,
        "items": [{
            "id": "string",
            "videoId": "string",
            "title": "string",
            "coverUrl": "string",
            "duration": 180.5,
            "watchedDuration": 120.2,
            "progress": 0.67,
            "watchedAt": "2024-03-16T10:00:00Z"
        }]
    }
}
```
- 错误情况:
  - 403: 无权访问

### 添加到收藏
- 请求方式: `POST`
- 路径: `/videos/:videoId/favorite`
- 请求头: `Authorization: Bearer {token}`
- 响应示例:
```json
{
    "code": 0,
    "msg": "收藏成功",
    "data": null
}
```
- 错误情况:
  - 403: 无权操作
  - 404: 视频不存在
  - 409: 已经收藏过

### 移除收藏
- 请求方式: `DELETE`
- 路径: `/videos/:videoId/favorite`
- 请求头: `Authorization: Bearer {token}`
- 响应示例:
```json
{
    "code": 0,
    "msg": "取消收藏成功",
    "data": null
}
```
- 错误情况:
  - 403: 无权操作
  - 404: 视频不存在或未收藏

### 获取收藏列表
- 请求方式: `GET`
- 路径: `/users/:userId/favorites`
- 请求头: `Authorization: Bearer {token}`
- 查询参数:
  - `page`: 页码，默认1
  - `pageSize`: 每页数量，默认20
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "total": 30,
        "page": 1,
        "pageSize": 20,
        "items": [{
            "id": "string",
            "videoId": "string",
            "title": "string",
            "coverUrl": "string",
            "duration": 180.5,
            "createdAt": "2024-03-16T10:00:00Z",
            "status": "public"
        }]
    }
}
```
- 错误情况:
  - 403: 无权访问

### 记录观看历史
- 请求方式: `POST`
- 路径: `/videos/:videoId/watch`
- 请求头: `Authorization: Bearer {token}`
- Content-Type: `application/json`
- 请求体:
```json
{
    "duration": 120.5,  // 已观看时长（秒）
    "progress": 0.65    // 进度百分比（0-1之间的小数）
}
```
- 响应示例:
```json
{
    "code": 0,
    "msg": "记录成功",
    "data": null
}
```
- 错误情况:
  - 400: 参数不合法
  - 403: 无权操作
  - 404: 视频不存在

## 视频相关

### 获取公开视频列表
- 请求方式: `GET`
- 路径: `/videos/public`
- 查询参数:
  - `page`: 页码，默认1
  - `pageSize`: 每页数量，默认10，最大50
  - `keyword`: 搜索关键词（可选）
  - `sortBy`: 排序字段（可选，支持：created_at/views/likes）
  - `sortOrder`: 排序方向（可选，asc/desc，默认desc）
  - `tags`: 标签筛选（可选，多个标签用逗号分隔）
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
            "thumbnailUrl": "string",
            "duration": 180.5,
            "status": "public",
            "tags": ["标签1", "标签2"],
            "stats": {
                "views": 1000,
                "likes": 100,
                "comments": 50
            },
            "createdAt": "2024-02-26T10:00:00Z"
        }]
    }
}
```

### 获取视频列表（需要认证）
- 请求方式: `GET`
- 路径: `/videos`
- 请求头: `Authorization: Bearer {token}`
- 查询参数:
  - `page`: 页码，默认1
  - `pageSize`: 每页数量，默认10，最大50
  - `status`: 视频状态（可选，支持：public/private/draft）
  - `keyword`: 搜索关键词（可选）
  - `sortBy`: 排序字段（可选，支持：created_at/views/likes/file_size）
  - `sortOrder`: 排序方向（可选，asc/desc，默认desc）
  - `startDate`: 开始日期（可选，格式：YYYY-MM-DD）
  - `endDate`: 结束日期（可选，格式：YYYY-MM-DD）
  - `tags`: 标签筛选（可选，多个标签用逗号分隔）
- 响应格式同上

### 上传视频
- 请求方式: `POST`
- 路径: `/videos`
- 请求头: `Authorization: Bearer {token}`
- Content-Type: `multipart/form-data`
- 请求参数:
  - `file`: 视频文件（必填，支持mp4/mov/avi/wmv/flv/mkv，最大1GB）
  - `cover`: 封面图（可选，支持jpg/jpeg/png，最大2MB）
  - `title`: 标题（必填，1-100个字符）
  - `description`: 描述（可选，最多500个字符）
  - `status`: 状态（可选，public/private/draft，默认private）
  - `duration`: 时长（必填，单位：秒，支持小数点后1位）
  - `tags`: 标签（可选，多个标签用逗号分隔，每个标签最多20个字符）
- 响应示例:
```json
{
    "code": 0,
    "msg": "上传成功",
    "data": {
        "id": "string",
        "title": "string",
        "description": "string",
        "fileName": "string",
        "fileSize": 1024000,
        "format": "mp4",
        "duration": 180.5,
        "status": "private",
        "coverUrl": "string",
        "thumbnailUrl": "string",
        "tags": ["标签1", "标签2"],
        "createdAt": "2024-02-26T10:00:00Z",
        "updatedAt": "2024-02-26T10:00:00Z"
    }
}
```
- 错误情况:
  - 400: 参数不合法
  - 413: 文件大小超过限制
  - 415: 不支持的文件格式
  - 507: 存储空间不足

### 获取视频详情
- 请求方式: `GET`
- 路径: `/videos/:id`
- 请求头: `Authorization: Bearer {token}` (对于非公开视频必须提供)
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "video": {
            "id": "string",
            "userId": "string",
            "title": "string",
            "description": "string",
            "fileName": "string",
            "fileSize": 1024000,
            "format": "mp4",
            "duration": 180.5,
            "status": "public",
            "coverUrl": "string",
            "thumbnailUrl": "string",
            "tags": ["标签1", "标签2"],
            "stats": {
                "views": 1000,
                "likes": 100,
                "comments": 50,
                "shares": 20
            },
            "createdAt": "2024-02-26T10:00:00Z",
            "updatedAt": "2024-02-26T10:00:00Z"
        },
        "isFavorite": true  // 当用户已登录时，返回该用户是否已收藏此视频
    }
}
```
- 错误情况:
  - 404: 视频不存在
  - 403: 无权访问

### 更新视频信息
- 请求方式: `PUT`
- 路径: `/videos/:id`
- 请求头: `Authorization: Bearer {token}`
- Content-Type: `application/json`
- 请求体:
```json
{
    "title": "string",         // 可选
    "description": "string",   // 可选
    "status": "string",        // 可选，public/private/draft
    "tags": ["string"]         // 可选
}
```
- 响应示例:
```json
{
    "code": 0,
    "msg": "更新成功",
    "data": null
}
```
- 错误情况:
  - 400: 参数不合法
  - 403: 无权操作
  - 404: 视频不存在

### 删除视频
- 请求方式: `DELETE`
- 路径: `/videos/:id`
- 请求头: `Authorization: Bearer {token}`
- 响应示例:
```json
{
    "code": 0,
    "msg": "删除成功",
    "data": null
}
```
- 错误情况:
  - 403: 无权操作
  - 404: 视频不存在

### 批量操作视频
- 请求方式: `POST`
- 路径: `/videos/batch`
- 请求头: `Authorization: Bearer {token}`
- Content-Type: `application/json`
- 请求体:
```json
{
    "ids": ["string"],                        // 视频ID列表
    "action": "string",                       // 操作类型：delete/update_status
    "status": "string"                        // 当action为update_status时需要
}
```
- 响应示例:
```json
{
    "code": 0,
    "msg": "操作成功",
    "data": {
        "successCount": 2,
        "failedCount": 0,
        "failedIds": []
    }
}
```

### 更新视频缩略图
- 请求方式: `POST`
- 路径: `/videos/:id/thumbnail`
- 请求头: `Authorization: Bearer {token}`
- Content-Type: `multipart/form-data`
- 请求参数:
  - `file`: 图片文件（必填，支持jpg/jpeg/png，最大2MB）
- 响应示例:
```json
{
    "code": 0,
    "msg": "更新成功",
    "data": {
        "thumbnailUrl": "string"
    }
}
```

### 获取视频统计信息
- 请求方式: `GET`
- 路径: `/videos/:id/stats`
- 请求头: `Authorization: Bearer {token}` (对于非公开视频必须提供)
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "views": 1000,
        "likes": 100,
        "comments": 50,
        "shares": 20
    }
}
```

### 视频流式播放
- 请求方式: `GET`
- 路径: `/videos/:id/stream`
- 请求头: 
  - `Authorization: Bearer {token}` (对于非公开视频必须提供)
  - `Range: bytes=start-end` (可选，支持范围请求)
- 响应头:
  - `Content-Type: video/mp4` (或其他视频格式)
  - `Accept-Ranges: bytes`
  - `Content-Length: size`
  - `Content-Range: bytes start-end/total` (范围请求时)
- 响应状态码:
  - 200: 完整内容
  - 206: 部分内容（范围请求）
  - 403: 无权访问
  - 404: 视频不存在
- 响应内容: 直接返回视频流

## 标记与笔记相关

### 添加标记
- 请求方式: `POST`
- 路径: `/marks`
- 请求头: `Authorization: Bearer {token}`
- Content-Type: `application/json`
- 请求体:
```json
{
    "videoId": "string",      // 视频ID
    "timestamp": 123.45,      // 时间戳，单位：秒
    "content": "string"       // 标记内容
}
```
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "id": "string",
        "userId": "string",
        "videoId": "string",
        "timestamp": 123.45,
        "content": "string",
        "annotations": [],
        "createdAt": "2024-03-16T10:00:00Z",
        "updatedAt": "2024-03-16T10:00:00Z"
    }
}
```
- 错误情况:
  - 400: 参数不合法
  - 403: 无权操作
  - 404: 视频不存在

### 获取标记列表
- 请求方式: `GET`
- 路径: `/marks`
- 请求头: `Authorization: Bearer {token}`
- 查询参数:
  - `videoId`: 视频ID(必填)
  - `page`: 页码，默认1
  - `pageSize`: 每页数量，默认20
  - `sortBy`: 排序字段（可选，默认timestamp）
  - `sortOrder`: 排序方向（可选，asc/desc，默认asc）
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "total": 10,
        "page": 1,
        "pageSize": 20,
        "items": [{
            "id": "string",
            "userId": "string",
            "videoId": "string",
            "timestamp": 123.45,
            "content": "string",
            "annotations": [{
                "id": "string",
                "content": "string",
                "createdAt": "2024-03-16T10:00:00Z"
            }],
            "createdAt": "2024-03-16T10:00:00Z",
            "updatedAt": "2024-03-16T10:00:00Z"
        }]
    }
}
```

### 获取标记详情
- 请求方式: `GET`
- 路径: `/marks/:markId`
- 请求头: `Authorization: Bearer {token}`
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "id": "string",
        "userId": "string",
        "videoId": "string",
        "timestamp": 123.45,
        "content": "string",
        "annotations": [{
            "id": "string",
            "content": "string",
            "createdAt": "2024-03-16T10:00:00Z"
        }],
        "createdAt": "2024-03-16T10:00:00Z",
        "updatedAt": "2024-03-16T10:00:00Z"
    }
}
```
- 错误情况:
  - 403: 无权访问
  - 404: 标记不存在

### 更新标记
- 请求方式: `PUT`
- 路径: `/marks/:markId`
- 请求头: `Authorization: Bearer {token}`
- Content-Type: `application/json`
- 请求体:
```json
{
    "timestamp": 123.45,      // 可选
    "content": "string"       // 可选
}
```
- 响应示例:
```json
{
    "code": 0,
    "msg": "更新成功",
    "data": null
}
```
- 错误情况:
  - 400: 参数不合法
  - 403: 无权操作
  - 404: 标记不存在

### 删除标记
- 请求方式: `DELETE`
- 路径: `/marks/:markId`
- 请求头: `Authorization: Bearer {token}`
- 响应示例:
```json
{
    "code": 0,
    "msg": "删除成功",
    "data": null
}
```
- 错误情况:
  - 403: 无权操作
  - 404: 标记不存在

### 添加注释
- 请求方式: `POST`
- 路径: `/marks/:markId/annotations`
- 请求头: `Authorization: Bearer {token}`
- Content-Type: `application/json`
- 请求体:
```json
{
    "content": "string"       // 注释内容
}
```
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "id": "string",
        "markId": "string",
        "content": "string",
        "createdAt": "2024-03-16T10:00:00Z",
        "updatedAt": "2024-03-16T10:00:00Z"
    }
}
```

### 获取注释列表
- 请求方式: `GET`
- 路径: `/marks/:markId/annotations`
- 请求头: `Authorization: Bearer {token}`
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "total": 5,
        "items": [{
            "id": "string",
            "markId": "string",
            "content": "string",
            "createdAt": "2024-03-16T10:00:00Z",
            "updatedAt": "2024-03-16T10:00:00Z"
        }]
    }
}
```

### 更新注释
- 请求方式: `PUT`
- 路径: `/marks/annotations/:annotationId`
- 请求头: `Authorization: Bearer {token}`
- Content-Type: `application/json`
- 请求体:
```json
{
    "content": "string"
}
```
- 响应示例:
```json
{
    "code": 0,
    "msg": "更新成功",
    "data": null
}
```

### 删除注释
- 请求方式: `DELETE`
- 路径: `/marks/annotations/:annotationId`
- 请求头: `Authorization: Bearer {token}`
- 响应示例:
```json
{
    "code": 0,
    "msg": "删除成功",
    "data": null
}
```

### 添加笔记
- 请求方式: `POST`
- 路径: `/notes`
- 请求头: `Authorization: Bearer {token}`
- Content-Type: `application/json`
- 请求体:
```json
{
    "videoId": "string",      // 视频ID
    "timestamp": 123.45,      // 时间戳，单位：秒，可选
    "content": "string",      // 笔记内容
    "title": "string"         // a笔记标题，可选
}
```
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "id": "string",
        "userId": "string",
        "videoId": "string",
        "timestamp": 123.45,
        "title": "string",
        "content": "string",
        "createdAt": "2024-03-16T10:00:00Z",
        "updatedAt": "2024-03-16T10:00:00Z"
    }
}
```

### 获取笔记列表
- 请求方式: `GET`
- 路径: `/notes`
- 请求头: `Authorization: Bearer {token}`
- 查询参数:
  - `videoId`: 视频ID(必填)
  - `page`: 页码，默认1
  - `pageSize`: 每页数量，默认20
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "total": 10,
        "page": 1,
        "pageSize": 20,
        "items": [{
            "id": "string",
            "userId": "string",
            "videoId": "string",
            "timestamp": 123.45,
            "title": "string",
            "content": "string",
            "createdAt": "2024-03-16T10:00:00Z",
            "updatedAt": "2024-03-16T10:00:00Z"
        }]
    }
}
```

### 更新笔记
- 请求方式: `PUT`
- 路径: `/notes/:noteId`
- 请求头: `Authorization: Bearer {token}`
- Content-Type: `application/json`
- 请求体:
```json
{
    "timestamp": 123.45,    // 可选
    "title": "string",      // 可选
    "content": "string"     // 可选
}
```
- 响应示例:
```json
{
    "code": 0,
    "msg": "更新成功",
    "data": null
}
```

### 删除笔记
- 请求方式: `DELETE`
- 路径: `/notes/:noteId`
- 请求头: `Authorization: Bearer {token}`
- 响应示例:
```json
{
    "code": 0,
    "msg": "删除成功",
    "data": null
}
```

### 导出标记和笔记
- 请求方式: `GET`
- 路径: `/videos/:videoId/export`
- 请求头: `Authorization: Bearer {token}`
- 查询参数:
  - `type`: 导出类型（marks/notes/all，默认all）
  - `format`: 导出格式（txt/json/pdf，默认txt）
- 响应状态码:
  - 200: 成功
  - 403: 无权操作
  - 404: 视频不存在
- 响应头:
  - `Content-Type: application/json` 或 `text/plain` 或 `application/pdf`
  - `Content-Disposition: attachment; filename="export_{videoId}_{timestamp}.{format}"`
- 响应内容: 导出的文件内容 