package jsontool

import (
	"encoding/json"
	"strings"

	"github.com/clbanning/mxj"
	"github.com/digisan/gotk/slice/tu8"
)

// IsValid :
func IsValid(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

// Fmt :
func Fmt(jsonstr, indent string) string {
	// jsonmap := make(map[string]interface{})
	var jsonmap interface{}
	json.Unmarshal([]byte(jsonstr), &jsonmap)
	bytes, err := json.MarshalIndent(&jsonmap, "", indent)
	failOnErr("%v", err)
	return string(bytes)
}

// Minimize :
func Minimize(jsonstr string) string {

	var (
		sb     = &strings.Builder{}
		pc     = byte(0)
		quotes = false
	)

	for i := 0; i < len(jsonstr); i++ {
		c := jsonstr[i]
		switch {
		case c == '"' && pc != '\\':
			quotes = !quotes
			sb.WriteByte(c)
		case !quotes:
			if tu8.NotIn(c, ' ', '\t', '\n', '\r') {
				sb.WriteByte(c)
			}
		case quotes:
			sb.WriteByte(c)
		}
		pc = c
	}

	return sb.String()
}

// Cvt2XML :
func Cvt2XML(jsonstr string, mav map[string]interface{}) string {

	var jsonmap interface{}
	json.Unmarshal([]byte(jsonstr), &jsonmap)
	bytes, err := mxj.AnyXmlIndent(jsonmap, "", "    ", "")
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
