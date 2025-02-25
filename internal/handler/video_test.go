package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
	"video-platform/config"
	"video-platform/internal/model"
	"video-platform/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	testMark = fmt.Sprintf("test_%d", time.Now().UnixNano())
	testLock sync.Mutex
)

func init() {
	// 设置测试环境
	testDB := fmt.Sprintf("video_platform_test_%d", time.Now().UnixNano())
	config.GlobalConfig.MongoDB.Database = testDB

	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)
}

// setupTestRouter 初始化测试环境并返回清理函数
func setupTestRouter(t *testing.T) (*gin.Engine, func()) {
	// 初始化配置和数据库连接
	if err := config.Init(); err != nil {
		t.Fatal(err)
	}
	if err := database.InitMongoDB(context.Background(), config.GlobalConfig.MongoDB, true); err != nil {
		t.Fatal(err)
	}

	// 清理测试数据
	cleanTestData(t)

	// 设置路由
	r := gin.Default()
	InitRoutes(r)

	// 返回清理函数
	return r, func() {
		testLock.Lock()
		defer testLock.Unlock()
		cleanTestData(t)
	}
}

// cleanTestData 清理测试数据
func cleanTestData(t *testing.T) {
	ctx := context.Background()

	// 删除测试数据库记录
	_, err := database.GetCollection("videos").DeleteMany(ctx, bson.M{
		"user_id": testMark,
	})
	if err != nil {
		t.Fatalf("清理测试数据失败: %v", err)
	}

	// 清理测试上传的文件
	files, err := os.ReadDir(config.GlobalConfig.Storage.UploadDir)
	if err == nil {
		for _, file := range files {
			// 清理所有测试相关文件
			name := file.Name()
			if strings.HasPrefix(name, "test-") ||
				strings.HasPrefix(name, "cover_") ||
				strings.Contains(name, "test") {
				os.Remove(filepath.Join(config.GlobalConfig.Storage.UploadDir, name))
			}
		}
	}
}

// func TestMain(m *testing.M) {
// 	// 初始化配置
// 	if err := config.Init(); err != nil {
// 		fmt.Printf("配置初始化失败: %v\n", err)
// 		os.Exit(1)
// 	}

// 	// 初始化数据库连接
// 	if err := database.InitMongoDB(context.Background(), config.GlobalConfig.MongoDB, true); err != nil {
// 		fmt.Printf("数据库连接失败: %v\n", err)
// 		os.Exit(1)
// 	}

// 	// 运行测试前清理
// 	cleanAllTestData()

// 	// 运行测试
// 	code := m.Run()

// 	// 测试结束后清理
// 	cleanAllTestData()

// 	// 关闭数据库连接
// 	database.CloseMongoDB()

// 	os.Exit(code)
// }

// cleanAllTestData 清理所有测试数据（不需要testing.T）
func cleanAllTestData() {
	// 清理数据库记录
	ctx := context.Background()
	if client := database.GetClient(); client != nil {
		// 使用更宽松的匹配条件
		_, _ = database.GetCollection("videos").DeleteMany(ctx, bson.M{})
	}

	// 清理上传的文件
	if uploadDir := config.GlobalConfig.Storage.UploadDir; uploadDir != "" {
		// 清理 handler/uploads 目录
		cleanDirectory(filepath.Join("internal/handler", uploadDir))
		// 清理 service/uploads 目录
		cleanDirectory(filepath.Join("internal/service", uploadDir))
	}
}

// cleanDirectory 清理指定目录下的所有文件和目录本身
func cleanDirectory(dir string) {
	// 获取项目根目录的绝对路径
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取当前目录失败: %v\n", err)
		return
	}

	// 如果当前在 internal/handler 或 internal/service 目录下，需要回到项目根目录
	if strings.Contains(rootDir, "internal/handler") {
		rootDir = filepath.Dir(filepath.Dir(rootDir))
	} else if strings.Contains(rootDir, "internal/service") {
		rootDir = filepath.Dir(filepath.Dir(rootDir))
	}

	// 构建完整的目录路径
	fullPath := filepath.Join(rootDir, dir)

	// 确保目录存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return
	}

	// 先删除目录中的所有文件
	files, err := os.ReadDir(fullPath)
	if err == nil {
		for _, file := range files {
			filePath := filepath.Join(fullPath, file.Name())
			if err := os.RemoveAll(filePath); err != nil {
				fmt.Printf("删除文件/目录 %s 失败: %v\n", filePath, err)
			}
		}
	}

	// 最后删除目录本身
	if err := os.Remove(fullPath); err != nil {
		fmt.Printf("删除目录 %s 失败: %v\n", fullPath, err)
	}
}

