package jsontool

import "testing"

func TestCvt2XML(t *testing.T) {
	out := MkSglEleBlk("ROOT", "~~~", true)
	fPln(out)
	mav := map[string]interface{}{"a": "b", "c": 12}
	xmlstr := Cvt2XML(out, mav)
	fPln(xmlstr)
}
