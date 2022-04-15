package jsontool

import (
	"os"
	"testing"
	"time"

	"github.com/digisan/gotk/misc"
)

func TestJSONBreakArrCont(t *testing.T) {
	defer misc.TrackTime(time.Now())

	bytes, err := os.ReadFile("./data/Activities.json")
	failOnErr("%v", err)
	str := string(bytes)

	values, ok := BreakArr(str)
	fPln(ok)
	for _, v := range values {
		fPln(v)
	}
}

func TestJSONBreakBlkContV2(t *testing.T) {
	defer misc.TrackTime(time.Now())

	if bytes, err := os.ReadFile("./data/Activity.json"); err == nil {
		// str := Fmt(string(bytes), "  ")
		str := string(bytes)
		_, cont := SglEleBlkCont(str)
		names, values := BreakMulEleBlkV2(cont)
		for i, name := range names {
			fPln(MkSglEleBlk(name, values[i], true))
			fPln(" ------------------------------------------ ")
		}
	}
}
