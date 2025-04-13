# 测试问题排查指南

## 常见错误与解决方案

### 1. 空指针异常 (nil pointer dereference)

**错误信息示例**：
```
panic: runtime error: invalid memory address or nil pointer dereference [recovered]
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0xf8 pc=0x81aa72]
```

**常见原因**：
- Gin上下文中的Request对象为nil
- 模拟对象未被正确初始化
- 访问未初始化的map或slice
- 类型断言错误

**解决方案**：
1. **Gin测试中初始化Request**：
   ```go
   // 在setupUserTest()中添加
   c.Request = httptest.NewRequest("GET", "/", nil)
   ```

2. **确保模拟对象初始化**：
   ```go
   mockService := new(MockUserService) // 不要忘记new
   ```

3. **初始化map或slice**：
   ```go
   // 错误示例
   var users map[string]string
   users["key"] = "value" // 这会导致nil指针异常
   
   // 正确示例
   users := make(map[string]string)
   users["key"] = "value"
   ```

4. **使用安全的类型断言**：
   ```go
   // 不安全的断言
   result := args.Get(0).(*model.User)
   
   // 安全的断言
   if result, ok := args.Get(0).(*model.User); ok {
       // 使用result
   }
   ```

### 2. 模拟调用验证失败

**错误信息示例**：
```
Error Trace:    user_test.go:123
Error:          mock: The method "GetUserProfile" was not called as expected:
                Expected: 1
                Actual: 0
```

**常见原因**：
- 预期的方法没有被调用
- 调用参数不匹配
- 调用次数不匹配

**解决方案**：
1. **检查方法名称**：确保`On()`方法中的方法名与实际调用的方法名完全一致

2. **使用匹配器放松参数匹配**：
   ```go
   // 严格匹配特定用户ID
   mockService.On("GetUserProfile", mock.Anything, "user123").Return(profileResp, nil)
   
   // 匹配任何字符串用户ID
   mockService.On("GetUserProfile", mock.Anything, mock.AnythingOfType("string")).Return(profileResp, nil)
   ```

3. **打印调试信息**：
   ```go
   fmt.Printf("调用GetUserProfile，参数: %v\n", userID)
   ```

4. **指定调用次数**：
   ```go
   // 预期刚好调用一次
   mockService.On("GetUserProfile", mock.Anything, mock.Anything).Return(profileResp, nil).Once()
   
   // 预期至少调用一次
   mockService.On("GetUserProfile", mock.Anything, mock.Anything).Return(profileResp, nil).AtLeast(1)
   ```

### 3. JSON解析错误

**错误信息示例**：
```
json: cannot unmarshal { into Go value of type response.Response
```

**常见原因**：
- JSON格式错误
- 响应结构与预期不匹配
- Content-Type头设置错误

**解决方案**：
1. **打印原始响应**：
   ```go
   fmt.Println("响应体:", w.Body.String())
   ```

2. **检查响应结构**：
   ```go
   // 正确的响应结构应该是
   var resp response.Response
   err := json.Unmarshal(w.Body.Bytes(), &resp)
   ```

3. **确保Content-Type正确**：
   ```go
   req.Header.Set("Content-Type", "application/json")
   c.Request = req
   ```

### 4. 模拟对象的返回值错误

**错误信息示例**：
```
panic: interface conversion: interface {} is nil, not *model.UserProfileResponse
```

**常见原因**：
- 模拟对象的Return方法返回nil
- 类型断言错误

**解决方案**：
1. **返回类型化的nil**：
   ```go
   // 错误示例
   mockService.On("GetUserProfile", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))
   
   // 正确示例
   mockService.On("GetUserProfile", mock.Anything, mock.Anything).Return((*model.UserProfileResponse)(nil), errors.New("not found"))
   ```

2. **正确处理错误情况**：
   ```go
   // 在模拟中
   if userIDStr == "invalid" {
       return nil, errors.New("user not found")
   }
   ```

### 5. MongoDB集合模拟错误

**错误信息示例**：
```
mock: This mock method call was unexpected:
FindOne(context.Background map[_id:ObjectID("..."])
```

**常见原因**：
- 预期调用与实际调用不匹配
- BSON/ObjectID转换问题

**解决方案**：
1. **使用更灵活的匹配器**：
   ```go
   // 不要过于严格地匹配filter
   mockCollection.On("FindOne", mock.Anything, mock.Anything).Return(...)
   ```

2. **正确处理ObjectID**：
   ```go
   // 字符串转ObjectID
   id, err := primitive.ObjectIDFromHex(userIDStr)
   if err != nil {
       return nil, err
   }
   
   // 在模拟中使用相同的ID
   objectID := primitive.NewObjectID()
   userIDStr := objectID.Hex()
   ```

### 6. 测试覆盖率低

**错误信息示例**：
```
coverage: 47.8% of statements
```

