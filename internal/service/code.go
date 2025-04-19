package service

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"video-platform/pkg/sms"
)

type CodeService interface {
	Send(ctx context.Context, biz, number string) error
	Verify(ctx context.Context, biz, number, code string) (bool, error)
}

type codeServiceImpl struct {
	sms sms.Service
}

func NewCodeSerivce(sms sms.Service) CodeService {
	return &codeServiceImpl{sms: sms}
}

func (c *codeServiceImpl) Send(ctx context.Context, biz, number string) error {
	code := c.genCode()
	if err := c.SetCode(ctx, biz, number, code); err != nil {
		slog.Error("set code error", "error", err.Error(), "biz", biz, "number", number, "code", code)
		return err
	}
	if err := c.sms.Send(ctx, "SMS_474870192", []sms.Param{{Name: "code", Value: code}}, number); err != nil {
		// redis set 成功，sms 发送失败 不能刪除 redis key 因为错误有可能是超时错误... 即短信发送成功，但是返回超时
		// 解决方案一：重试 让调用者自己决定重试方案 即sms 缺陷：用户重复收到验证码；一直重复一直失败，系统负载高
		// 解决方案二：
		slog.Error("send sms error", "error", err.Error(), "biz", biz, "number", number, "code", code)
		return err
	}
	return nil
}

func (c *codeServiceImpl) Verify(ctx context.Context, biz, number, code string) (bool, error) {
	return c.VerifyCode(ctx, biz, number, code)
}

func (c *codeServiceImpl) genCode() string {
	// 生成6位數隨機驗證碼 0 - 999999
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (r *codeServiceImpl) SetCode(ctx context.Context, biz, number, code string) error {
	return nil
}

func (r *codeServiceImpl) VerifyCode(ctx context.Context, biz, number, code string) (bool, error) {
	return true, nil
}
