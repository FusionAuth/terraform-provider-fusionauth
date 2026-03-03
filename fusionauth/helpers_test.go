package fusionauth

import (
	"reflect"
	"testing"
)

func Test_handleStringSliceFromList(t *testing.T) {
	tests := []struct {
		name      string
		list      []any
		want      []string
		wantPanic bool
	}{
		{
			name:      "empty list",
			list:      []any{},
			want:      []string{},
			wantPanic: false,
		},
		{
			name:      "all strings",
			list:      []any{"hello", "world"},
			want:      []string{"hello", "world"},
			wantPanic: false,
		},
		{
			name:      "mixed types",
			list:      []any{"string1", 42, "string2"},
			want:      nil,
			wantPanic: true,
		},
		{
			name:      "nil element",
			list:      []any{"valid", nil, "also valid"},
			want:      []string{"valid", "also valid"},
			wantPanic: false,
		},
		{
			name:      "multiple nil elements",
			list:      []any{nil, "middle", nil},
			want:      []string{"middle"},
			wantPanic: false,
		},
		{
			name:      "all nil elements",
			list:      []any{nil, nil, nil},
			want:      []string{},
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []string
			var didPanic bool

			defer func() {
				if r := recover(); r != nil {
					didPanic = true
				}

				if didPanic != tt.wantPanic {
					t.Errorf("handleStringSliceFromList() panic = %v, wantPanic %v", didPanic, tt.wantPanic)
					return
				}

				if !didPanic && !reflect.DeepEqual(got, tt.want) {
					t.Errorf("handleStringSliceFromList() = %v, want %v", got, tt.want)
				}
			}()

			got = handleStringSliceFromList(tt.list)
		})
	}
}

func Test_intMapToStringMap(t *testing.T) {
	type args struct {
		intMap map[string]any
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "FA Issues #1482",
			args: args{
				intMap: map[string]any{
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
				intMap: map[string]any{
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
