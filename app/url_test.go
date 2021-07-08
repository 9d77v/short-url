package app

import "testing"

func TestConvertIntToStr(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{10000}, "2Bi"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertIntToStr(tt.args.id); got != tt.want {
				t.Errorf("ConvertIntToStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
