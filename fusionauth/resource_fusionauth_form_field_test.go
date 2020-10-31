package fusionauth

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_validateKey(t *testing.T) {
	type args struct {
		i interface{}
		k string
	}
	tests := []struct {
		name         string
		args         args
		wantWarnings []string
		wantErrors   []error
	}{
		{
			name: "predefined",
			args: args{
				i: "user.birthDate",
				k: "key",
			},
			wantErrors:   nil,
			wantWarnings: nil,
		},
		{
			name: "custom",
			args: args{
				i: "user.data.",
				k: "key",
			},
			wantErrors:   nil,
			wantWarnings: nil,
		},
		{
			name: "invalid type",
			args: args{
				i: false,
				k: "key",
			},

			wantErrors:   []error{fmt.Errorf(`expected type of "key" to be string`)},
			wantWarnings: nil,
		},
		{
			name: "invalid custom",
			args: args{
				i: "user.invalid.",
				k: "key",
			},
			wantErrors:   []error{fmt.Errorf(`valid options for "key" are: ["registration.username" "user.birthDate" "user.email" "user.firstname" "user.fullName" "user.lastName" "user.middleName" "user.mobilePhone" "user.password" "user.username"] or start with ["user.data." "registration.data."]`)},
			wantWarnings: nil,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			gotWarnings, gotErrors := validateKey(tt.args.i, tt.args.k)
			if !reflect.DeepEqual(gotWarnings, tt.wantWarnings) {
				t.Errorf("validateKey() gotWarnings = %v, want %v", gotWarnings, tt.wantWarnings)
			}
			if len(gotErrors) != len(tt.wantErrors) {
				t.Errorf("validateKey() gotErrors = %v, want %v", gotErrors, tt.wantErrors)
			}
		})
	}
}

func Test_validateRegex(t *testing.T) {
	type args struct {
		i interface{}
		k string
	}
	tests := []struct {
		name         string
		args         args
		wantWarnings []string
		wantErrors   []error
	}{
		{
			name: "valid",
			args: args{
				i: "a*b",
				k: "test",
			},
			wantErrors:   nil,
			wantWarnings: nil,
		},
		{
			name: "invalid",
			args: args{
				i: "[",
				k: "test",
			},
			wantErrors:   []error{fmt.Errorf("error parsing regexp: missing closing ]: `[`")},
			wantWarnings: nil,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			gotWarnings, gotErrors := validateRegex(tt.args.i, tt.args.k)
			if !reflect.DeepEqual(gotWarnings, tt.wantWarnings) {
				t.Errorf("validateRegex() gotWarnings = %v, want %v", gotWarnings, tt.wantWarnings)
			}
			if len(gotErrors) != len(tt.wantErrors) {
				t.Errorf("validateKey() gotErrors = %v, want %v", gotErrors, tt.wantErrors)
			}
		})
	}
}
