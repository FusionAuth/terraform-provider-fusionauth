package fusionauth

import (
	"reflect"
	"testing"
)

func Test_intMapToStringMap(t *testing.T) {
	type args struct {
		intMap map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "FA Issues #1482",
			args: args{
				intMap: map[string]interface{}{
					"ar": "Test",
				},
			},
			want: map[string]string{
				"ar": "Test",
			},
		},
		{
			name: "FA Issues #1482",
			args: args{
				intMap: map[string]interface{}{
					"ar":    "Test",
					"aaass": 2,
				},
			},
			want: map[string]string{
				"ar": "Test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := intMapToStringMap(tt.args.intMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("intMapToStringMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
