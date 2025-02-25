package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"video-platform/config"
	"video-platform/internal/model"
	"video-platform/internal/service"
	"video-platform/pkg/database"

	"github.com/gin-gonic/gin"
)

var (
	router          *gin.Engine
	markHandler     *MarkHandler
	testMarkService service.MarkService
)

func setupRouter() {
	gin.SetMode(gin.TestMode)
	router = gin.Default()

	// 创建服务实例
	testMarkService = service.NewMarkService()
	markHandler = NewMarkHandler(testMarkService)

	// API v1 分组
	v1 := router.Group("/api/v1")
	{
		// 标记相关路由
		marks := v1.Group("/marks/:userId/:id")
		{
			marks.POST("", markHandler.AddMark)                          // 添加标记
			marks.GET("", markHandler.GetMarks)                          // 获取标记列表
			marks.PUT("/:markId", markHandler.UpdateMark)                // 更新标记
			marks.DELETE("/:markId", markHandler.DeleteMark)             // 删除标记
			marks.POST("/annotations/:markId", markHandler.AddAnnotation) // 添加注释
			marks.GET("/annotations/:markId", markHandler.GetAnnotations) // 获取注释
			marks.PUT("/annotations/:annotationId", markHandler.UpdateAnnotation)    // 更新注释
			marks.DELETE("/annotations/:annotationId", markHandler.DeleteAnnotation) // 删除注释
		}

		// 笔记相关路由
		notes := v1.Group("/notes/:userId/:id")
		{
			notes.POST("", markHandler.AddNote)              // 添加笔记
			notes.GET("", markHandler.GetNotes)              // 获取笔记列表
			notes.PUT("/:noteId", markHandler.UpdateNote)    // 更新笔记
			notes.DELETE("/:noteId", markHandler.DeleteNote) // 删除笔记
		}

		// 导出相关路由
		videos := v1.Group("/videos")
		{
			videos.GET("/export/:userId/:id", markHandler.ExportMarks) // 导出标记、注释和笔记
		}
	}
}

// 添加辅助函数
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func checkResponseFormat(t *testing.T, w *httptest.ResponseRecorder) {
	var response Response
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if w.Code == http.StatusOK {
		if response.Code != 0 {
			t.Errorf("Expected code 0, got %d", response.Code)
		}
		if response.Msg != "success" {
			t.Errorf("Expected msg 'success', got %s", response.Msg)
		}
	} else {
		if response.Code != 1 {
			t.Errorf("Expected code 1, got %d", response.Code)
		}
		if response.Data != nil {
			t.Error("Expected data to be nil for error response")
		}
	}
}

func setupTest(t *testing.T) func() {
	// 初始化配置
	if err := config.Init(); err != nil {
		t.Fatal(err)
	}

	// 初始化数据库
	if err := database.InitMongoDB(context.Background(), config.GlobalConfig.MongoDB, true); err != nil {
		t.Fatal(err)
	}

	// 创建服务和handler实例
	testMarkService = service.NewMarkService()
	markHandler = NewMarkHandler(testMarkService)

	// 返回清理函数
	return func() {
		if err := database.CleanupTestData(context.Background()); err != nil {
			t.Error(err)
		}
	}
}

func TestAddMark(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	mark := &model.Mark{
		UserID:    TestUserID,
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Mark",
	}

	err := testMarkService.AddMark(context.Background(), mark.UserID, mark)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// ... 其他测试代码
}

