package providerregistrysdk

import (
	"testing"
)

func TestProvider_String(t *testing.T) {
	type fields struct {
		Name      string
		Publisher string
		Version   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ok",
			fields: fields{
				Name:      "provider",
				Publisher: "common-fate",
				Version:   "v1.0.0",
			},
			want: "common-fate/provider@v1.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Provider{
				Name:      tt.fields.Name,
				Publisher: tt.fields.Publisher,
				Version:   tt.fields.Version,
			}
			if got := p.String(); got != tt.want {
				t.Errorf("Provider.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProviderDetail_String(t *testing.T) {
	type fields struct {
		CfnTemplateS3Arn string
		LambdaAssetS3Arn string
		Name             string
		Publisher        string
		Schema           Schema
		Version          string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ok",
			fields: fields{
				Name:      "provider",
				Publisher: "common-fate",
				Version:   "v1.0.0",
			},
			want: "common-fate/provider@v1.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ProviderDetail{
				CfnTemplateS3Arn: tt.fields.CfnTemplateS3Arn,
				LambdaAssetS3Arn: tt.fields.LambdaAssetS3Arn,
				Name:             tt.fields.Name,
				Publisher:        tt.fields.Publisher,
				Schema:           tt.fields.Schema,
				Version:          tt.fields.Version,
			}
			if got := p.String(); got != tt.want {
				t.Errorf("ProviderDetail.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
