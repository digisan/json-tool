package jsontool

import (
	"log"
	"regexp"
)

// js need to be formatted
func TrimFields(js string, rmNull, rmEmptyStr, rmEmptyObj, rmEmptyArr bool) string {

	rNull := regexp.MustCompile(`[,{](\n)(\s)+"[^"]+":(\s)+null[,\n]`)
	rEmptyStr := regexp.MustCompile(`[,{](\n)(\s)+"[^"]+":(\s)+""[,\n]`)
	rEmptyObj := regexp.MustCompile(`[,{](\n)(\s)+"[^"]+":(\s)+\{[\n\s]*\}[,\n]`)
	rEmptyArr := regexp.MustCompile(`[,{](\n)(\s)+"[^"]+":(\s)+\[[\n\s]*\][,\n]`)

	rTrim := []*regexp.Regexp{}
	if rmNull {
		rTrim = append(rTrim, rNull)
	}
	if rmEmptyStr {
		rTrim = append(rTrim, rEmptyStr)
	}
	if rmEmptyObj {
		rTrim = append(rTrim, rEmptyObj)
	}
	if rmEmptyArr {
		rTrim = append(rTrim, rEmptyArr)
	}

	trimmed := js
	for _, re := range rTrim {
	AGAIN:
		rm := false
		trimmed = re.ReplaceAllStringFunc(trimmed, func(s string) string {
			rm = true
			first, last := s[0], s[len(s)-1]
			switch {
			case first == '{' && last == ',': // first, NOT single
				return "{"
			case first == '{' && last == '\n': // first, single
				return "{"
			case first == ',' && last == ',': // NOT first, NOT last
				return ","
			case first == ',' && last == '\n': // NOT first, the last
				return "\n"
			default:
				log.Fatalln(s)
			}
			return ""
		})
		if rm {
			goto AGAIN
		}
	}

	return trimmed
}
