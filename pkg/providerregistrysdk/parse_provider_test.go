package providerregistrysdk

import (
	"reflect"
	"testing"
)

func TestParseProvider(t *testing.T) {
	tests := []struct {
		name    string
		give    string
		want    Provider
		wantErr bool
	}{
		{
			name: "ok",
			give: "common-fate/provider@v0.1.0",
			want: Provider{
				Publisher: "common-fate",
				Name:      "provider",
				Version:   "v0.1.0",
			},
		},
		{
			name:    "invalid with spaces",
			give:    "common-fate/some provider@v0.1.0",
			wantErr: true,
		},
		{
			name:    "invalid with no version",
			give:    "common-fate/provider@",
			wantErr: true,
		},
		{
			name:    "invalid with bad format",
			give:    "invalid",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseProvider(tt.give)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}
