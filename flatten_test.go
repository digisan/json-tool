package jsontool

import (
	"fmt"
	"os"
	"testing"

	"github.com/digisan/go-generics/so"
	"github.com/digisan/go-generics/str"
)

func TestFlattenObject(t *testing.T) {

	data, err := os.ReadFile("./data/FlattenTest.json")
	if err != nil {
		panic(err)
	}
	js := string(data)

	m, err := FlattenStr(js)
	fmt.Println(len(m), err)
	// I := 0
	// for k, v := range m {
	// 	fmt.Printf("%02d --- %v: %v\n", I, k, v)
	// 	I++
	// }
	ks, _ := so.Map2KVs(m, nil, nil)
	fmt.Println(len(ks))

	fmt.Println("--------------------------------------------------")

	paths, _ := GetLeavesPathOrderly(js)
	fmt.Println(len(paths))
	// for i, p := range paths {
	// 	fmt.Printf("%02d --- %v: %v\n", i, p, "tbd")
	// }

	fmt.Println("--------------------------------------------------")

	fmt.Println(str.Equal(paths, ks))
}
