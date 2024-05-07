package scan

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

	const DIR = "../data"
	des, err := os.ReadDir(DIR)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, de := range des {

		if !strings.HasSuffix(de.Name(), ".json") {
			continue
		}

		fPath := filepath.Join(DIR, de.Name())
		fOut := filepath.Join("./", de.Name())

		fmt.Printf("processing... %s\n", fPath)
		paths, err := ScanJsonLine(fPath, fOut, opt) // *** //
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("paths count... %d\n", len(paths))

		// original copying check
		data1, err := os.ReadFile(fPath)
		if err != nil {
			log.Fatalln(err)
		}
		data2, err := os.ReadFile(fOut)
		if err != nil {
			log.Fatalln(err)
		}
		if string(data1) != string(data2) {
			log.Fatalln("NOT copying equally")
		} else {
			fmt.Println("successful")
		}
	}
}
