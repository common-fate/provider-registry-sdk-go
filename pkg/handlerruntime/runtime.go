package handlerruntime

import (
	"context"

	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

// Enforce build errors if the runtimes don't meet the interface
var _ Runtime = &Lambda{}
var _ Runtime = &Local{}

type Runtime interface {
	FetchResources(ctx context.Context, name string, contx interface{}) (resources LoadResourceResponse, err error)
	Describe(ctx context.Context) (describeResponse *providerregistrysdk.DescribeResponse, err error)
	Grant(ctx context.Context, subject string, target Target) (err error)
	Revoke(ctx context.Context, subject string, target Target) (err error)
}
