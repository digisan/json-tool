package jsontool

import (
	"fmt"
	"os"
	"testing"
)

func TestFieldName(t *testing.T) {
	name := FieldName("children.1.children.2.children.2.children.0.dcterms_title")
	fmt.Println(name)
	name = FieldName("children")
	fmt.Println(name)
}

func TestHasSiblings(t *testing.T) {
	data, err := os.ReadFile("./data/FlattenTest.json")
	if err != nil {
		panic(err)
	}
	js := string(data)
	mSibling, mFamilyTree := FamilyTree(js)

	ok := HasSiblings(js, "array.0.subarray.0.c", mSibling, "g", "f", "h")
	fmt.Println(ok)
	ok = HasSiblings(js, "array.0.subarray.1.aa", mSibling, "cc", "ee", "aa")
	fmt.Println(ok)
	ok = HasSiblings(js, "array.0.subarray", mSibling, "c", "a", "e")
	fmt.Println(ok)

	fmt.Println(PathExists(js, "array.0.subarray.1.aa", mFamilyTree))
	fmt.Println(PathExists(js, "array.0.subarray.2.aaa", mFamilyTree))
	fmt.Println(PathExists(js, "object.aa", mFamilyTree))
}

func TestConditionalMod(t *testing.T) {

	data, err := os.ReadFile("./data/test.json")
	if err != nil {
		panic(err)
	}
	js := string(data)
	mSibling, _ := FamilyTree(js)

	// paths := GetFieldPaths(js, "dcterms_title", mSibling) // get all paths which contains field 'dcterms_title'
	// fmt.Println(paths)

	// mFS := GetSiblingPath(js, "dcterms_title", "asn_statementLabel", mSibling) // get all valid siblings for each 'dcterms_title' path
	// fmt.Println(mFS)

	mFSs := GetSiblingsPath(js, "dcterms_title", mSibling, "asn_statementLabel")
	for fp, sps := range mFSs {
		fmt.Println(fp, sps)
	}

	// for k, v := range mFS {
	// 	js, _ = sjson.Set(js, k, 1000)                   // modify existing field 1
	// 	js, _ = sjson.Set(js, v, "2002-09-13")           // modify existing field 2
	// 	js, _ = sjson.Set(js, NewSibling(k, "HELLO"), 1) // create new field
	// }

	// os.WriteFile("./data/test_out.json", []byte(js), os.ModePerm)
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
