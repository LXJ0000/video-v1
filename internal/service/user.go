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

// UserService 用户服务接口
type UserService interface {
	Register(ctx context.Context, req *model.RegisterRequest) (*model.User, error)
	Login(ctx context.Context, req *model.LoginRequest) (*model.User, string, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	UpdateProfile(ctx context.Context, id string, profile *model.UserProfile) error
}

type userService struct {
	collection string
}

// NewUserService 创建用户服务实例
func NewUserService() UserService {
	return &userService{
		collection: "users",
	}
}

// Register 用户注册
func (s *userService) Register(ctx context.Context, req *model.RegisterRequest) (*model.User, error) {
	// 检查用户名是否已存在
	collection := database.GetCollection(s.collection)
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
func (s *userService) Login(ctx context.Context, req *model.LoginRequest) (*model.User, string, error) {
	collection := database.GetCollection(s.collection)
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
func (s *userService) GetByID(ctx context.Context, id string) (*model.User, error) {
	collection := database.GetCollection(s.collection)
	var user model.User
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userService) UpdateProfile(ctx context.Context, id string, profile *model.UserProfile) error {
	// Implementation needed
	return nil
}
