package service

import (
	"context"
	"video-platform/internal/model"
	"video-platform/pkg/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MarkService 标记服务
type MarkService struct{}

// NewMarkService 创建新的标记服务
func NewMarkService() *MarkService {
	return &MarkService{}
}

// AddMark 添加标记
func (s *MarkService) AddMark(ctx context.Context, userID string, mark *model.Mark) error {
	mark.ID = primitive.NewObjectID()
	mark.UserID = userID // 设置用户ID
	collection := database.GetCollection("marks")
	_, err := collection.InsertOne(ctx, mark)
	return err
}

// GetMarks 获取标记列表
func (s *MarkService) GetMarks(ctx context.Context, userID string, videoID string) ([]model.Mark, error) {
	collection := database.GetCollection("marks")
	cursor, err := collection.Find(ctx, gin.H{"user_id": userID, "video_id": videoID})
	if err != nil {
		return nil, err
	}
	var marks []model.Mark
	if err := cursor.All(ctx, &marks); err != nil {
		return nil, err
	}
	return marks, nil
}

// AddAnnotation 添加注释
func (s *MarkService) AddAnnotation(ctx context.Context, annotation *model.Annotation) error {
	annotation.ID = primitive.NewObjectID()
	collection := database.GetCollection("annotations")
	_, err := collection.InsertOne(ctx, annotation)
	return err
}

// GetAnnotations 获取注释
func (s *MarkService) GetAnnotations(ctx context.Context, markID primitive.ObjectID) ([]model.Annotation, error) {
	collection := database.GetCollection("annotations")
	cursor, err := collection.Find(ctx, gin.H{"mark_id": markID})
	if err != nil {
		return nil, err
	}
	var annotations []model.Annotation
	if err := cursor.All(ctx, &annotations); err != nil {
		return nil, err
	}
	return annotations, nil
}

// AddNote 添加笔记
func (s *MarkService) AddNote(ctx context.Context, note *model.Note) error {
	note.ID = primitive.NewObjectID()
	collection := database.GetCollection("notes")
	_, err := collection.InsertOne(ctx, note)
	return err
}

// GetNotes 获取笔记列表
func (s *MarkService) GetNotes(ctx context.Context, videoID string) ([]model.Note, error) {
	collection := database.GetCollection("notes")
	cursor, err := collection.Find(ctx, gin.H{"video_id": videoID})
	if err != nil {
		return nil, err
	}
	var notes []model.Note
	if err := cursor.All(ctx, &notes); err != nil {
		return nil, err
	}
	return notes, nil
}

// UpdateMark 更新标记
func (s *MarkService) UpdateMark(ctx context.Context, userID string, markID primitive.ObjectID, mark *model.Mark) error {
	collection := database.GetCollection("marks")
	_, err := collection.UpdateOne(ctx, gin.H{"_id": markID, "user_id": userID}, gin.H{"$set": mark})
	return err
}

// DeleteMark 删除标记
func (s *MarkService) DeleteMark(ctx context.Context, userID string, markID primitive.ObjectID) error {
	collection := database.GetCollection("marks")
	_, err := collection.DeleteOne(ctx, gin.H{"_id": markID, "user_id": userID})
	return err
}

// UpdateAnnotation 更新注释
func (s *MarkService) UpdateAnnotation(ctx context.Context, userID string, annotationID primitive.ObjectID, annotation *model.Annotation) error {
	collection := database.GetCollection("annotations")
	_, err := collection.UpdateOne(ctx, gin.H{"_id": annotationID, "user_id": userID}, gin.H{"$set": annotation})
	return err
}

// DeleteAnnotation 删除注释
func (s *MarkService) DeleteAnnotation(ctx context.Context, userID string, annotationID primitive.ObjectID) error {
	collection := database.GetCollection("annotations")
	_, err := collection.DeleteOne(ctx, gin.H{"_id": annotationID, "user_id": userID})
	return err
}

// UpdateNote 更新笔记
func (s *MarkService) UpdateNote(ctx context.Context, userID string, noteID primitive.ObjectID, note *model.Note) error {
	collection := database.GetCollection("notes")
	_, err := collection.UpdateOne(ctx, gin.H{"_id": noteID, "user_id": userID}, gin.H{"$set": note})
	return err
}

// DeleteNote 删除笔记
func (s *MarkService) DeleteNote(ctx context.Context, userID string, noteID primitive.ObjectID) error {
	collection := database.GetCollection("notes")
	_, err := collection.DeleteOne(ctx, gin.H{"_id": noteID, "user_id": userID})
	return err
}
