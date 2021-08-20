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
	m, err := Flatten(data)
	fmt.Println(len(m), err)
	str := Composite(m)
	os.WriteFile("./data/FlattenTest_Composite.json", []byte(str), os.ModePerm)
}
