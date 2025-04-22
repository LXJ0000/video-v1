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
- 路径: `/users/send_sms_code`
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

### 获取用户个人资料
- 请求方式: `GET`
- 路径: `/users/:userId/profile`
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "id": "string",
        "username": "string",
        "nickname": "string",
        "email": "string",
        "avatar": "string",
        "bio": "string",
        "stats": {
            "uploadedVideos": 10,
            "totalWatchTime": 120,
            "totalLikes": 50
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
    "username": "string",  // 可选
    "email": "string",     // 可选
    "nickname": "string",  // 可选
    "bio": "string",       // 可选，个人简介
    "avatar": "string"     // 可选，base64编码的头像
}
```
- 响应示例:
```json
{
    "code": 0,
    "msg": "更新成功",
    "data": {
        "id": "string",
        "username": "string",
        "nickname": "string",
        "email": "string",
        "avatar": "string",
        "bio": "string",
        "stats": {
            "uploadedVideos": 10,
            "totalWatchTime": 120,
            "totalLikes": 50
        },
        "createdAt": "2024-02-26T10:00:00Z"
    }
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
  - `size`: 每页数量，默认12
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "history": [
            {
                "id": "string",
                "videoId": "string",
                "videoTitle": "string",
                "coverUrl": "string",
                "watchedAt": "2024-02-26T10:00:00Z",
                "progress": 60,
                "videoDuration": 120
            }
        ],
        "total": 50,
        "page": 1,
        "size": 12
    }
}
```

### 获取收藏列表
- 请求方式: `GET`
- 路径: `/users/:userId/favorites`
- 请求头: `Authorization: Bearer {token}`
- 查询参数:
  - `page`: 页码，默认1
  - `size`: 每页数量，默认12
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "favorites": [
            {
                "id": "string",
                "videoId": "string",
                "videoTitle": "string",
                "coverUrl": "string",
                "addedAt": "2024-02-26T10:00:00Z",
                "videoDuration": 120
            }
        ],
        "total": 30,
        "page": 1,
        "size": 12
    }
}
```

### 添加视频到收藏
- 请求方式: `POST`
- 路径: `/videos/:videoId/favorite`
- 请求头: `Authorization: Bearer {token}`
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "message": "添加收藏成功"
    }
}
```

### 从收藏中移除视频
- 请求方式: `DELETE`
- 路径: `/videos/:videoId/favorite`
- 请求头: `Authorization: Bearer {token}`
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "message": "取消收藏成功"
    }
}
```

### 记录观看历史
- 请求方式: `POST`
- 路径: `/videos/:videoId/watch`
- 请求头: `Authorization: Bearer {token}`
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "message": "记录观看历史成功"
    }
}
```

## 视频相关接口

### 获取公开视频列表
- 请求方式: `GET`
- 路径: `/videos/public`
- 查询参数:
  - `page`: 页码，默认1
  - `size`: 每页数量，默认12
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "videos": [
            {
                "id": "string",
                "title": "string",
                "coverUrl": "string",
                "videoUrl": "string",
                "description": "string",
                "duration": 120,
                "userId": "string",
                "username": "string",
                "stats": {
                    "views": 100,
                    "likes": 10,
                    "comments": 5
                },
                "createdAt": "2024-02-26T10:00:00Z"
            }
        ],
        "total": 100,
        "page": 1,
        "size": 12
    }
}
```