// 修改所有测试用例，使用 defer 调用清理函数
func TestVideoHandler_Upload(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	token := getTestToken(t)

	// 创建测试文件
	tempFile, err := os.CreateTemp("", "test-video-*.mp4")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// 准备请求
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.mp4")
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, tempFile)
	if err != nil {
		t.Fatal(err)
	}
	writer.WriteField("title", "Test Video")
	writer.WriteField("duration", "180.5")
	writer.Close()

	// 发送请求
	req := httptest.NewRequest("POST", "/api/v1/videos", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var response struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data model.Video `json:"data"`
	}
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 0, response.Code)
	assert.Equal(t, 180.5, response.Data.Duration)
}

func TestVideoHandler_GetList(t *testing.T) {
	r, cleanup := setupTestRouter(t)
	defer cleanup()

	token := getTestToken(t)

	// 创建请求
	req := httptest.NewRequest("GET", "/api/v1/videos?page=1&pageSize=10", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// 执行请求
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
}

func TestVideoHandler_GetByID(t *testing.T) {
	r, cleanup := setupTestRouter(t)
	defer cleanup()

	token := getTestToken(t)

	// 先上传一个视频
	video := uploadTestVideo(t, r)

	// 创建请求
	req := httptest.NewRequest("GET", "/api/v1/videos/"+video.ID.Hex(), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// 执行请求
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
}

func TestVideoHandler_Stream(t *testing.T) {
	r, cleanup := setupTestRouter(t)
	defer cleanup()

	token := getTestToken(t)

	// 先上传一个视频
	video := uploadTestVideo(t, r)

	// 创建请求
	req := httptest.NewRequest("GET", "/api/v1/videos/"+video.ID.Hex()+"/stream", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// 执行请求
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "video/")
	assert.NotEmpty(t, w.Header().Get("Content-Length"))
	assert.Equal(t, "bytes", w.Header().Get("Accept-Ranges"))
}

// 辅助函数：创建测试视频文件
func createTestVideoFile(t *testing.T) *os.File {
	// 确保上传目录存在
	if err := os.MkdirAll(config.GlobalConfig.Storage.UploadDir, 0755); err != nil {
		t.Fatalf("创建上传目录失败: %v", err)
	}

	// 在上传目录中创建测试文件
	tempFile, err := os.CreateTemp(config.GlobalConfig.Storage.UploadDir, "test-*.mp4")
	if err != nil {
		t.Fatal(err)
	}

	// 写入一些测试数据
	content := []byte("fake video content")
	if _, err := tempFile.Write(content); err != nil {
		t.Fatal(err)
	}

	if err := tempFile.Sync(); err != nil {
		t.Fatal(err)
	}

	if _, err := tempFile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	return tempFile
}

// 辅助函数：上传测试视频并返回视频信息
func uploadTestVideo(t *testing.T, r *gin.Engine) *model.Video {
	tempFile := createTestVideoFile(t)
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()

	token := getTestToken(t)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// 创建文件表单字段
	part, err := writer.CreateFormFile("file", filepath.Base(tempFile.Name()))
	if err != nil {
		t.Fatal(err)
	}

	// 读取文件内容
	fileContent, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// 写入文件内容
	if _, err := part.Write(fileContent); err != nil {
		t.Fatal(err)
	}

	// 添加其他字段
	writer.WriteField("title", "测试视频")
	writer.WriteField("description", "这是一个测试视频")
	writer.WriteField("userId", TestUserID)
	writer.WriteField("duration", "180.5")
	writer.Close()

	// 创建请求
	req := httptest.NewRequest("POST", "/api/v1/videos", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// 执行请求
	r.ServeHTTP(w, req)

	// 检查响应
	if w.Code != http.StatusOK {
		t.Fatalf("上传视频失败: %s", w.Body.String())
	}

	var response struct {
		Code int         `json:"code"`
		Data model.Video `json:"data"`
	}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	return &response.Data
}

// 添加新的测试用例

func TestVideoHandler_BatchOperation(t *testing.T) {
	r, cleanup := setupTestRouter(t)
	defer cleanup()

	token := getTestToken(t)

	// 先上传两个测试视频
	video1 := uploadTestVideo(t, r)
	video2 := uploadTestVideo(t, r)

	// 创建批量操作请求
	reqBody := model.BatchOperationRequest{
		IDs:    []string{video1.ID.Hex(), video2.ID.Hex()},
		Action: "update_status",
		Status: model.VideoStatusPrivate,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}

	// 创建请求
	req := httptest.NewRequest("POST", "/api/v1/videos/batch", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// 执行请求
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
}

func TestVideoHandler_UpdateThumbnail(t *testing.T) {
	r, cleanup := setupTestRouter(t)
	defer cleanup()

	token := getTestToken(t)

	// 先上传一个测试视频
	video := uploadTestVideo(t, r)

	// 创建测试缩略图
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.jpg")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("fake image content"))
	writer.Close()

	// 创建请求
	req := httptest.NewRequest("POST", "/api/v1/videos/"+video.ID.Hex()+"/thumbnail", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// 执行请求
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
}

func TestVideoHandler_GetStats(t *testing.T) {
	r, cleanup := setupTestRouter(t)
	defer cleanup()

	token := getTestToken(t)

	// 先上传一个测试视频
	video := uploadTestVideo(t, r)

	// 创建请求
	req := httptest.NewRequest("GET", "/api/v1/videos/"+video.ID.Hex()+"/stats", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// 执行请求
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
}

// TestVideoHandler_UploadWithCover 测试上传带封面图的视频
func TestVideoHandler_UploadWithCover(t *testing.T) {
	r, cleanup := setupTestRouter(t)
	defer cleanup()

	token := getTestToken(t)

	// 创建测试文件
	videoContent := []byte("fake video content")
	coverContent := []byte("fake image content")

	// 创建multipart请求
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// 添加视频文件
	videoPart, err := writer.CreateFormFile("file", "test.mp4")
	if err != nil {
		t.Fatal(err)
	}
	videoPart.Write(videoContent)

	// 添加封面图文件
	coverPart, err := writer.CreateFormFile("cover", "cover.jpg")
	if err != nil {
		t.Fatal(err)
	}
	coverPart.Write(coverContent)

	// 添加其他字段
	writer.WriteField("title", "测试视频")
	writer.WriteField("description", "测试描述")
	writer.WriteField("status", "private")
	writer.WriteField("tags", "测试,单元测试")
	writer.WriteField("duration", "180.5")
	writer.Close()

	// 创建请求
	req := httptest.NewRequest("POST", "/api/v1/videos", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// 执行请求
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var response struct {
		Code int         `json:"code"`
		Data model.Video `json:"data"`
	}
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 0, response.Code)
	assert.Equal(t, 180.5, response.Data.Duration)
	assert.NotEmpty(t, response.Data.CoverURL)
	assert.True(t, strings.HasPrefix(response.Data.CoverURL, "/uploads/cover_"))
}

// TestVideoHandler_UploadInvalidCover 测试上传无效封面图
func TestVideoHandler_UploadInvalidCover(t *testing.T) {
	r, cleanup := setupTestRouter(t)
	defer cleanup()

	token := getTestToken(t)

	// 创建过大的封面图
	largeCoverContent := make([]byte, 3*1024*1024) // 3MB

	// 创建multipart请求
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// 添加视频文件
	videoPart, err := writer.CreateFormFile("file", "test.mp4")
	if err != nil {
		t.Fatal(err)
	}
	videoPart.Write([]byte("fake video content"))

	// 添加过大的封面图
	coverPart, err := writer.CreateFormFile("cover", "large.jpg")
	if err != nil {
		t.Fatal(err)
	}
	coverPart.Write(largeCoverContent)

	// 添加其他字段
	writer.WriteField("title", "测试视频")
	writer.WriteField("duration", "180.5")
	writer.Close()

	// 创建请求
	req := httptest.NewRequest("POST", "/api/v1/videos", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// 执行请求
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]interface{}
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), response["code"])
	assert.Contains(t, response["msg"], "封面图大小不能超过2MB")
}
