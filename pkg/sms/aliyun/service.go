package aliyun

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"

	"video-platform/config"
	sms2 "video-platform/pkg/sms"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
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
func NewAliyunClient() *sms.Client {
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
	// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
	String := func(str string) *string {
		return &str
	}
	config := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId: String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")),
		//// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
		AccessKeySecret: String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")),

		Endpoint: String(config.GlobalConfig.SMS.Endpoint),
		RegionId: String(config.GlobalConfig.SMS.RegionID),
	}
	_result, err := sms.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	return _result
}
