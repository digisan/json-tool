package jsontool

import (
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
			jsonstr := string(bytes)

			if !IsValid(jsonstr) {
				os.WriteFile(fSf("debug_%s.json", jsonfile), []byte(jsonstr), 0666)
				panic("error on MkJSON")
			}

			//if jsonfile == "CensusCollection_0.xml" {
			// os.WriteFile(fSf("record_%s.json", jsonfile), []byte(jsonstr), 0666)
			//}
		}
		return nil
	})
}

func TestMinimize(t *testing.T) {
	type args struct {
		jsonstr string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			args: args {
				jsonstr: `{    "Id": 1,    "Name": "Ahmad,   Ahmad", "Age": "2	1" 		 }`,
			},
			want: `{"Id":1,"Name":"Ahmad,   Ahmad","Age":"2	1"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Minimize(tt.args.jsonstr); got != tt.want {
				t.Errorf("Minimize() = %v, want %v", got, tt.want)
			}
		})
	}
}
