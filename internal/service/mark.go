package service

import (
	"context"
	"time"
	"video-platform/internal/model"
	"video-platform/pkg/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MarkService 标记服务接口
type MarkService interface {
	AddMark(ctx context.Context, userID string, mark *model.Mark) error
	GetMarks(ctx context.Context, userID, videoID string) ([]model.Mark, error)
	UpdateMark(ctx context.Context, userID string, markID primitive.ObjectID, mark *model.Mark) error
	DeleteMark(ctx context.Context, userID string, markID primitive.ObjectID) error

	AddAnnotation(ctx context.Context, annotation *model.Annotation) error
	GetAnnotations(ctx context.Context, markID primitive.ObjectID) ([]model.Annotation, error)
	UpdateAnnotation(ctx context.Context, userID string, annotationID primitive.ObjectID, annotation *model.Annotation) error
	DeleteAnnotation(ctx context.Context, userID string, annotationID primitive.ObjectID) error

	AddNote(ctx context.Context, note *model.Note) error
	GetNotes(ctx context.Context, videoID string) ([]model.Note, error)
	UpdateNote(ctx context.Context, userID string, noteID primitive.ObjectID, note *model.Note) error
	DeleteNote(ctx context.Context, userID string, noteID primitive.ObjectID) error
}

// markServiceImpl 标记服务实现
type markServiceImpl struct {
	collection string
}

// NewMarkService 创建标记服务实例
func NewMarkService() MarkService {
	return &markServiceImpl{
		collection: "marks",
	}
}

// AddMark 添加标记
func (s *markServiceImpl) AddMark(ctx context.Context, userID string, mark *model.Mark) error {
	mark.ID = primitive.NewObjectID()
	mark.UserID = userID // 设置用户ID
	collection := database.GetCollection(s.collection)
	mark.CreatedAt = time.Now()
	mark.UpdatedAt = time.Now()
	_, err := collection.InsertOne(ctx, mark)
	return err
}

// GetMarks 获取标记列表
func (s *markServiceImpl) GetMarks(ctx context.Context, userID string, videoID string) ([]model.Mark, error) {
	collection := database.GetCollection(s.collection)
	cursor, err := collection.Find(ctx, gin.H{"user_id": userID, "video_id": videoID})
	if err != nil {
		return nil, err
	}
	var marks []model.Mark
	if err := cursor.All(ctx, &marks); err != nil {
		return nil, err
	}

	// 获取每个标记对应的注释
	for i := range marks {
		annotations, err := s.GetAnnotations(ctx, marks[i].ID)
		if err != nil {
			// 如果获取注释失败，不影响标记的返回，只是该标记的注释为空
			marks[i].Annotations = []model.Annotation{}
			continue
		}
		marks[i].Annotations = annotations
	}

	return marks, nil
}

// AddAnnotation 添加注释
func (s *markServiceImpl) AddAnnotation(ctx context.Context, annotation *model.Annotation) error {
	annotation.ID = primitive.NewObjectID()
	collection := database.GetCollection("annotations")
	annotation.CreatedAt = time.Now()
	annotation.UpdatedAt = time.Now()
	_, err := collection.InsertOne(ctx, annotation)
	return err
}

// GetAnnotations 获取注释
func (s *markServiceImpl) GetAnnotations(ctx context.Context, markID primitive.ObjectID) ([]model.Annotation, error) {
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
func (s *markServiceImpl) AddNote(ctx context.Context, note *model.Note) error {
	note.ID = primitive.NewObjectID()
	collection := database.GetCollection("notes")
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()
	_, err := collection.InsertOne(ctx, note)
	return err
}

// GetNotes 获取笔记列表
func (s *markServiceImpl) GetNotes(ctx context.Context, videoID string) ([]model.Note, error) {
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
func (s *markServiceImpl) UpdateMark(ctx context.Context, userID string, markID primitive.ObjectID, mark *model.Mark) error {
	collection := database.GetCollection(s.collection)
	mark.UpdatedAt = time.Now()
	_, err := collection.UpdateOne(ctx, gin.H{"_id": markID, "user_id": userID}, gin.H{"$set": mark})
	return err
}

// DeleteMark 删除标记
func (s *markServiceImpl) DeleteMark(ctx context.Context, userID string, markID primitive.ObjectID) error {
	collection := database.GetCollection(s.collection)
	_, err := collection.DeleteOne(ctx, gin.H{"_id": markID, "user_id": userID})
	return err
}

// UpdateAnnotation 更新注释
func (s *markServiceImpl) UpdateAnnotation(ctx context.Context, userID string, annotationID primitive.ObjectID, annotation *model.Annotation) error {
	collection := database.GetCollection("annotations")
	annotation.UpdatedAt = time.Now()
	_, err := collection.UpdateOne(ctx, gin.H{"_id": annotationID, "user_id": userID}, gin.H{"$set": annotation})
	return err
}

// DeleteAnnotation 删除注释
func (s *markServiceImpl) DeleteAnnotation(ctx context.Context, userID string, annotationID primitive.ObjectID) error {
	collection := database.GetCollection("annotations")
	_, err := collection.DeleteOne(ctx, gin.H{"_id": annotationID, "user_id": userID})
	return err
}

// UpdateNote 更新笔记
func (s *markServiceImpl) UpdateNote(ctx context.Context, userID string, noteID primitive.ObjectID, note *model.Note) error {
	collection := database.GetCollection("notes")
	note.UpdatedAt = time.Now()
	_, err := collection.UpdateOne(ctx, gin.H{"_id": noteID, "user_id": userID}, gin.H{"$set": note})
	return err
}

// DeleteNote 删除笔记
func (s *markServiceImpl) DeleteNote(ctx context.Context, userID string, noteID primitive.ObjectID) error {
	collection := database.GetCollection("notes")
	_, err := collection.DeleteOne(ctx, gin.H{"_id": noteID, "user_id": userID})
	return err
}
