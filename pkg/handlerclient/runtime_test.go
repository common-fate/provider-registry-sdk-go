package handlerclient

import (
	"context"
	"testing"

	"github.com/common-fate/provider-registry-sdk-go/pkg/msg"
	"github.com/stretchr/testify/assert"
)

func TestRuntime_FetchResources(t *testing.T) {
	tests := []struct {
		name         string
		invokeResult *msg.Result
		invokeErr    error
		want         *msg.LoadResponse
		wantErr      bool
	}{
		{
			name: "ok",
			invokeResult: &msg.Result{
				Response: []byte(`{"resources": [{"type": "Test", "id": "123"}]}`),
			},
			want: &msg.LoadResponse{
				Resources: []msg.Resource{
					{Type: "Test", ID: "123"},
				},
			},
		},
		{
			name: "invalid json",
			invokeResult: &msg.Result{
				Response: []byte(`bad`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Client{
				Executor: MockExecutor{Result: tt.invokeResult, Err: tt.invokeErr},
			}
			ctx := context.Background()

			got, err := r.FetchResources(ctx, msg.LoadResources{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Runtime.FetchResources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRuntime_Grant(t *testing.T) {
	tests := []struct {
		name         string
		invokeResult *msg.Result
		invokeErr    error
		want         *msg.GrantResponse
		wantErr      bool
	}{
		{
			name: "ok",
			invokeResult: &msg.Result{
				Response: []byte(`{"access_instructions": "test", "state": {"foo": "bar"}}`),
			},
			want: &msg.GrantResponse{
				AccessInstructions: "test",
				State: map[string]any{
					"foo": "bar",
				},
			},
		},
		{
			name: "invalid json",
			invokeResult: &msg.Result{
				Response: []byte(`bad`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Client{
				Executor: MockExecutor{Result: tt.invokeResult, Err: tt.invokeErr},
			}
			ctx := context.Background()

			got, err := r.Grant(ctx, msg.Grant{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Runtime.FetchResources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
