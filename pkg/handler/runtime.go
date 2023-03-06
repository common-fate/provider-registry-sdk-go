package handler

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

func (r *Client) FetchResources(ctx context.Context, name string, contx map[string]any) (*msg.LoadResponse, error) {
	response, err := r.Executor.Execute(ctx, msg.LoadResources{Task: name, Ctx: contx})
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

func (r *Client) Grant(ctx context.Context, subject string, target msg.Target) error {
	_, err := r.Executor.Execute(ctx, msg.Grant{Subject: subject, Target: target})
	return err
}

func (r *Client) Revoke(ctx context.Context, subject string, target msg.Target) error {
	_, err := r.Executor.Execute(ctx, msg.Revoke{Subject: subject, Target: target})
	return err
}
