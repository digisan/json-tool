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

	str := CompositeIncl(m, "array.0.c", "object.a", "object1.object11.a")
	os.WriteFile("./data/FlattenTest_CompositeIncl.json", []byte(str), os.ModePerm)

	str = CompositeExcl(m, "array.0.c", "object.a", "object1.object11.a")
	os.WriteFile("./data/FlattenTest_CompositeExcl.json", []byte(str), os.ModePerm)
}
