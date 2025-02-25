package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"video-platform/internal/model"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func setupRouter() {
	gin.SetMode(gin.TestMode)
	router = gin.Default()
	InitRoutes(router)
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

func TestAddMarkHandlerWithUserID(t *testing.T) {
	setupRouter()

	body := bytes.NewBufferString(`{"videoId":"test_video_id","timestamp":123.45,"content":"Test Mark"}`)
	req, _ := http.NewRequest("POST", "/api/v1/marks/test_user_id/test_video_id", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

func TestGetMarksHandlerWithUserID(t *testing.T) {
	setupRouter()

	mark := model.Mark{
		UserID:    "test_user_id",
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Mark",
	}
	markService.AddMark(context.Background(), mark.UserID, &mark)

	req, _ := http.NewRequest("GET", "/api/v1/marks/test_user_id/test_video_id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

func TestAddAnnotationHandler(t *testing.T) {
	setupRouter()

	mark := model.Mark{
		UserID:    "test_user_id",
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Mark",
	}
	markService.AddMark(context.Background(), mark.UserID, &mark)

	body := bytes.NewBufferString(`{"markId":"` + mark.ID.Hex() + `","content":"Test Annotation"}`)
	req, _ := http.NewRequest("POST", "/api/v1/marks/"+mark.ID.Hex()+"/annotations", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

func TestGetAnnotationsHandler(t *testing.T) {
	setupRouter()

	mark := model.Mark{
		UserID:    "test_user_id",
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Mark",
	}
	markService.AddMark(context.Background(), mark.UserID, &mark)

	annotation := model.Annotation{
		MarkID:  mark.ID,
		Content: "Test Annotation",
	}
	markService.AddAnnotation(context.Background(), &annotation)

	req, _ := http.NewRequest("GET", "/api/v1/marks/"+mark.ID.Hex()+"/annotations", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

func TestAddNoteHandler(t *testing.T) {
	setupRouter()

	body := bytes.NewBufferString(`{"videoId":"test_video_id","timestamp":123.45,"content":"Test Note"}`)
	req, _ := http.NewRequest("POST", "/api/v1/notes/test_user_id/test_video_id", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

func TestGetNotesHandler(t *testing.T) {
	setupRouter()

	note := model.Note{
		UserID:    "test_user_id",
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Note",
	}
	markService.AddNote(context.Background(), &note)

	req, _ := http.NewRequest("GET", "/api/v1/notes/test_user_id/test_video_id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

func TestExportMarksHandler(t *testing.T) {
	setupRouter()

	mark := model.Mark{
		UserID:    "test_user_id",
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Mark",
	}
	markService.AddMark(context.Background(), mark.UserID, &mark)

	req, _ := http.NewRequest("GET", "/api/v1/videos/export/test_user_id/test_video_id", nil)
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

	mark := model.Mark{
		UserID:    "test_user_id",
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Mark",
	}
	markService.AddMark(context.Background(), mark.UserID, &mark)

	body := bytes.NewBufferString(`{"content":"Updated Mark"}`)
	req, _ := http.NewRequest("PUT", "/api/v1/marks/test_user_id/test_video_id/"+mark.ID.Hex(), body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

func TestDeleteMarkHandler(t *testing.T) {
	setupRouter()

	mark := model.Mark{
		UserID:    "test_user_id",
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Mark",
	}
	markService.AddMark(context.Background(), mark.UserID, &mark)

	req, _ := http.NewRequest("DELETE", "/api/v1/marks/test_user_id/test_video_id/"+mark.ID.Hex(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

func TestUpdateAnnotationHandler(t *testing.T) {
	setupRouter()

	mark := model.Mark{
		UserID:    "test_user_id",
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Mark",
	}
	markService.AddMark(context.Background(), mark.UserID, &mark)

	annotation := model.Annotation{
		MarkID:  mark.ID,
		Content: "Test Annotation",
	}
	markService.AddAnnotation(context.Background(), &annotation)

	body := bytes.NewBufferString(`{"content":"Updated Annotation"}`)
	req, _ := http.NewRequest("PUT", "/api/v1/marks/test_user_id/test_video_id/annotations/"+annotation.ID.Hex(), body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

func TestDeleteAnnotationHandler(t *testing.T) {
	setupRouter()

	mark := model.Mark{
		UserID:    "test_user_id",
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Mark",
	}
	markService.AddMark(context.Background(), mark.UserID, &mark)

	annotation := model.Annotation{
		MarkID:  mark.ID,
		Content: "Test Annotation",
	}
	markService.AddAnnotation(context.Background(), &annotation)

	req, _ := http.NewRequest("DELETE", "/api/v1/marks/test_user_id/test_video_id/annotations/"+annotation.ID.Hex(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

func TestUpdateNoteHandler(t *testing.T) {
	setupRouter()

	note := model.Note{
		UserID:    "test_user_id",
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Note",
	}
	markService.AddNote(context.Background(), &note)

	body := bytes.NewBufferString(`{"content":"Updated Note"}`)
	req, _ := http.NewRequest("PUT", "/api/v1/notes/test_user_id/test_video_id/"+note.ID.Hex(), body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}

func TestDeleteNoteHandler(t *testing.T) {
	setupRouter()

	note := model.Note{
		UserID:    "test_user_id",
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Note",
	}
	markService.AddNote(context.Background(), &note)

	req, _ := http.NewRequest("DELETE", "/api/v1/notes/test_user_id/test_video_id/"+note.ID.Hex(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	checkResponseFormat(t, w)
}
