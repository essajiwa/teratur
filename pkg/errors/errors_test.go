// Package errors is a custom error message implementation with file path and file number
package errors

import (
	"os"
	"testing"
)

func TestSet(t *testing.T) {
	gp := os.Getenv("GOPATH")
	os.Setenv("GOPATH", "")
	defer os.Setenv("GOPATH", gp)

	type args struct {
		e error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Success",
			args{
				e: New("error euy"),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Wrap(tt.args.e); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
