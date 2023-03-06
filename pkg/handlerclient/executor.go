package handlerclient

import (
	"context"

	"github.com/common-fate/provider-registry-sdk-go/pkg/msg"
)

// Executors can execute RPC calls to Handlers.
type Executor interface {
	Execute(ctx context.Context, request msg.Request) (*msg.Result, error)
}
