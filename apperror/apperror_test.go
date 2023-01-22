package apperror

import (
	"reflect"
	"testing"
)

func TestFromError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 Error
	}{
		{
			name: "Convert error to its original value",
			args: args{
				err: Error{
					Message: "test",
				},
			},
			want: true,
			want1: Error{
				Message: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := FromError(tt.args.err)
			if got != tt.want {
				t.Errorf("FromError() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("FromError() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
