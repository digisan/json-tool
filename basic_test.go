package jsontool

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestIsValid(t *testing.T) {
	dir := "./data/"
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		failOnErrWhen(info == nil, "%v", err)
		if jsonfile := info.Name(); sHasSuffix(jsonfile, ".json") {
			fPln("--->", jsonfile)

			bytes, _ := os.ReadFile(dir + jsonfile)
			str := string(bytes)

			if !IsValidStr(str) {
				os.WriteFile(fSf("debug_%s.json", jsonfile), []byte(str), 0666)
				panic("error on MkJSON")
			}

			//if jsonfile == "CensusCollection_0.xml" {
			// os.WriteFile(fSf("record_%s.json", jsonfile), []byte(str), 0666)
			//}
		}
		return nil
	})
}

func TestTryFmt(t *testing.T) {
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
				str: `{    "Id": 1,    "Name": "Ahmad,   Ahmad", "Age": "2	1" 	 }`,
			},
			want: `{"Id":1,"Name":"Ahmad,   Ahmad","Age":"2	1"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TryFmtStr(tt.args.str, "  ")
			fmt.Println(got)
			got = FmtStr(tt.args.str, "  ")
			fmt.Println(got)
		})
	}
}

func TestMinimize(t *testing.T) {
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
				str: `{    "Id": 1,    "Name": "Ahmad,   Ahmad", "Age": "2	1" 		 }`,
			},
			want: `{"Id":1,"Name":"Ahmad,   Ahmad","Age":"2	1"}`,
		},
		{
			name: "OK1",
			args: args{
				str: `   `,
			},
			want: ``,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TryMinimize(tt.args.str)
			fmt.Println(got)
			got = Minimize(tt.args.str, false)
			fmt.Println(got)
			got = Minimize(tt.args.str, true)
			fmt.Println(got)
		})
	}
}

func TestMarshalRemove(t *testing.T) {

	type Object struct {
		A int
		B string
		C bool
	}

	object := Object{
		A: 1,
		B: "b",
		C: true,
	}

	type Obj struct {
		A int    `json:"a"`
		B string `json:"b"`
		C bool   `json:"c"`
	}

	obj := Obj{
		A: 2,
		B: "bb",
		C: false,
	}

	type args struct {
		v            interface{}
		mFieldOldNew map[string]string
		fields       []string
	}
	tests := []struct {
		name      string
		args      args
		wantBytes []byte
		wantErr   bool
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			args: args{
				v:            object,
				mFieldOldNew: map[string]string{"B": "BB"},
				fields:       []string{"A"},
			},
		},
		{
			name: "OK",
			args: args{
				v:            obj,
				mFieldOldNew: map[string]string{"b": "bbb"},
				fields:       []string{"a"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBytes, _ := MarshalRemove(tt.args.v, tt.args.mFieldOldNew, tt.args.fields...)
			r := string(gotBytes)
			fmt.Println(r)
		})
	}
}
