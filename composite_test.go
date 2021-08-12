package jsontool

import (
	"fmt"
	"os"
	"testing"
)

func TestComposite(t *testing.T) {
	data, err := os.ReadFile("./data/FlattenTest.json")
	if err != nil {
		panic(err)
	}
	jsonStr := string(data)
	m, err := FlattenObject(jsonStr)
	fmt.Println(len(m), err)
	jsonstr := Composite(m)
	os.WriteFile("./data/FlattenTest_Composite.json", []byte(jsonstr), os.ModePerm)
}
