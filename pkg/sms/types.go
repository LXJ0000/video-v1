package sms

import "context"

type Service interface {
	Send(ctx context.Context, templateID string, args []Param, numbers ...string) error
}

type Param struct {
	Name  string
	Value string
}
