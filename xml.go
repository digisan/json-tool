package jsontool

import (
	"encoding/json"

	"github.com/clbanning/mxj"
)

// Cvt2XML :
func Cvt2XML(str string, mav map[string]interface{}) string {

	var m interface{}
	json.Unmarshal([]byte(str), &m)
	bytes, err := mxj.AnyXmlIndent(m, "", "    ", "")
	failOnErr("%v", err)
	xmlstr := string(bytes)
	xmlstr = sReplaceAll(xmlstr, "<>", "")
	xmlstr = sReplaceAll(xmlstr, "</>", "")
	xmlstr = sTrim(xmlstr, " \t\n")

	attrs := []string{}
	for a, v := range mav {
		attrs = append(attrs, fSf(`%s="%v"`, a, v))
	}
	if p := sIndex(xmlstr, ">"); len(attrs) > 0 {
		xmlstr = xmlstr[:p] + " " + sJoin(attrs, " ") + xmlstr[p:]
	}

	return xmlstr
}