### 获取用户视频列表
- 请求方式: `GET`
- 路径: `/videos`
- 请求头: `Authorization: Bearer {token}`
- 查询参数:
  - `page`: 页码，默认1
  - `size`: 每页数量，默认12
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "videos": [
            {
                "id": "string",
                "title": "string",
                "coverUrl": "string",
                "videoUrl": "string",
                "description": "string",
                "duration": 120,
                "stats": {
                    "views": 100,
                    "likes": 10,
                    "comments": 5
                },
                "status": "public",
                "createdAt": "2024-02-26T10:00:00Z",
                "updatedAt": "2024-02-26T10:00:00Z"
            }
        ],
        "total": 50,
        "page": 1,
        "size": 12
    }
}
```

### 上传视频
- 请求方式: `POST`
- 路径: `/videos`
- 请求头: `Authorization: Bearer {token}`
- Content-Type: `multipart/form-data`
- 表单参数:
  - `title`: 视频标题
  - `description`: 视频描述
  - `video`: 视频文件
  - `cover`: 封面图片
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "id": "string",
        "title": "string",
        "coverUrl": "string",
        "videoUrl": "string",
        "description": "string",
        "duration": 120,
        "userId": "string",
        "username": "string",
        "stats": {
            "views": 0,
            "likes": 0,
            "comments": 0
        },
        "status": "public",
        "createdAt": "2024-02-26T10:00:00Z",
        "updatedAt": "2024-02-26T10:00:00Z"
    }
}
```

### 获取视频详情
- 请求方式: `GET`
- 路径: `/videos/:videoId`
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "id": "string",
        "title": "string",
        "coverUrl": "string",
        "videoUrl": "string",
        "description": "string",
        "duration": 120,
        "userId": "string",
        "username": "string",
        "stats": {
            "views": 100,
            "likes": 10,
            "comments": 5
        },
        "status": "public",
        "favorited": true,
        "createdAt": "2024-02-26T10:00:00Z",
        "updatedAt": "2024-02-26T10:00:00Z"
    }
}
```

### 更新视频信息
- 请求方式: `PUT`
- 路径: `/videos/:videoId`
- 请求头: `Authorization: Bearer {token}`
- Content-Type: `application/json`
- 请求体:
```json
{
    "title": "string",       // 可选
    "description": "string", // 可选
    "status": "string"       // 可选，public或private
}
```
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "id": "string",
        "title": "string",
        "coverUrl": "string",
        "videoUrl": "string",
        "description": "string",
        "duration": 120,
        "userId": "string",
        "username": "string",
        "stats": {
            "views": 100,
            "likes": 10,
            "comments": 5
        },
        "status": "public",
        "createdAt": "2024-02-26T10:00:00Z",
        "updatedAt": "2024-02-26T10:00:00Z"
    }
}
```

### 删除视频
- 请求方式: `DELETE`
- 路径: `/videos/:videoId`
- 请求头: `Authorization: Bearer {token}`
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": null
}
```

### 视频流播放
- 请求方式: `GET`
- 路径: `/videos/:videoId/stream`
- 参数:
  - 可选的 Range 头，支持断点续传
- 响应:
  - Content-Type: video/mp4
  - 支持范围请求(206 Partial Content)

### 更新视频缩略图
- 请求方式: `POST`
- 路径: `/videos/:videoId/thumbnail`
- 请求头: `Authorization: Bearer {token}`
- Content-Type: `multipart/form-data`
- 表单参数:
  - `cover`: 新封面图片
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "coverUrl": "string"
    }
}
```

### 获取视频统计信息
- 请求方式: `GET`
- 路径: `/videos/:videoId/stats`
- 请求头: `Authorization: Bearer {token}`
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "views": 100,
        "likes": 10,
        "comments": 5
    }
}
```

### 批量操作视频
- 请求方式: `POST`
- 路径: `/videos/batch`
- 请求头: `Authorization: Bearer {token}`
- Content-Type: `application/json`
- 请求体:
```json
{
    "operation": "string", // delete, public, private
    "videoIds": ["string"]
}
```
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "success": true,
        "message": "批量操作成功"
    }
}
```

## 标记相关接口

### 添加标记
- 请求方式: `POST`
- 路径: `/marks`
- 请求头: `Authorization: Bearer {token}`
- Content-Type: `application/json`
- 请求体:
```json
{
    "videoId": "string",
    "content": "string",
    "timestamp": 30,
    "type": "comment"
}
```
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "id": "string",
        "videoId": "string",
        "userId": "string",
        "username": "string",
        "content": "string",
        "timestamp": 30,
        "type": "comment",
        "createdAt": "2024-02-26T10:00:00Z",
        "updatedAt": "2024-02-26T10:00:00Z"
    }
}
```

### 获取标记列表
- 请求方式: `GET`
- 路径: `/marks`
- 请求头: `Authorization: Bearer {token}`
- 查询参数:
  - `videoId`: 视频ID
  - `page`: 页码，默认1
  - `size`: 每页数量，默认20
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "marks": [
            {
                "id": "string",
                "videoId": "string",
                "userId": "string",
                "username": "string",
                "content": "string",
                "timestamp": 30,
                "type": "comment",
                "createdAt": "2024-02-26T10:00:00Z",
                "updatedAt": "2024-02-26T10:00:00Z"
            }
        ],
        "total": 20,
        "page": 1,
        "size": 20
    }
}
```

