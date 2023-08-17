package jsontool

// Root :
// LIKE { "only-one-element": { ... } }
func Root(str string) string {
	root, _ := SglEleBlkCont(str)
	return root
}

// MkSglEleBlk :
// LIKE { "only-one-element": { ... } }
func MkSglEleBlk(name string, value any, fmt bool) string {
	// string type value to be added "quotes"
	switch vt := value.(type) {
	case string:
		if !(len(vt) >= 2 && (vt[0] == '{' || vt[0] == '[')) {
			value = fSf(`"%s"`, sTrim(vt, `"`))
		}
	case nil:
		value = "null"
	}

	str := fSf(`{ "%s": %v }`, name, value)
	failOnErrWhen(!IsValidStr(str), "%v", fEf("Error in Making JSON Block")) // test mode open
	if fmt {
		return FmtStr(str, "  ")
	}
	return str
}

// SglEleBlkCont :
// LIKE { "only-one-element": { ... } }
func SglEleBlkCont(str string) (string, string) {
	qtIdx1, qtIdx2 := -1, -1
	for i := 0; i < len(str); i++ {
		if qtIdx1 == -1 && str[i] == '"' {
			qtIdx1 = i
			continue
		}
		if qtIdx1 != -1 && str[i] == '"' {
			qtIdx2 = i
			break
		}
	}
	failOnErrWhen(str[qtIdx2+1] != ':', "%v", fEf("error (format) json"))
	failOnErrWhen(str[qtIdx2+2] != ' ', "%v", fEf("error (format) json"))
	ebIdx := sLastIndex(str, "}")
	return str[qtIdx1+1 : qtIdx2], sTrimRight(str[qtIdx2+3:ebIdx], " \t\n\r")
}

// SglEleAttrVal : attributes MUST be ahead of other sub-elements
func SglEleAttrVal(str, attr, attrPrefix string) (val string, ok bool) {
	lookFor := fSf(`%s%s`, attrPrefix, attr)
	dqGrp := []int{}
SCAN:
	for i := 0; i < len(str); i++ {
		switch str[i] {
		case '}':
			break SCAN
		case '"':
			dqGrp = append(dqGrp, i)
		}
	}
	dqV1, dqV2 := 0, 0
	for i := 0; i < len(dqGrp); i += 2 {
		dq1, dq2 := dqGrp[i], dqGrp[i+1]
		if str[dq1+1:dq2] == lookFor {
			dqV1, dqV2 = dqGrp[i+2], dqGrp[i+3]
			ok = true
			break
		}
	}
	if !ok {
		return "", ok
	}
	return str[dqV1+1 : dqV2], ok
}
