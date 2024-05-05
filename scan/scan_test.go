package scan

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	. "github.com/digisan/go-generics"
)

func TestScanJsonLine(t *testing.T) {

	opt := OptLineProc{
		Fn_KV:          nil,
		Fn_KV_Str:      nil,
		Fn_KV_Obj_Open: nil,
		Fn_KV_Arr_Open: nil,
		Fn_Obj:         nil,
		Fn_Arr:         nil,
		Fn_Elem:        nil,
		Fn_Elem_Str:    nil,
	}

	des, err := os.ReadDir("../data")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, de := range des {

		if In(de.Name(), "FlattenTest.json") {
			continue
		}

		fPath := filepath.Join("../data", de.Name())
		fOut := filepath.Join("./", de.Name())

		fmt.Printf("testing... %s\n", fPath)
		fmt.Println(ScanJsonLine(fPath, fOut, opt)) // *** //

		// original copying check
		data1, err := os.ReadFile(fPath)
		if err != nil {
			panic(err)
		}
		data2, err := os.ReadFile(fOut)
		if err != nil {
			panic(err)
		}
		if string(data1) != string(data2) {
			panic("NOT copying equally")
		}
	}
}