### 更新标记
- 请求方式: `PUT`
- 路径: `/marks/:markId`
- 请求头: `Authorization: Bearer {token}`
- Content-Type: `application/json`
- 请求体:
```json
{
    "content": "string",
    "timestamp": 35
}
```
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "id": "string",
        "videoId": "string",
        "userId": "string",
        "username": "string",
        "content": "string",
        "timestamp": 35,
        "type": "comment",
        "createdAt": "2024-02-26T10:00:00Z",
        "updatedAt": "2024-02-26T10:00:00Z"
    }
}
```

### 删除标记
- 请求方式: `DELETE`
- 路径: `/marks/:markId`
- 请求头: `Authorization: Bearer {token}`
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": null
}
```

### 添加注释
- 请求方式: `POST`
- 路径: `/marks/:markId/annotations`
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
    "msg": "success",
    "data": {
        "id": "string",
        "markId": "string",
        "userId": "string",
        "username": "string",
        "content": "string",
        "createdAt": "2024-02-26T10:00:00Z",
        "updatedAt": "2024-02-26T10:00:00Z"
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
    "data": [
        {
            "id": "string",
            "markId": "string",
            "userId": "string",
            "username": "string",
            "content": "string",
            "createdAt": "2024-02-26T10:00:00Z",
            "updatedAt": "2024-02-26T10:00:00Z"
        }
    ]
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
    "msg": "success",
    "data": {
        "id": "string",
        "markId": "string",
        "userId": "string",
        "username": "string",
        "content": "string",
        "createdAt": "2024-02-26T10:00:00Z",
        "updatedAt": "2024-02-26T10:00:00Z"
    }
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
    "msg": "success",
    "data": null
}
```

## 笔记相关接口

### 添加笔记
- 请求方式: `POST`
- 路径: `/notes`
- 请求头: `Authorization: Bearer {token}`
- Content-Type: `application/json`
- 请求体:
```json
{
    "videoId": "string",
    "title": "string",
    "content": "string"
}
```
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "id": "string",
        "videoId": "string",
        "userId": "string",
        "title": "string",
        "content": "string",
        "createdAt": "2024-02-26T10:00:00Z",
        "updatedAt": "2024-02-26T10:00:00Z"
    }
}
```

### 获取笔记列表
- 请求方式: `GET`
- 路径: `/notes`
- 请求头: `Authorization: Bearer {token}`
- 查询参数:
  - `videoId`: 视频ID
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": [
        {
            "id": "string",
            "videoId": "string",
            "userId": "string",
            "title": "string",
            "content": "string",
            "createdAt": "2024-02-26T10:00:00Z",
            "updatedAt": "2024-02-26T10:00:00Z"
        }
    ]
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
    "title": "string",
    "content": "string"
}
```
- 响应示例:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "id": "string",
        "videoId": "string",
        "userId": "string",
        "title": "string",
        "content": "string",
        "createdAt": "2024-02-26T10:00:00Z",
        "updatedAt": "2024-02-26T10:00:00Z"
    }
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
    "msg": "success",
    "data": null
}
```

## 导出相关接口

### 导出标记、注释和笔记
- 请求方式: `GET`
- 路径: `/videos/export`
- 请求头: `Authorization: Bearer {token}`
- 查询参数:
  - `videoId`: 视频ID
  - `format`: 导出格式，支持 `json`, `csv`, `pdf`，默认为 `json`
- 响应:
  - Content-Type: 取决于format参数
  - Content-Disposition: attachment; filename="export.{format}" 