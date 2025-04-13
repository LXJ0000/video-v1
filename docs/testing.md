# 视频平台后端测试指南

## 目录
- [测试架构](#测试架构)
- [测试编写指南](#测试编写指南)
  - [服务层测试](#服务层测试)
  - [处理器层测试](#处理器层测试)
- [模拟技术](#模拟技术)
  - [MongoDB模拟](#mongodb模拟)
  - [服务层模拟](#服务层模拟)
- [常见问题解决方案](#常见问题解决方案)
- [测试覆盖率](#测试覆盖率)
- [持续集成](#持续集成)

## 测试架构

本项目采用分层测试策略，主要测试以下几个层面：

```
+------------------+
|     端到端测试     |  <- 测试整个API流程
+------------------+
|    处理器测试      |  <- 测试HTTP请求处理
+------------------+
|    服务层测试      |  <- 测试业务逻辑
+------------------+
|    单元测试       |  <- 测试独立函数
+------------------+
```

### 测试工具与库

- **标准库**：`testing` - Go标准测试框架
- **断言库**：`github.com/stretchr/testify/assert` - 提供丰富的断言功能
- **模拟库**：`github.com/stretchr/testify/mock` - 支持接口模拟
- **HTTP测试**：`net/http/httptest` - HTTP请求模拟
- **Gin测试**：`github.com/gin-gonic/gin` - Gin框架提供的测试工具

## 测试编写指南

### 服务层测试

服务层测试主要测试业务逻辑，需要模拟数据库操作。

#### 测试文件命名

服务层测试文件应命名为 `{服务名}_test.go`，放在对应服务的同一目录下。

示例：
- 服务文件：`internal/service/user.go`
- 测试文件：`internal/service/user_test.go`

#### 测试函数命名

测试函数命名应遵循 `Test{函数名}{场景描述}` 的格式，清晰表达测试的功能和场景。

示例：
```go
// 正常场景
func TestGetUserProfileSuccess(t *testing.T) {
    // 测试正常获取用户资料的情况
}

// 异常场景
func TestGetUserProfileNotFound(t *testing.T) {
    // 测试用户不存在的情况
}
```

#### 服务层测试模板

```go
func TestServiceFunction(t *testing.T) {
    // 1. 准备测试数据
    // 创建测试输入和预期输出
    
    // 2. 模拟外部依赖
    // 如数据库操作、第三方服务等
    
    // 3. 执行被测试函数
    // 调用服务函数并获取结果
    
    // 4. 验证结果
    // 使用assert断言结果符合预期
    
    // 5. 验证模拟调用
    // 确保模拟对象的方法被正确调用
}
```

### 处理器层测试

处理器层测试主要测试HTTP请求处理逻辑，需要模拟服务层和HTTP请求。

#### 测试文件命名

处理器层测试文件应命名为 `{处理器名}_test.go`，放在对应处理器的同一目录下。

示例：
- 处理器文件：`internal/handler/user.go`
- 测试文件：`internal/handler/user_test.go`

#### 处理器测试模板

```go
func TestHandlerFunction(t *testing.T) {
    // 1. 设置测试环境
    c, w, mockService, handler := setupTest()
    
    // 2. 模拟请求参数
    // 设置请求参数、路径参数、查询参数等
    
    // 3. 模拟服务层响应
    // 设置mockService的预期行为和返回值
    
    // 4. 执行处理器函数
    handler.HandleFunction(c)
    
    // 5. 验证HTTP响应
    // 检查状态码、响应体等
    
    // 6. 验证服务层调用
    // 确保服务层方法被正确调用
}
```

#### 重要！处理器测试注意事项

1. **请求对象初始化**：
   ```go
   // 确保请求对象不为nil，否则c.Request.Context()会导致空指针异常
   c.Request = httptest.NewRequest("GET", "/", nil)
   ```

2. **模拟JSON请求体**：
   ```go
   jsonStr := `{"username": "testuser", "password": "password123"}`
   req := httptest.NewRequest("POST", "/", strings.NewReader(jsonStr))
   req.Header.Set("Content-Type", "application/json")
   c.Request = req
   ```

3. **模拟表单请求**：
   ```go
   form := url.Values{}
   form.Add("username", "testuser")
   form.Add("password", "password123")
   req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   c.Request = req
   ```

4. **模拟认证用户**：
   ```go
   // 模拟已登录用户
   c.Set("userId", "user123")
   c.Set("username", "testuser")
   ```

5. **模拟路径参数**：
   ```go
   c.Params = []gin.Param{
       {Key: "userId", Value: "user123"},
   }
   ```

6. **模拟查询参数**：
   ```go
   // 方法一：直接设置Request
   c.Request = httptest.NewRequest("GET", "/?page=1&size=10", nil)
   
   // 方法二：使用Query方法
   c.Request = httptest.NewRequest("GET", "/", nil)
   q := c.Request.URL.Query()
   q.Add("page", "1")
   q.Add("size", "10")
   c.Request.URL.RawQuery = q.Encode()
   ```

## 模拟技术

### MongoDB模拟

MongoDB模拟使用自定义的`MockCollection`结构体，实现与MongoDB集合相同的接口方法。

```go
// 模拟Collection接口
type MockCollection struct {
    mock.Mock
    *mongo.Collection
}

// 模拟FindOne方法
func (m *MockCollection) FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult {
    args := m.Called(ctx, filter)
    return args.Get(0).(*mongo.SingleResult)
}

// 模拟数据库响应
func mockFindOneResponse(doc interface{}) *mongo.SingleResult {
    return mongo.NewSingleResultFromDocument(doc, nil, nil)
}

// 模拟空结果
func mockEmptyResult() *mongo.SingleResult {
    return mongo.NewSingleResultFromDocument(nil, mongo.ErrNoDocuments, nil)
}
```

### 服务层模拟

服务层模拟使用`testify/mock`库，实现与实际服务相同的接口。

```go
// 创建UserService的Mock
type MockUserService struct {
    mock.Mock
}

// 实现GetUserProfile方法
func (m *MockUserService) GetUserProfile(ctx context.Context, id string) (*model.UserProfileResponse, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*model.UserProfileResponse), args.Error(1)
}

// 使用方法
mockService := new(MockUserService)
profileResp := &model.UserProfileResponse{...}
mockService.On("GetUserProfile", mock.Anything, "user123").Return(profileResp, nil)

// 验证调用
mockService.AssertExpectations(t)
```

## 常见问题解决方案

### 1. 空指针异常：`invalid memory address or nil pointer dereference`

**问题原因**：通常是因为在测试中访问了未初始化的对象

**解决方案**：
- 检查Request对象是否为nil：`c.Request = httptest.NewRequest("GET", "/", nil)`
- 检查模拟对象是否正确初始化：`mockSvc := new(MockService)`
- 检查返回值断言是否正确：`return args.Get(0).(*model.User), args.Error(1)`

### 2. 模拟方法未被调用：`mock: The method "GetUserProfile" was not called as expected`

**问题原因**：代码中没有调用模拟对象的方法，或者调用参数不匹配

**解决方案**：
- 确保参数匹配：使用`mock.Anything`或准确的参数值
- 检查函数调用路径：使用调试输出确认函数是否被调用
- 检查`On()`方法的参数：应与实际调用参数完全匹配

```go
// 更灵活的参数匹配
mockService.On("GetUserProfile", mock.Anything, mock.AnythingOfType("string")).Return(profileResp, nil)
```

### 3. JSON解析错误：`invalid character '{' looking for beginning of value`

**问题原因**：JSON响应格式不正确或解析时出错

**解决方案**：
- 使用`w.Body.String()`打印原始响应
- 检查JSON生成是否正确：`json.Marshal`、`c.JSON`等
- 确保Content-Type正确：`"Content-Type": "application/json"`

### 4. 测试覆盖率低：`coverage: 42.3% of statements`

**问题原因**：测试案例不够全面，没有覆盖全部代码路径

**解决方案**：
- 使用`go test -cover`查看哪些代码未被覆盖
- 添加更多测试用例，包括正常和异常场景
- 使用`go test -coverprofile=coverage.out ./...`生成详细报告

## 测试覆盖率

项目要求服务层和处理器层的测试覆盖率达到以下标准：

- **核心业务逻辑**：>= 80%
- **错误处理路径**：>= 70%
- **边界条件**：>= 60%

### 生成覆盖率报告

```bash
# 生成覆盖率数据
go test -coverprofile=coverage.out ./...

# 在浏览器中查看HTML格式报告
go tool cover -html=coverage.out

# 查看命令行摘要
go tool cover -func=coverage.out
```

### 忽略不需要测试的代码

对于不需要测试的代码，可以添加注释：

```go
// 忽略测试覆盖率
//go:build !test

// 或者使用注释
// test:ignore-next-line
func someUtilityFunction() {
    // ...
}
```

## 持续集成

项目使用GitHub Actions进行持续集成，每次提交代码都会自动运行测试。

### CI流程

1. 检出代码
2. 设置Go环境
3. 安装依赖
4. 运行测试
5. 生成覆盖率报告
6. 发布测试结果

### CI配置示例

```yaml
name: Go Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.20
    - name: Install dependencies
      run: go mod download
    - name: Run tests
      run: go test -v -coverprofile=coverage.out ./...
    - name: Upload coverage report
      uses: codecov/codecov-action@v1
      with:
        file: ./coverage.out
        fail_ci_if_error: true
``` 