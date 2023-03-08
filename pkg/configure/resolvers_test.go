package configure

import (
	"context"
	"testing"

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
