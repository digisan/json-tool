package jsontool

import "testing"

func TestCvt2XML(t *testing.T) {
	out := MkSglEleBlk("ROOT", "~~~", true)
	fPln(out)
	mav := map[string]any{"a": "b", "c": 12}
	xmlStr := Cvt2XML(out, mav)
	fPln(xmlStr)
}
