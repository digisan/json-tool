package jsontool

import (
	"os"
	"testing"
	"time"

	"github.com/digisan/gotk"
)

func TestJSONBlkCont(t *testing.T) {
	defer gotk.TrackTime(time.Now())

	bytes, err := os.ReadFile("./data/Activity.json")
	failOnErr("%v", err)
	str := string(bytes)

	val, ok := SglEleAttrVal(str, "RefId", "-")
	fPln(val, ok)

	name, cont := SglEleBlkCont(str)
	fPln("root", name)
	fPln(cont)
	fPln(" ------------------------- ")

	// names, values := JSONBreakBlkCont(cont)
	// for i, name := range names {
	// 	fPln(i, name, ":", values[i])
	// }
}

func TestMkSglEleBlk(t *testing.T) {
	fPln(MkSglEleBlk("name", nil, true))
}
