package jsontool

import (
	"fmt"
	"testing"
)

func TestFmt(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			args: args{
				str: `{    "Id": 1,    "Name": "Ahmad,   Ahmad", "Age": "21" 	 }`,
			},
			want: `{"Id":1,"Name":"Ahmad,   Ahmad","Age":"2	1"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := TryFmtStr(tt.args.str, "  ")
			fmt.Println("try:", got)

			got = FmtStr(tt.args.str, "  ")
			fmt.Println("fmt:", got)

			got, err := FmtJS(tt.args.str)
			fmt.Println("fmtJS:", got)
			fmt.Println("err:", err)
		})
	}
}
