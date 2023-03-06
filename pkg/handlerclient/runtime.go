package handlerclient

import (
	"context"
	"encoding/json"

	"github.com/common-fate/provider-registry-sdk-go/pkg/msg"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

// Enforce build errors if the runtimes don't meet the interface
var _ Executor = &Lambda{}
var _ Executor = &Local{}

type Client struct {
	Executor Executor
}

func (r *Client) FetchResources(ctx context.Context, req msg.LoadResources) (*msg.LoadResponse, error) {
	response, err := r.Executor.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	var lr msg.LoadResponse
	err = json.Unmarshal(response.Response, &lr)
	if err != nil {
		return nil, err
	}

	return &lr, nil
}

func (r *Client) Describe(ctx context.Context) (*providerregistrysdk.DescribeResponse, error) {
	response, err := r.Executor.Execute(ctx, msg.Describe{})
	if err != nil {
		return nil, err
	}

	var dr providerregistrysdk.DescribeResponse

	err = json.Unmarshal(response.Response, &dr)
	if err != nil {
		return nil, err
	}

	return &dr, nil
}

func (r *Client) Grant(ctx context.Context, req msg.Grant) (*msg.GrantResponse, error) {
	response, err := r.Executor.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	var gr msg.GrantResponse

	err = json.Unmarshal(response.Response, &gr)
	if err != nil {
		return nil, err
	}

	return &gr, nil
}

func (r *Client) Revoke(ctx context.Context, req msg.Revoke) error {
	_, err := r.Executor.Execute(ctx, req)
	return err
}
