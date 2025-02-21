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
- 方法: `POST`
- 路径: `/videos`
- Content-Type: `multipart/form-data`
- 参数:
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
#### 请求
- 方法: `GET`
- 路径: `/videos`
- 参数:
  - `page`: 页码（默认1）
  - `pageSize`: 每页数量（默认10，最大50）
  - `keyword`: 关键词搜索，匹配标题和描述
  - `status`: 视频状态筛选
    - 不传：显示所有视频
    - `public`: 只显示公开视频
    - `private`: 只显示私有视频
    - `draft`: 只显示草稿
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
#### 请求
- 方法: `GET`
- 路径: `/videos/:id`

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
#### 请求
- 方法: `PUT`
- 路径: `/videos/:id`
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
#### 请求
- 方法: `DELETE`
- 路径: `/videos/:id`

#### 响应
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

### 6. 批量操作视频
#### 请求
- 方法: `POST`
- 路径: `/videos/batch`
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
#### 请求
- 方法: `POST`
- 路径: `/videos/:id/thumbnail`
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
#### 请求
- 方法: `GET`
- 路径: `/videos/:id/stats`

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
#### 请求
- 方法: `GET`
- 路径: `/videos/:id/stream`
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