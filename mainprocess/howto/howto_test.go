package howto

import (
	"testing"
)

func Test_ditDahToHowTo(t *testing.T) {
	type args struct {
		ditdah string
	}
	tests := []struct {
		name      string
		args      args
		wantHowto string
	}{
		// TODO: Add test cases.
		{
			name: "ab",
			args: args{
				ditdah: ".- -...",
			},
			wantHowto: "down up, down 2 3 up, 2 3, down 2 3 up, down up, down up, down up",
		},
		{
			name: "cq",
			args: args{
				ditdah: "-.-. --.-",
			},
			wantHowto: "down 2 3 up, down up, down 2 3 up, down up, 2 3, down 2 3 up, down 2 3 up, down up, down 2 3 up",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotHowto := ditDahToHowTo(tt.args.ditdah); gotHowto != tt.wantHowto {
				t.Errorf("ditDahToHowTo() = %v, want %v", gotHowto, tt.wantHowto)
			}
		})
	}
}
