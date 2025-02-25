package service

import (
	"context"
	"fmt"
	"os"
	"testing"
	"video-platform/config"
	"video-platform/internal/model"
	"video-platform/pkg/database"
)

var markService MarkService

func setup() {
	// 初始化配置
	if err := config.Init(); err != nil {
		fmt.Printf("配置初始化失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化测试数据库连接
	ctx := context.Background()
	if err := database.InitMongoDB(ctx, config.GlobalConfig.MongoDB, true); err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	markService = NewMarkService()
}

func teardown() {
	// 清理测试数据
	ctx := context.Background()
	if err := database.CleanupTestData(ctx); err != nil {
		fmt.Printf("清理测试数据失败: %v\n", err)
	}
}

func TestAddMark(t *testing.T) {
	setup()
	defer teardown()

	mark := &model.Mark{
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Mark",
	}

	err := markService.AddMark(context.Background(), "test_user_id", mark)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if mark.ID.IsZero() {
		t.Fatal("Expected mark ID to be set")
	}
}

func TestGetMarks(t *testing.T) {
	setup()
	defer teardown()

	mark := &model.Mark{
		UserID:    "test_user_id",
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Mark",
	}
	err := markService.AddMark(context.Background(), mark.UserID, mark)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	marks, err := markService.GetMarks(context.Background(), mark.UserID, mark.VideoID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(marks) == 0 {
		t.Fatal("Expected to get marks, got none")
	}
}

func TestAddAnnotation(t *testing.T) {
	setup()
	defer teardown()

	mark := &model.Mark{
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Mark",
	}
	markService.AddMark(context.Background(), "test_user_id", mark)

	annotation := &model.Annotation{
		MarkID:  mark.ID,
		Content: "Test Annotation",
	}

	err := markService.AddAnnotation(context.Background(), annotation)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if annotation.ID.IsZero() {
		t.Fatal("Expected annotation ID to be set")
	}
}

func TestGetAnnotations(t *testing.T) {
	setup()
	defer teardown()

	mark := &model.Mark{
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Mark",
	}
	markService.AddMark(context.Background(), "test_user_id", mark)

	annotation := &model.Annotation{
		MarkID:  mark.ID,
		Content: "Test Annotation",
	}
	markService.AddAnnotation(context.Background(), annotation)

	annotations, err := markService.GetAnnotations(context.Background(), mark.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(annotations) == 0 {
		t.Fatal("Expected to get annotations, got none")
	}
}

func TestAddNote(t *testing.T) {
	setup()
	defer teardown()

	note := &model.Note{
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Note",
	}

	err := markService.AddNote(context.Background(), note)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if note.ID.IsZero() {
		t.Fatal("Expected note ID to be set")
	}
}

func TestGetNotes(t *testing.T) {
	setup()
	defer teardown()

	note := &model.Note{
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Note",
	}
	markService.AddNote(context.Background(), note)

	notes, err := markService.GetNotes(context.Background(), "test_video_id")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(notes) == 0 {
		t.Fatal("Expected to get notes, got none")
	}
}

func TestAddMarkWithUserID(t *testing.T) {
	setup()
	defer teardown()

	mark := &model.Mark{
		UserID:    "test_user_id",
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Mark",
	}

	err := markService.AddMark(context.Background(), mark.UserID, mark)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if mark.ID.IsZero() {
		t.Fatal("Expected mark ID to be set")
	}
}

func TestGetMarksWithUserID(t *testing.T) {
	setup()
	defer teardown()

	mark := &model.Mark{
		UserID:    "test_user_id",
		VideoID:   "test_video_id",
		Timestamp: 123.45,
		Content:   "Test Mark",
	}
	markService.AddMark(context.Background(), mark.UserID, mark)

	marks, err := markService.GetMarks(context.Background(), mark.UserID, "test_video_id")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(marks) == 0 {
		t.Fatal("Expected to get marks, got none")
	}
}