func TestAddMarkHandlerWithUserID(t *testing.T) {
	setupRouter()
	cleanup := setupTest(t)
	defer cleanup()

	token := getTestToken(t)

	body := bytes.NewBufferString(`{"videoId":"test_video_id","timestamp":123.45,"content":"Test Mark"}`)
	req, _ := http.NewRequest("POST", "/api/v1/marks/test_user_id/test_video_id", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

func TestGetMarksHandlerWithUserID(t *testing.T) {
	setupRouter()
	cleanup := setupTest(t)
	defer cleanup()

	token := getTestToken(t)

	mark := model.Mark{
		UserID:    TestUserID,
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Mark",
	}
	testMarkService.AddMark(context.Background(), mark.UserID, &mark)

	req, _ := http.NewRequest("GET", "/api/v1/marks/test_user_id/test_video_id", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

func TestAddAnnotationHandler(t *testing.T) {
	setupRouter()
	cleanup := setupTest(t)
	defer cleanup()

	token := getTestToken(t)

	mark := model.Mark{
		UserID:    TestUserID,
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Mark",
	}
	testMarkService.AddMark(context.Background(), mark.UserID, &mark)

	body := bytes.NewBufferString(`{"markId":"` + mark.ID.Hex() + `","content":"Test Annotation"}`)
	req, _ := http.NewRequest("POST", "/api/v1/marks/"+mark.ID.Hex()+"/annotations", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

func TestGetAnnotationsHandler(t *testing.T) {
	setupRouter()
	cleanup := setupTest(t)
	defer cleanup()

	token := getTestToken(t)

	mark := model.Mark{
		UserID:    TestUserID,
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Mark",
	}
	testMarkService.AddMark(context.Background(), mark.UserID, &mark)

	annotation := model.Annotation{
		MarkID:  mark.ID,
		Content: "Test Annotation",
	}
	testMarkService.AddAnnotation(context.Background(), &annotation)

	req, _ := http.NewRequest("GET", "/api/v1/marks/"+mark.ID.Hex()+"/annotations", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

func TestAddNoteHandler(t *testing.T) {
	setupRouter()
	cleanup := setupTest(t)
	defer cleanup()

	token := getTestToken(t)

	body := bytes.NewBufferString(`{"videoId":"test_video_id","timestamp":123.45,"content":"Test Note"}`)
	req, _ := http.NewRequest("POST", "/api/v1/notes/test_user_id/test_video_id", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

func TestGetNotesHandler(t *testing.T) {
	setupRouter()
	cleanup := setupTest(t)
	defer cleanup()

	token := getTestToken(t)

	note := model.Note{
		UserID:    TestUserID,
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Note",
	}
	testMarkService.AddNote(context.Background(), &note)

	req, _ := http.NewRequest("GET", "/api/v1/notes/test_user_id/test_video_id", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

func TestExportMarksHandler(t *testing.T) {
	setupRouter()
	cleanup := setupTest(t)
	defer cleanup()

	token := getTestToken(t)

	mark := model.Mark{
		UserID:    TestUserID,
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Mark",
	}
	testMarkService.AddMark(context.Background(), mark.UserID, &mark)

	req, _ := http.NewRequest("GET", "/api/v1/videos/export/test_user_id/test_video_id", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

// 添加新的测试用例
func TestUpdateMarkHandler(t *testing.T) {
	setupRouter()
	cleanup := setupTest(t)
	defer cleanup()

	token := getTestToken(t)

	mark := model.Mark{
		UserID:    TestUserID,
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Mark",
	}
	testMarkService.AddMark(context.Background(), mark.UserID, &mark)

	body := bytes.NewBufferString(`{"content":"Updated Mark"}`)
	req, _ := http.NewRequest("PUT", "/api/v1/marks/test_user_id/test_video_id/"+mark.ID.Hex(), body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

func TestDeleteMarkHandler(t *testing.T) {
	setupRouter()
	cleanup := setupTest(t)
	defer cleanup()

	token := getTestToken(t)

	mark := model.Mark{
		UserID:    TestUserID,
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Mark",
	}
	testMarkService.AddMark(context.Background(), mark.UserID, &mark)

	req, _ := http.NewRequest("DELETE", "/api/v1/marks/test_user_id/test_video_id/"+mark.ID.Hex(), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

func TestUpdateAnnotationHandler(t *testing.T) {
	setupRouter()
	cleanup := setupTest(t)
	defer cleanup()

	token := getTestToken(t)

	mark := model.Mark{
		UserID:    TestUserID,
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Mark",
	}
	testMarkService.AddMark(context.Background(), mark.UserID, &mark)

	annotation := model.Annotation{
		MarkID:  mark.ID,
		Content: "Test Annotation",
	}
	testMarkService.AddAnnotation(context.Background(), &annotation)

	body := bytes.NewBufferString(`{"content":"Updated Annotation"}`)
	req, _ := http.NewRequest("PUT", "/api/v1/marks/test_user_id/test_video_id/annotations/"+annotation.ID.Hex(), body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

func TestDeleteAnnotationHandler(t *testing.T) {
	setupRouter()
	cleanup := setupTest(t)
	defer cleanup()

	token := getTestToken(t)

	mark := model.Mark{
		UserID:    TestUserID,
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Mark",
	}
	testMarkService.AddMark(context.Background(), mark.UserID, &mark)

	annotation := model.Annotation{
		MarkID:  mark.ID,
		Content: "Test Annotation",
	}
	testMarkService.AddAnnotation(context.Background(), &annotation)

	req, _ := http.NewRequest("DELETE", "/api/v1/marks/test_user_id/test_video_id/annotations/"+annotation.ID.Hex(), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

func TestUpdateNoteHandler(t *testing.T) {
	setupRouter()
	cleanup := setupTest(t)
	defer cleanup()

	token := getTestToken(t)

	note := model.Note{
		UserID:    TestUserID,
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Note",
	}
	testMarkService.AddNote(context.Background(), &note)

	body := bytes.NewBufferString(`{"content":"Updated Note"}`)
	req, _ := http.NewRequest("PUT", "/api/v1/notes/test_user_id/test_video_id/"+note.ID.Hex(), body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

func TestDeleteNoteHandler(t *testing.T) {
	setupRouter()
	cleanup := setupTest(t)
	defer cleanup()

	token := getTestToken(t)

	note := model.Note{
		UserID:    TestUserID,
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Note",
	}
	testMarkService.AddNote(context.Background(), &note)

	req, _ := http.NewRequest("DELETE", "/api/v1/notes/test_user_id/test_video_id/"+note.ID.Hex(), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

func TestMain(m *testing.M) {
	// 设置串行执行测试
	flag.Parse()
	if testing.Short() {
		flag.Set("test.parallel", "1")
	}

	// 运行测试
	os.Exit(m.Run())
}
