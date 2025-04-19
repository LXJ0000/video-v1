package local

import (
	"context"
	"fmt"
	"video-platform/pkg/sms"
)

type Service struct {
	appID    string
	signName string
}

func NewService() sms.Service {
	return &Service{}
}

func (s *Service) Send(ctx context.Context, templateID string, args []sms.Param, numbers ...string) error {
	fmt.Println(args)
	return nil
}
