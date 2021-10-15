package jsontool

import (
	"fmt"
	"os"
	"testing"

	"github.com/tidwall/gjson"
)

func TestGetPathByLoc(t *testing.T) {
	data, err := os.ReadFile("./data/Activity.json")
	if err != nil {
		panic(err)
	}
	js := string(data)
	mSibling, mFT := FamilyTree(js)

	mPathRoom := make(map[string][][]int)

	for i := 0; i < 1024; i++ {
		if fields := mSibling[i]; len(fields) > 0 {
			// fmt.Println(fields)
			for _, path := range fields {
				if r := gjson.Get(js, path); r.Exists() {
					value := r.String()
					mPathRoom[path] = ZipInts(IndexAll(js, value))
				}
			}
		}
	}

	fmt.Println()

	// fmt.Println(mPathRoom)
	for k, v := range mPathRoom {
		fmt.Println(k, v)
	}

	fmt.Println()

	for k, v := range mFT {
		if len(v) > 0 {
			fmt.Println(k, v)
		}
	}

}
