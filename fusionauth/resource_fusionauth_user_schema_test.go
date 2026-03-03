package fusionauth

import (
	"context"
	"reflect"
	"testing"
)

func Test_upgradeUserSchemaV0ToV1(t *testing.T) {
	type args struct {
		rawState map[string]any
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]any
		wantErr bool
	}{
		{
			name: "Should handle nil state",
			args: args{
				rawState: nil,
			},
			want:    map[string]any{},
			wantErr: false,
		},
		{
			name: "Should handle empty state",
			args: args{
				rawState: map[string]any{},
			},
			want:    map[string]any{},
			wantErr: false,
		},
		{
			name: "Should not touch other properties",
			args: args{
				rawState: map[string]any{
					"first_name": "John",
					"last_name":  "Doe",
					"username":   "user@example.com",
				},
			},
			want: map[string]any{
				"first_name": "John",
				"last_name":  "Doe",
				"username":   "user@example.com",
			},
			wantErr: false,
		},
		{
			name: "Should remove deprecated state properties",
			args: args{
				rawState: map[string]any{
					"two_factor_delivery": "TextMessage",
					"two_factor_enabled":  "false",
					"two_factor_secret":   "UEBzc3cwcmQ=",
				},
			},
			want:    map[string]any{},
			wantErr: false,
		},
		{
			name: "Should upgrade user.data from TypeMap to TypeString",
			args: args{
				rawState: map[string]any{
					"data": map[string]any{
						"test":                   "string",
						"should":                 "upgrade",
						"numbersAreStillStringy": "2",
					},
				},
			},
			want: map[string]any{
				"data": "{\"numbersAreStillStringy\":\"2\",\"should\":\"upgrade\",\"test\":\"string\"}",
			},
			wantErr: false,
		},
		{
			name: "Should upgrade from V0 to V1",
			args: args{
				rawState: map[string]any{
					"first_name": "John",
					"last_name":  "Doe",
					"username":   "user@example.com",
					"data": map[string]any{
						"test":                   "string",
						"should":                 "upgrade",
						"numbersAreStillStringy": "2",
					},
					"two_factor_delivery": "TextMessage",
					"two_factor_enabled":  "false",
					"two_factor_secret":   "UEBzc3cwcmQ=",
				},
			},
			want: map[string]any{
				"first_name": "John",
				"last_name":  "Doe",
				"username":   "user@example.com",
				"data":       "{\"numbersAreStillStringy\":\"2\",\"should\":\"upgrade\",\"test\":\"string\"}",
			},
			wantErr: false,
		},
		{
			name: "Should handle empty user.data from TypeMap to TypeString",
			args: args{
				rawState: map[string]any{
					"data": map[string]any{},
				},
			},
			want: map[string]any{
				"data": "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := upgradeUserSchemaV0ToV1(context.Background(), tt.args.rawState, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("upgradeUserSchemaV0ToV1() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("upgradeUserSchemaV0ToV1() got = %v, want %v", got, tt.want)
			}
		})
	}
}
