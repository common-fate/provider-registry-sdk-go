package bootstrapper

import (
	_ "embed"
	"testing"

	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

func TestBootstrapStackOutput_CloudFormationURL(t *testing.T) {
	type fields struct {
		AssetsBucket string
		Region       string
	}
	type args struct {
		p providerregistrysdk.Provider
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "ok",
			fields: fields{
				AssetsBucket: "common-fate-bootstrap-123",
				Region:       "us-west-2",
			},
			args: args{
				p: providerregistrysdk.Provider{
					Name:      "test",
					Publisher: "common-fate",
					Version:   "v1.0.0",
				},
			},
			want: "https://common-fate-bootstrap-123.s3.us-west-2.amazonaws.com/registry.commonfate.io/v1alpha1/providers/common-fate/test/v1.0.0/cloudformation.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bso := BootstrapStackOutput{
				AssetsBucket: tt.fields.AssetsBucket,
				Region:       tt.fields.Region,
			}
			if got := bso.CloudFormationURL(tt.args.p); got != tt.want {
				t.Errorf("BootstrapStackOutput.CloudFormationURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