**常见原因**：
- 测试用例不足
- 逻辑分支覆盖不全面
- 错误处理路径未测试

**解决方案**：
1. **查看未覆盖的代码**：
   ```bash
   go test -coverprofile=coverage.out ./...
   go tool cover -html=coverage.out
   ```

2. **添加边界条件测试**：
   ```go
   // 测试正常情况
   TestFunctionSuccess(t *testing.T) { ... }
   
   // 测试参数无效情况
   TestFunctionInvalidInput(t *testing.T) { ... }
   
   // 测试资源不存在情况
   TestFunctionNotFound(t *testing.T) { ... }
   ```

3. **测试错误处理路径**：
   ```go
   // 模拟数据库错误
   mockCollection.On("FindOne", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
   ```

## 特定问题示例与解决

### 问题1：`c.Request.Context()` 导致空指针异常

**错误信息**：
```
panic: runtime error: invalid memory address or nil pointer dereference
...
net/http.(*Request).Context(...)
        /go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.6.linux-amd64/src/net/http/request.go:352
video-platform/internal/handler.(*UserHandler).GetUserProfile(0xc0001f6550, 0xc0001fa100)
```

**解决方案**：
```go
// 在setupTest函数中
func setupUserTest() (*gin.Context, *httptest.ResponseRecorder, *MockUserService, *UserHandler) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    
    // 添加这一行代码
    c.Request = httptest.NewRequest("GET", "/", nil)
    
    mockService := new(MockUserService)
    handler := NewUserHandler(mockService)
    return c, w, mockService, handler
}
```

### 问题2：模拟`UpdateOne`方法返回更新结果

**错误信息**：
```
panic: interface conversion: interface {} is nil, not *mongo.UpdateResult
```

**解决方案**：
```go
// 创建一个UpdateResult对象
updateResult := &mongo.UpdateResult{
    MatchedCount:  1,
    ModifiedCount: 1,
}

// 在模拟中正确返回
mockCollection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
    Return(updateResult, nil)
```

### 问题3：测试时需要模拟文件上传

**问题描述**：需要测试用户头像上传功能，但不知道如何在测试中模拟文件上传

**解决方案**：
```go
// 创建一个临时文件
tempFile, err := ioutil.TempFile("", "test-*.jpg")
if err != nil {
    t.Fatal(err)
}
defer os.Remove(tempFile.Name())
defer tempFile.Close()

// 写入一些测试数据
if _, err := tempFile.Write([]byte("fake image data")); err != nil {
    t.Fatal(err)
}
tempFile.Close()

// 创建multipart文件
file, err := os.Open(tempFile.Name())
if err != nil {
    t.Fatal(err)
}
defer file.Close()

// 创建multipart请求
body := &bytes.Buffer{}
writer := multipart.NewWriter(body)
part, err := writer.CreateFormFile("avatar", "test.jpg")
if err != nil {
    t.Fatal(err)
}
io.Copy(part, file)
writer.Close()

// 创建请求
req := httptest.NewRequest("POST", "/", body)
req.Header.Set("Content-Type", writer.FormDataContentType())
c.Request = req

// 在handler中访问上传的文件
c.Request.ParseMultipartForm(10 << 20) // 10 MB
```

## 测试性能优化

如果测试运行速度较慢，可以考虑以下优化：

1. **并行测试**：
   ```go
   func TestFunction(t *testing.T) {
       t.Parallel() // 允许并行运行测试
       // ...
   }
   ```

2. **缓存测试结果**：
   ```bash
   # 正常运行测试（支持缓存）
   go test ./...
   
   # 强制重新运行所有测试（不使用缓存）
   go test -count=1 ./...
   ```

3. **跳过慢测试**：
   ```go
   func TestSlowFunction(t *testing.T) {
       if testing.Short() {
           t.Skip("跳过慢速测试")
       }
       // ...测试逻辑
   }
   ```
   使用`go test -short ./...`运行时会跳过这些测试

4. **测试超时控制**：
   ```go
   func TestFunction(t *testing.T) {
       ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
       defer cancel()
       
       // 使用ctx调用被测试函数
   }
   ```

## 测试最佳实践总结

1. **测试文件结构清晰**：
   - 每个测试函数专注于一个功能点
   - 测试名称清晰表达测试内容
   - 使用表格驱动测试测试多个场景

2. **模拟外部依赖**：
   - 所有外部依赖都应该被模拟
   - 模拟响应应该与真实情况一致
   - 模拟错误应该与真实错误场景对应

3. **测试断言明确**：
   - 每个测试函数都应该有明确的断言
   - 使用`assert`或`require`包简化断言
   - 错误信息应该清晰描述预期行为

4. **清理测试资源**：
   - 使用`defer`清理资源
   - 测试前后保持数据库一致
   - 删除测试期间创建的临时文件

5. **持续改进测试**：
   - 定期检查测试覆盖率
   - 新功能必须有对应测试
   - 修复bug时添加回归测试 