package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
	"video-platform/config"
	"video-platform/internal/model"
	"video-platform/pkg/database"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	testDB   string
	testLock sync.Mutex
	testMark = fmt.Sprintf("test_%d", time.Now().UnixNano())
)

func init() {
	// 使用唯一的测试数据库名
	testDB = fmt.Sprintf("video_platform_test_%d", time.Now().UnixNano())
	config.GlobalConfig.MongoDB.Database = testDB
}

// setupTest 初始化测试环境并返回清理函数
func setupTest(t *testing.T) func() {
	testLock.Lock()
	defer testLock.Unlock()

	// 清理测试数据
	cleanTestData(t)

	// 确保上传目录存在
	uploadDir := config.GlobalConfig.Storage.UploadDir
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		t.Fatalf("创建上传目录失败: %v", err)
	}

	// 返回清理函数
	return func() {
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

// TestMain 测试主函数
func TestMain(m *testing.M) {
	// 初始化配置
	if err := config.Init(); err != nil {
		fmt.Printf("配置初始化失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化数据库连接
	if err := database.InitMongoDB(); err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}

	// 运行测试前清理
	cleanAllTestData()

	// 运行测试
	code := m.Run()

	// 测试结束后清理
	cleanAllTestData()

	// 关闭数据库连接
	database.CloseMongoDB()

	os.Exit(code)
}

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

// cleanDirectory 清理指定目录下的所有文件
func cleanDirectory(dir string) {
	files, err := os.ReadDir(dir)
	if err == nil {
		for _, file := range files {
			if !file.IsDir() { // 只删除文件，不删除目录
				os.Remove(filepath.Join(dir, file.Name()))
			}
		}
	}
}

func TestVideoService_Upload(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	service := NewVideoService()
	ctx := context.Background()

	// 创建测试文件
	testFile := createTestFile(t)
	defer os.Remove(testFile.Name())

	// 创建一个真实的文件句柄
	file, err := os.Open(testFile.Name())
	if err != nil {
		t.Fatalf("打开测试文件失败: %v", err)
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		t.Fatalf("获取文件信息失败: %v", err)
	}

	// 创建 multipart.FileHeader
	fileHeader := &multipart.FileHeader{
		Filename: "test.mp4",
		Size:     fileInfo.Size(),
	}

	// 将真实文件路径存储在上下文中
	ctx = context.WithValue(ctx, "testFilePath", testFile.Name())

	info := model.Video{
		Title:       "测试视频",
		Description: "这是一个测试视频",
	}

	video, err := service.Upload(ctx, fileHeader, nil, info)
	if err != nil {
		t.Fatalf("上传视频失败: %v", err)
	}

	if video.Title != info.Title {
		t.Errorf("期望标题为 %s，实际为 %s", info.Title, video.Title)
	}
}

func TestVideoService_GetList(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	service := NewVideoService()
	ctx := context.Background()

	query := model.VideoQuery{
		Page:      1,
		PageSize:  10,
		SortBy:    "created_at",
		SortOrder: "desc",
	}

	list, err := service.GetList(ctx, query)
	if err != nil {
		t.Fatalf("获取视频列表失败: %v", err)
	}

	if list == nil {
		t.Error("返回的视频列表为空")
	}
}

func TestVideoService_GetByID(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	service := NewVideoService()
	ctx := context.Background()

	// 先创建一个视频记录
	testFile := createTestFile(t)
	defer os.Remove(testFile.Name())

	// 获取文件信息
	fileInfo, err := testFile.Stat()
	if err != nil {
		t.Fatalf("获取文件信息失败: %v", err)
	}

	// 创建 multipart.FileHeader
	fileHeader := &multipart.FileHeader{
		Filename: "test.mp4",
		Size:     fileInfo.Size(),
	}

	// 将测试文件路径存储在上下文中
	ctx = context.WithValue(ctx, "testFilePath", testFile.Name())

	info := model.Video{
		Title:       "测试视频",
		Description: "这是一个测试视频",
	}

	video, err := service.Upload(ctx, fileHeader, nil, info)
	if err != nil {
		t.Fatalf("上传视频失败: %v", err)
	}

	// 测试获取视频
	result, err := service.GetByID(ctx, video.ID.Hex())
	if err != nil {
		t.Fatalf("获取视频失败: %v", err)
	}

	if result.Title != info.Title {
		t.Errorf("期望标题为 %s，实际为 %s", info.Title, result.Title)
	}
}

func TestVideoService_Update(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	service := NewVideoService()
	ctx := context.Background()

	// 先创建一个视频记录
	video := createTestVideo(t, service)

	// 更新视频信息
	updateInfo := model.Video{
		Title:       "更新后的标题",
		Description: "更新后的描述",
	}

	err := service.Update(ctx, video.ID.Hex(), updateInfo)
	if err != nil {
		t.Fatalf("更新视频失败: %v", err)
	}

	// 验证更新结果
	updated, err := service.GetByID(ctx, video.ID.Hex())
	if err != nil {
		t.Fatalf("获取视频失败: %v", err)
	}

	if updated.Title != updateInfo.Title {
		t.Errorf("期望标题为 %s，实际为 %s", updateInfo.Title, updated.Title)
	}
}

func TestVideoService_Delete(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	service := NewVideoService()
	ctx := context.Background()

	// 先创建一个视频记录
	video := createTestVideo(t, service)

	// 删除视频
	err := service.Delete(ctx, video.ID.Hex())
	if err != nil {
		t.Fatalf("删除视频失败: %v", err)
	}

	// 验证视频已被删除
	_, err = service.GetByID(ctx, video.ID.Hex())
	if err == nil {
		t.Error("视频应该已被删除")
	}
}

// createTestFile 创建测试文件
func createTestFile(t *testing.T) *os.File {
	// 确保上传目录存在
	uploadDir := config.GlobalConfig.Storage.UploadDir
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		t.Fatalf("创建上传目录失败: %v", err)
	}

	// 创建测试文件
	file, err := os.CreateTemp(uploadDir, "test-*.mp4")
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 写入测试内容
	content := []byte("fake video content")
	if _, err := file.Write(content); err != nil {
		t.Fatalf("写入测试文件失败: %v", err)
	}

	if err := file.Sync(); err != nil {
		t.Fatalf("同步文件失败: %v", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		t.Fatalf("重置文件指针失败: %v", err)
	}

	return file
}

// createTestVideoFile 创建测试视频文件
func createTestVideoFile(t *testing.T) *os.File {
	// 确保上传目录存在
	uploadDir := config.GlobalConfig.Storage.UploadDir
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		t.Fatalf("创建上传目录失败: %v", err)
	}

	// 创建测试文件
	file, err := os.CreateTemp(uploadDir, "test-*.mp4")
	if err != nil {
		t.Fatalf("创建测试视频文件失败: %v", err)
	}

	// 写入测试内容
	content := []byte("fake video content")
	if _, err := file.Write(content); err != nil {
		t.Fatalf("写入测试视频文件失败: %v", err)
	}

	if err := file.Sync(); err != nil {
		t.Fatalf("同步文件失败: %v", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		t.Fatalf("重置文件指针失败: %v", err)
	}

	return file
}

// 辅助函数：创建测试视频
func createTestVideo(t *testing.T, service VideoService) *model.Video {
	ctx := context.Background()
	testFile := createTestFile(t)
	defer os.Remove(testFile.Name())

	// 获取文件信息
	fileInfo, err := testFile.Stat()
	if err != nil {
		t.Fatalf("获取文件信息失败: %v", err)
	}

	fileHeader := &multipart.FileHeader{
		Filename: "test.mp4",
		Size:     fileInfo.Size(),
	}

	// 将测试文件路径存储在上下文中
	ctx = context.WithValue(ctx, "testFilePath", testFile.Name())

	info := model.Video{
		Title:       "测试视频",
		Description: "这是一个测试视频",
		UserID:      testMark,
		Tags:        []string{"test"},
		Status:      model.VideoStatusPublic,
	}

	video, err := service.Upload(ctx, fileHeader, nil, info)
	if err != nil {
		t.Fatalf("上传视频失败: %v", err)
	}

	return video
}

// 添加新的测试用例

func TestVideoService_BatchOperation(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()
	service := NewVideoService()
	ctx := context.Background()

	// 先创建两个测试视频
	video1 := createTestVideo(t, service)
	video2 := createTestVideo(t, service)

	// 测试批量更新状态
	req := model.BatchOperationRequest{
		IDs:    []string{video1.ID.Hex(), video2.ID.Hex()},
		Action: "update_status",
		Status: model.VideoStatusPrivate,
	}

	result, err := service.BatchOperation(ctx, req)
	if err != nil {
		t.Fatalf("批量操作失败: %v", err)
	}

	if result.SuccessCount != 2 {
		t.Errorf("期望成功数量为2，实际为%d", result.SuccessCount)
	}

	// 验证状态是否更新
	video, err := service.GetByID(ctx, video1.ID.Hex())
	if err != nil {
		t.Fatalf("获取视频失败: %v", err)
	}
	if video.Status != model.VideoStatusPrivate {
		t.Errorf("期望状态为%s，实际为%s", model.VideoStatusPrivate, video.Status)
	}
}

func TestVideoService_UpdateThumbnail(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()
	service := NewVideoService()
	ctx := context.Background()

	// 创建测试视频
	video := createTestVideo(t, service)

	// 创建测试缩略图
	coverFile := createTestImageFile(t, "thumbnail_test.jpg")
	defer os.Remove(coverFile.Name())

	// 获取文件信息
	fileInfo, err := coverFile.Stat()
	if err != nil {
		t.Fatalf("获取文件信息失败: %v", err)
	}

	fileHeader := &multipart.FileHeader{
		Filename: "test.jpg",
		Size:     fileInfo.Size(),
	}

	// 将测试文件路径存储在上下文中
	ctx = context.WithValue(ctx, "testFilePath", coverFile.Name())

	// 更新缩略图
	thumbnailURL, err := service.UpdateThumbnail(ctx, video.ID.Hex(), fileHeader)
	if err != nil {
		t.Fatalf("更新缩略图失败: %v", err)
	}

	if thumbnailURL == "" {
		t.Error("缩略图URL不应为空")
	}
}

func TestVideoService_GetStats(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()
	service := NewVideoService()
	ctx := context.Background()

	// 创建测试视频
	video := createTestVideo(t, service)

	// 增加一些统计数据
	err := service.IncrementStats(ctx, video.ID.Hex(), "views")
	if err != nil {
		t.Fatalf("增加观看次数失败: %v", err)
	}

	// 获取统计信息
	stats, err := service.GetStats(ctx, video.ID.Hex())
	if err != nil {
		t.Fatalf("获取统计信息失败: %v", err)
	}

	if stats.Views != 1 {
		t.Errorf("期望观看次数为1，实际为%d", stats.Views)
	}
}

// 辅助函数：创建测试图片文件
func createTestImageFile(t *testing.T, name string) *os.File {
	file, err := os.CreateTemp(config.GlobalConfig.Storage.UploadDir, name)
	if err != nil {
		t.Fatalf("创建测试图片文件失败: %v", err)
	}

	content := []byte("fake image content")
	if _, err := file.Write(content); err != nil {
		t.Fatalf("写入测试图片文件失败: %v", err)
	}

	if err := file.Sync(); err != nil {
		t.Fatalf("同步文件失败: %v", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		t.Fatalf("重置文件指针失败: %v", err)
	}

	return file
}

// TestVideoService_Upload_WithCover 测试上传带封面图的视频
func TestVideoService_Upload_WithCover(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()
	service := NewVideoService()
	ctx := context.Background()

	// 创建测试视频文件
	videoFile := createTestVideoFile(t)
	defer os.Remove(videoFile.Name())

	// 创建测试封面图文件
	coverFile := createTestImageFile(t, "cover_test.jpg")
	defer os.Remove(coverFile.Name())

	// 获取文件信息
	videoInfo, err := videoFile.Stat()
	if err != nil {
		t.Fatalf("获取视频文件信息失败: %v", err)
	}
	coverInfo, err := coverFile.Stat()
	if err != nil {
		t.Fatalf("获取封面图文件信息失败: %v", err)
	}

	// 创建文件头
	videoHeader := &multipart.FileHeader{
		Filename: "test.mp4",
		Size:     videoInfo.Size(),
	}
	coverHeader := &multipart.FileHeader{
		Filename: "cover.jpg",
		Size:     coverInfo.Size(),
	}

	// 设置测试文件路径
	ctx = context.WithValue(ctx, "testFilePath", videoFile.Name())
	ctx = context.WithValue(ctx, "testCoverPath", coverFile.Name())

	// 创建视频信息
	info := model.Video{
		Title:       "测试视频",
		Description: "测试描述",
		Status:      model.VideoStatusPrivate,
		Tags:        []string{"测试", "单元测试"},
	}

	// 上传视频和封面图
	video, err := service.Upload(ctx, videoHeader, coverHeader, info)
	if err != nil {
		t.Fatalf("上传视频失败: %v", err)
	}

	// 验证结果
	if video.Title != info.Title {
		t.Errorf("期望标题为 %s，实际为 %s", info.Title, video.Title)
	}
	if video.CoverURL == "" {
		t.Error("封面图URL不应为空")
	}
	if !strings.HasPrefix(video.CoverURL, "/uploads/cover_") {
		t.Errorf("封面图URL格式错误: %s", video.CoverURL)
	}
}

// TestVideoService_Upload_InvalidCover 测试上传无效封面图
func TestVideoService_Upload_InvalidCover(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()
	service := NewVideoService()
	ctx := context.Background()

	// 创建测试视频文件
	videoFile := createTestVideoFile(t)
	defer os.Remove(videoFile.Name())

	// 获取视频文件信息
	videoInfo, err := videoFile.Stat()
	if err != nil {
		t.Fatalf("获取视频文件信息失败: %v", err)
	}

	// 创建视频文件头
	videoHeader := &multipart.FileHeader{
		Filename: "test.mp4",
		Size:     videoInfo.Size(),
	}

	// 测试无效格式的封面图
	invalidCoverHeader := &multipart.FileHeader{
		Filename: "test.txt",
		Size:     1024, // 1KB，确保不会触发大小限制
	}

	// 设置测试文件路径
	ctx = context.WithValue(ctx, "testFilePath", videoFile.Name())

	// 创建视频信息
	info := model.Video{
		Title:       "测试视频",
		Description: "测试描述",
	}

	// 尝试上传无效格式的封面图
	_, err = service.Upload(ctx, videoHeader, invalidCoverHeader, info)
	if err == nil {
		t.Error("应该返回无效图片格式的错误")
	}
	if !strings.Contains(err.Error(), "不支持的图片格式") {
		t.Errorf("期望错误信息包含'不支持的图片格式'，实际为：%v", err)
	}

	// 创建过大的封面图文件
	largeCoverFile := createLargeTestFile(t, 3*1024*1024) // 3MB
	defer os.Remove(largeCoverFile.Name())

	// 获取大文件信息
	coverInfo, err := largeCoverFile.Stat()
	if err != nil {
		t.Fatalf("获取封面图文件信息失败: %v", err)
	}

	// 创建大文件的文件头
	largeCoverHeader := &multipart.FileHeader{
		Filename: "large.jpg",
		Size:     coverInfo.Size(),
	}

	ctx = context.WithValue(ctx, "testCoverPath", largeCoverFile.Name())

	// 尝试上传过大的封面图
	_, err = service.Upload(ctx, videoHeader, largeCoverHeader, info)
	if err == nil {
		t.Error("应该返回封面图过大的错误")
	}
	if !strings.Contains(err.Error(), "封面图大小不能超过2MB") {
		t.Errorf("期望错误信息包含'封面图大小不能超过2MB'，实际为：%v", err)
	}
}

// 辅助函数：创建指定大小的测试文件
func createLargeTestFile(t *testing.T, size int64) *os.File {
	file, err := os.CreateTemp(config.GlobalConfig.Storage.UploadDir, "test-large-*.jpg")
	if err != nil {
		t.Fatalf("创建大文件失败: %v", err)
	}

	// 写入指定大小的数据
	if err := file.Truncate(size); err != nil {
		t.Fatalf("调整文件大小失败: %v", err)
	}

	if err := file.Sync(); err != nil {
		t.Fatalf("同步文件失败: %v", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		t.Fatalf("重置文件指针失败: %v", err)
	}

	return file
}
