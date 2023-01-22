package utils

import (
	"fmt"
	"log"
	"testing"
)

func TestCryptRandString(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name    string
		args    args
		wantLen int
	}{
		{
			name:    "Generate random string",
			args:    args{n: 32},
			wantLen: 32,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CryptRandString(tt.args.n)
			fmt.Println(got)
			log.Println("random string:", got)
			if len(got) != tt.wantLen {
				t.Errorf("Want length: %d, got length: %d", tt.wantLen, len(got))
			}
		})
	}
}
