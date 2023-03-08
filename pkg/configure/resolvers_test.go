package configure

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/stretchr/testify/assert"
)

func TestConfig_Fill(t *testing.T) {
	type fields struct {
		Values map[string]ConfigValue
	}
	tests := []struct {
		name       string
		fields     fields
		opts       FillOpts
		wantValues map[string]ConfigValue
		wantErr    bool
	}{
		{
			name: "ok",
			fields: fields{
				Values: map[string]ConfigValue{
					"api_url": {},
				},
			},
			opts: FillOpts{
				ConfigResolvers: []Resolver{
					MapResolver{kv: map[string]string{"api_url": "test"}},
				},
			},
			wantValues: map[string]ConfigValue{
				"api_url": {Value: "test"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Values: tt.fields.Values,
			}
			ctx := context.Background()
			if err := c.Fill(ctx, tt.opts); (err != nil) != tt.wantErr {
				t.Errorf("Config.Fill() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.wantValues, c.Values)
		})
	}
}

func TestConfig_CfnParams(t *testing.T) {
	type fields struct {
		Values map[string]ConfigValue
	}
	tests := []struct {
		name   string
		fields fields
		want   []types.Parameter
	}{
		{
			name: "config",
			fields: fields{
				Values: map[string]ConfigValue{
					"api_url": {
						Value: "test",
					},
				},
			},
			want: []types.Parameter{
				{
					ParameterKey:   aws.String("ApiUrl"),
					ParameterValue: aws.String("test"),
				},
			},
		},
		{
			name: "secret",
			fields: fields{
				Values: map[string]ConfigValue{
					"api_url": {
						Secret: true,
						Ref:    "awsssm://some/secret",
					},
				},
			},
			want: []types.Parameter{
				{
					ParameterKey:   aws.String("ApiUrlSecret"),
					ParameterValue: aws.String("awsssm://some/secret"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cv := Config{
				Values: tt.fields.Values,
			}
			got := cv.CfnParams()

			assert.Equal(t, tt.want, got)
		})
	}
}
