package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"video-platform/internal/model"

	"github.com/gin-gonic/gin"
)

func TestRegister(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	r := gin.Default()
	r.POST("/register", Register)

	// 测试成功注册
	body := model.RegisterRequest{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, body: %s", w.Code, w.Body.String())
		return
	}

	var response struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data *model.User `json:"data"`
	}
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
		return
	}

	if response.Code != 0 {
		t.Errorf("Expected code 0, got %d", response.Code)
	}
	if response.Data == nil {
		t.Error("Expected user data not to be nil")
		return
	}
	if response.Data.Username != body.Username {
		t.Errorf("Expected username %s, got %s", body.Username, response.Data.Username)
	}
}

func TestLogin(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	r := gin.Default()
	r.POST("/register", Register)
	r.POST("/login", Login)

	// 先注册一个用户
	registerBody := model.RegisterRequest{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
	}
	jsonRegisterBody, _ := json.Marshal(registerBody)
	reqRegister := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonRegisterBody))
	reqRegister.Header.Set("Content-Type", "application/json")
	wRegister := httptest.NewRecorder()
	r.ServeHTTP(wRegister, reqRegister)

	// 测试登录
	loginBody := model.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	jsonLoginBody, _ := json.Marshal(loginBody)
	reqLogin := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonLoginBody))
	reqLogin.Header.Set("Content-Type", "application/json")
	wLogin := httptest.NewRecorder()
	r.ServeHTTP(wLogin, reqLogin)

	if wLogin.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", wLogin.Code)
	}

	var response struct {
		Code int `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			User  *model.User `json:"user"`
			Token string      `json:"token"`
		} `json:"data"`
	}
	err := json.Unmarshal(wLogin.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Code != 0 {
		t.Errorf("Expected code 0, got %d", response.Code)
	}
	if response.Data.Token == "" {
		t.Error("Expected token not to be empty")
	}
} 