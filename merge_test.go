package jsontool

import (
	"os"
	"testing"
)

func TestMerge(t *testing.T) {
	bytes1, err := os.ReadFile("./data/single1.json")
	failOnErr("%v", err)
	bytes2, err := os.ReadFile("./data/single2.json")
	failOnErr("%v", err)
	bytes3, err := os.ReadFile("./data/single3.json")
	failOnErr("%v", err)

	merged := MergeSgl(string(bytes1), string(bytes2), string(bytes3))
	failOnErrWhen(!IsValid(merged), "%v", fEf("Invalid JSON"))
	fPln(merged)
}

func TestScalarSel(t *testing.T) {
	bytes, err := os.ReadFile("./data/itemResults.json")
	failOnErr("%v", err)
	jsonstr := string(bytes)
	result := ScalarSelX(jsonstr, "School", "YrLevel", "Test Item RefID")
	fPln(result)
}
