package jsontool

import (
	"fmt"
	"os"
	"testing"

	"github.com/tidwall/sjson"
)

func TestConditionalMod(t *testing.T) {

	data, err := os.ReadFile("./data/Activity.json")
	if err != nil {
		panic(err)
	}
	js := string(data)
	mSibling, _ := FamilyTree(js)

	paths := GetFieldPaths(js, "ActivityWeight", mSibling) // get all paths which contains field 'ActivityWeight'
	fmt.Println(paths)

	mFS := GetSiblingPath(js, "ActivityWeight", "StartDate", mSibling) // get all valid siblings for each 'ActivityWeight' path
	fmt.Println(mFS)

	mFSs := GetSiblingsPath(js, "ActivityWeight", mSibling, "StartDate", "FinishDate", "CreationDate")
	fmt.Println(mFSs)

	for k, v := range mFS {
		js, _ = sjson.Set(js, k, 1000)                   // modify existing field 1
		js, _ = sjson.Set(js, v, "2002-09-13")           // modify existing field 2
		js, _ = sjson.Set(js, NewSibling(k, "HELLO"), 1) // create new field
	}

	os.WriteFile("./data/ATest.json", []byte(js), os.ModePerm)
}

func TestFamilyTree(t *testing.T) {

	data, err := os.ReadFile("./data/Activity.json")
	if err != nil {
		panic(err)
	}

	js := string(data)
	mSibling, mFT := FamilyTree(js)

	for l := 0; l < 1024; l++ {
		if len(mSibling[l]) > 0 {
			fmt.Println(mSibling[l])
		}
	}

	fmt.Println()

	for k, v := range mFT {
		if len(v) > 0 {
			fmt.Println(k, v)
		}
	}

	// for i := 0; i < 1024; i++ {
	// 	if fields := mSibling[i]; len(fields) > 0 {
	// 		// fmt.Println(fields)
	// 		for _, path := range fields {
	// 			if r := gjson.Get(js, path); r.Exists() {
	// 				value := r.String()

	// 				literal := regexp.QuoteMeta(value) // ***

	// 				re := regexp.MustCompile(literal)
	// 				locPairsAll := re.FindAllStringIndex(js, -1)
	// 				fmt.Println(path, locPairsAll)

	// 				// simple
	// 				field := pathLeaf(path)
	// 				re = regexp.MustCompile(fmt.Sprintf(`"%s":\s*"?%s"?`, field, literal))
	// 				locPairsSimple := re.FindAllStringIndex(js, -1)
	// 				for i := 0; i < len(locPairsSimple); i++ {
	// 					if js[locPairsSimple[i][1]-1] == '"' {
	// 						locPairsSimple[i][1] -= 1
	// 					}
	// 				}
	// 				fmt.Println(path, locPairsSimple)

	// 				// array
	// 				// re = regexp.MustCompile(fmt.Sprintf(`"%s":\s*%s`, field, literal))

	// 				fmt.Println()
	// 			}
	// 		}
	// 	}
	// }
}
