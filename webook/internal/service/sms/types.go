package sms

import "context"

type Service interface {
	Send(ctx context.Context, tpl string, numbers ...string) error
}
