package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"` // 响应码：0成功，1失败
	Msg  string      `json:"msg"`  // 响应信息
	Data interface{} `json:"data"` // 响应数据
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, httpCode int, msg string) {
	c.JSON(httpCode, Response{
		Code: 1,
		Msg:  msg,
		Data: nil,
	})
}
