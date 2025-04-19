package aliyun

import (
	"context"
	"encoding/json"
	"log/slog"

	sms2 "github.com/LXJ0000/go-backend/internal/usecase/sms"
	sms "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
)

type Service struct {
	appID    string
	signName string
	client   sms.Client
}

func NewService(appID, signName string, client *sms.Client) sms2.Service {
	return &Service{
		appID:    appID,
		signName: signName,
		client:   *client,
	}
}

func (s *Service) Send(ctx context.Context, templateID string, args []sms2.Param, numbers ...string) error {
	for _, number := range numbers {
		argsMap := make(map[string]string, len(args))
		for _, arg := range args {
			argsMap[arg.Name] = arg.Value
		}
		templateParam, err := json.Marshal(argsMap)
		if err != nil {
			slog.Error("failed to marshal template param", "error", err.Error())
			continue
		}
		templateParamStr := string(templateParam)
		req := sms.SendSmsRequest{
			PhoneNumbers:  &number,
			SignName:      &s.signName,
			TemplateCode:  &templateID,
			TemplateParam: &templateParamStr, // eg. json - "{\"code\":\"1234\"}"
		}
		resp, err := s.client.SendSms(&req)
		if err != nil {
			slog.Error("failed to send sms", "error", err.Error())
			continue
		}
		slog.InfoContext(ctx, "send sms info", "response", resp)
		if resp.Body != nil && *(resp.Body.Message) != "OK" {
			slog.Error("send sms failed", "error", "send sms failed")
		}
	}
	return nil
}
