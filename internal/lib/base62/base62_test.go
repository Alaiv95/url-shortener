package base62

import "testing"

func TestConvertNum(t *testing.T) {
	type args struct {
		num int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "0",
			args: args{num: 0},
			want: "",
		},
		{
			name: "1",
			args: args{num: 1},
			want: "1",
		},
		{
			name: "1523",
			args: args{num: 1523},
			want: "oz",
		},
		{
			name: "1231423423423",
			args: args{num: 1231423423423},
			want: "lG9xbNJ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertNum(tt.args.num); got != tt.want {
				t.Errorf("ConvertNum() = %v, want %v", got, tt.want)
			}
		})
	}
}
