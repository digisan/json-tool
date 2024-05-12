package scan

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/digisan/go-generics"
	jt "github.com/digisan/json-tool"
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

		if de.Name() != "FlattenTest.json" {
			continue
		}

		fPath := filepath.Join(DIR, de.Name())
		fOut := filepath.Join("./", de.Name())

		fmt.Printf("processing... %s\n", fPath)
		_, paths, values, err := AnalyzeJson(fPath) // *** //
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("paths count... %d\nvalues count... %d\n", len(paths), len(values))

		// print each (path, value)
		{
			for i, path := range paths {
				value := values[i]
				fmt.Printf("==> %s -- %v\n", path, value)
			}
		}

		// *** //
		if err := ScanJsonLine(fPath, fOut, opt); err != nil {
			log.Fatalln(err)
		}

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

func TestFlattenJson(t *testing.T) {

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

		// if In(de.Name(),
		// 	"Activities.json",
		// 	"ModulePrerequisites.json",
		// 	"Modules.json",
		// 	"Questions.json",
		// 	"StudentMastery.json",
		// 	"Substrands.json",
		// 	"data.json",
		// 	"example.json",
		// 	"itemResults.json",
		// 	"mixed.json",
		// 	"otflevel.json") {
		// 	continue
		// }

		if de.Name() != "FlattenTest.json" {
			continue
		}

		////////////////////////////////////////////////

		fPath := filepath.Join(DIR, de.Name())
		fmt.Printf("processing... %s\n", fPath)

		m1, err := FlattenJson(fPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("FlattenJson: m1 length %d\n\n", len(m1))

		////////////////////////////////////////////////

		data, err := os.ReadFile(fPath)
		if err != nil {
			panic(err)
		}
		js := string(data)

		m2, err := jt.FlattenStr(js)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("FlattenStr: m2 length %d\n\n", len(m2))

		////////////////////////////////////////////////

		fmt.Println("searching... k1 isn't in m2(FlattenStr)")
		for k1, v1 := range m1 {
			if _, ok := m2[k1]; !ok {
				fmt.Printf("  %s -- %v\n", k1, v1)
			}
		}
		fmt.Println()

		fmt.Println("searching... k2 isn't in m1(FlattenJson)")
		for k2, v2 := range m2 {
			if _, ok := m1[k2]; !ok {
				fmt.Printf("  %s -- %v\n", k2, v2)
			}
		}
		fmt.Println()

		////////////////////////////////////////////////

		oPaths, _ := jt.GetLeavesPathOrderly(js)
		fmt.Printf("GetLeavesPathOrderly: paths_orderly %d\n", len(oPaths))
		// for _, op := range oPaths {
		// 	fmt.Println("ordered path ==>", op)
		// }

		fmt.Println("searching... k1 isn't in oPaths(GetLeavesPathOrderly)")
		for k1, v1 := range m1 {
			if NotIn(k1, oPaths...) {
				fmt.Printf("  %s -- %v\n", k1, v1)
			}
		}
		fmt.Println()

		fmt.Println("searching... k2 isn't in oPaths(GetLeavesPathOrderly)")
		for k2, v2 := range m2 {
			if NotIn(k2, oPaths...) {
				fmt.Printf("  %s -- %v\n", k2, v2)
			}
		}
		fmt.Println()
	}
}

func TestLastFieldOrElem(t *testing.T) {

	_, paths, _, err := AnalyzeJson("../data/FlattenTest.json")
	if err != nil {
		log.Fatalln(err)
	}
	// for _, path := range paths {
	// 	fmt.Println(path)
	// }

	lPaths := LastFieldLines("ee", paths)
	for _, path := range lPaths {
		fmt.Println(path)
	}
}
