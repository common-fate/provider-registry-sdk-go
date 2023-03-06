package handler

import (
	"context"

	"github.com/common-fate/provider-registry-sdk-go/pkg/msg"
)

// MockExecutor can be used to test the handler runtime client.
type MockExecutor struct {
	Result *msg.Result
	Err    error
}

func (mi MockExecutor) Execute(ctx context.Context, request msg.Request) (*msg.Result, error) {
	return mi.Result, mi.Err
}
