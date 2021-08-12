package jsontool

import (
	"fmt"
	"os"
	"testing"
)

func TestFlattenObject(t *testing.T) {
	data, err := os.ReadFile("./data/FlattenTest.json")
	if err != nil {
		panic(err)
	}
	jsonStr := string(data)
	m, err := FlattenStr(jsonStr)
	fmt.Println(len(m), err)
	I := 0
	for k, v := range m {
		fmt.Printf("%02d --- %v: %v\n", I, k, v)
		I++
	}
}
