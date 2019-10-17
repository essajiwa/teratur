package grpc

import "testing"

func TestServe(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			"Test 1 [Success]",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Serve(); (err != nil) != tt.wantErr {
				t.Errorf("Serve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
