package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mark 用户自定义标记模型
type Mark struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"userId"`          // 用户ID
	VideoID     string             `bson:"video_id" json:"videoId"`        // 视频ID
	Timestamp   float64            `bson:"timestamp" json:"timestamp"`     // 时间戳
	Content     string             `bson:"content" json:"content"`         // 标记内容
	Annotations []Annotation       `bson:"annotations" json:"annotations"` // 关联的注释
}

// Annotation 注释模型
type Annotation struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID  string             `bson:"user_id" json:"userId"`  // 用户ID
	MarkID  primitive.ObjectID `bson:"mark_id" json:"markId"`  // 关联的标记ID
	Content string             `bson:"content" json:"content"` // 注释内容
}

// Note 笔记模型
type Note struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id" json:"userId"`      // 用户ID
	VideoID   string             `bson:"video_id" json:"videoId"`    // 视频ID
	Timestamp float64            `bson:"timestamp" json:"timestamp"` // 时间戳
	Content   string             `bson:"content" json:"content"`     // 笔记内容
}
