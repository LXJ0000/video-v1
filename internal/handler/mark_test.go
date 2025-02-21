package handler

import (
	"bytes"
	"context"
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
}
