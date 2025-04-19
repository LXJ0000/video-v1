package aliyun

import (
	"context"
	sms2 "github.com/LXJ0000/go-backend/internal/usecase/sms"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sms "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	"os"
	"testing"
)

func TestService_Send(t *testing.T) {
	client, err := CreateClient()
	if err != nil {
		t.Error(err)
	}
	service := NewService("1889161073986679", "小南的编程世界", client)
	ctx := context.Background()
	if err := service.Send(ctx, "SMS_474870192", []sms2.Param{{Name: "code", Value: "666666"}}, "18126934563"); err != nil {
		t.Error(err)
	}

}

func CreateClient() (_result *sms.Client, _err error) {
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
	// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
	config := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId: String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")),
		//// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
		AccessKeySecret: String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")),

		Endpoint: String("dysmsapi.aliyuncs.com"),
		RegionId: String("cn-shenzhen"),
	}
	_result = &sms.Client{}
	_result, _err = sms.NewClient(config)
	return _result, _err
}

func String(str string) *string {
	return &str
}
