package jsontool

import (
	"os"
	"testing"
)

func TestTrimFields(t *testing.T) {

	data, err := os.ReadFile("./data/otflevel.json")
	if err != nil {
		panic(err)
	}

	js := TrimFields(string(data), true, true, true, true)
	os.WriteFile("./data/otflevel.json", []byte(js), os.ModePerm)
}
