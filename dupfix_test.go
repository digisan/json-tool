package jsontool

import (
	"os"
	"strings"
	"testing"
)

func TestDupFixOne(t *testing.T) {
	fpath := "./data/out/la-Technologies"
AGAIN:
	data, err := os.ReadFile(fpath + ".json")
	if err != nil {
		panic(err)
	}
	n := 30
	prefix := "\n" + strings.Repeat(" ", n) + "\"asn_skillEmbodied\":"
	jsx, ok := FixOneDupKeyOnce(string(data), prefix)
	os.WriteFile(fpath+".json", []byte(jsx), os.ModePerm)
	if !ok {
		goto AGAIN
	}
}

func TestFixOneDupKey(t *testing.T) {
	fpath := "./data/dupkey"
	data, err := os.ReadFile(fpath + ".json")
	if err != nil {
		panic(err)
	}
	n := 6
	prefix := "\n" + strings.Repeat(" ", n) + "\"Age\":"
	fixed := FixOneDupKey(string(data), prefix)
	os.WriteFile(fpath+"1.json", []byte(fixed), os.ModePerm)
}

func TestRmDupEleOnce(t *testing.T) {
	fpath := "./data/dupele"
	data, err := os.ReadFile(fpath + ".json")
	if err != nil {
		panic(err)
	}
	fixed := RmDupEleOnce(string(data), "root.0.Age")
	os.WriteFile(fpath+"1.json", []byte(fixed), os.ModePerm)
}
