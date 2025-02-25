package service

import (
	"context"
	"errors"
	"time"
	"video-platform/internal/model"
	"video-platform/pkg/database"
	"video-platform/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *model.RegisterRequest) (*model.User, error) {
	// 检查用户名是否已存在
	collection := database.GetCollection("users")
	var existUser model.User
	err := collection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&existUser)
	if err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 创建新用户
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:        primitive.NewObjectID(),
		Username:  req.Username,
		Password:  hashedPassword,
		Email:     req.Email,
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *model.LoginRequest) (*model.User, string, error) {
	collection := database.GetCollection("users")
	var user model.User
	err := collection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
	if err != nil {
		return nil, "", errors.New("用户不存在")
	}

	// 验证密码
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, "", errors.New("密码错误")
	}

	// 生成 JWT token
	token, err := utils.GenerateToken(user.ID.Hex(), user.Username)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

// GetUserByID 根据ID获取用户信息
func (s *UserService) GetUserByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	collection := database.GetCollection("users")
	var user model.User
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
} 