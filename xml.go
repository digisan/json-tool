package jsontool

import (
	"encoding/json"

	"github.com/clbanning/mxj"
)

// Cvt2XML :
func Cvt2XML(str string, mav map[string]any) string {

	var m any
	json.Unmarshal([]byte(str), &m)
	bytes, err := mxj.AnyXmlIndent(m, "", "    ", "")
	failOnErr("%v", err)
	xmlStr := string(bytes)
	xmlStr = sReplaceAll(xmlStr, "<>", "")
	xmlStr = sReplaceAll(xmlStr, "</>", "")
	xmlStr = sTrim(xmlStr, " \t\n")

	attrs := []string{}
	for a, v := range mav {
		attrs = append(attrs, fSf(`%s="%v"`, a, v))
	}
	if p := sIndex(xmlStr, ">"); len(attrs) > 0 {
		xmlStr = xmlStr[:p] + " " + sJoin(attrs, " ") + xmlStr[p:]
	}

	return xmlStr
}
