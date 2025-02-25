package model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 用户模型
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string            `bson:"username" json:"username"`         // 用户名
	Password  string            `bson:"password" json:"-"`               // 密码（加密存储）
	Email     string            `bson:"email" json:"email"`             // 邮箱
	Status    int              `bson:"status" json:"status"`           // 状态 1:正常 2:禁用
	CreatedAt time.Time         `bson:"created_at" json:"createdAt"`    // 创建时间
	UpdatedAt time.Time         `bson:"updated_at" json:"updatedAt"`    // 更新时间
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Email    string `json:"email" binding:"required,email"`
}

// UserProfile 用户资料
type UserProfile struct {
	Nickname    string    `json:"nickname" bson:"nickname"`
	Avatar      string    `json:"avatar" bson:"avatar"`
	Description string    `json:"description" bson:"description"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
} 