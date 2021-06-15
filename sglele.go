package jsontool

// Root :
// LIKE { "only-one-element": { ... } }
func Root(jsonstr string) string {
	root, _ := SglEleBlkCont(jsonstr)
	return root
}

// MkSglEleBlk :
// LIKE { "only-one-element": { ... } }
func MkSglEleBlk(name string, value interface{}, fmt bool) string {
	// string type value to be added "quotes"
	switch vt := value.(type) {
	case string:
		if !(len(vt) >= 2 && (vt[0] == '{' || vt[0] == '[')) {
			value = fSf(`"%s"`, sTrim(vt, `"`))
		}
	case nil:
		value = "null"
	}

	jsonstr := fSf(`{ "%s": %v }`, name, value)
	failOnErrWhen(!IsValid(jsonstr), "%v", fEf("Error in Making JSON Block")) // test mode open
	if fmt {
		return Fmt(jsonstr, "  ")
	}
	return jsonstr
}

// SglEleBlkCont :
// LIKE { "only-one-element": { ... } }
func SglEleBlkCont(jsonstr string) (string, string) {
	qtIdx1, qtIdx2 := -1, -1
	for i := 0; i < len(jsonstr); i++ {
		if qtIdx1 == -1 && jsonstr[i] == '"' {
			qtIdx1 = i
			continue
		}
		if qtIdx1 != -1 && jsonstr[i] == '"' {
			qtIdx2 = i
			break
		}
	}
	failOnErrWhen(jsonstr[qtIdx2+1] != ':', "%v", fEf("error (format) json"))
	failOnErrWhen(jsonstr[qtIdx2+2] != ' ', "%v", fEf("error (format) json"))
	ebIdx := sLastIndex(jsonstr, "}")
	return jsonstr[qtIdx1+1 : qtIdx2], sTrimRight(jsonstr[qtIdx2+3:ebIdx], " \t\n\r")
}

// SglEleAttrVal : attributes MUST be ahead of other sub-elements
func SglEleAttrVal(jsonstr, attr, attrprefix string) (val string, ok bool) {
	lookfor := fSf(`%s%s`, attrprefix, attr)
	dqGrp := []int{}
SCAN:
	for i := 0; i < len(jsonstr); i++ {
		switch jsonstr[i] {
		case '}':
			break SCAN
		case '"':
			dqGrp = append(dqGrp, i)
		}
	}
	dqV1, dqV2 := 0, 0
	for i := 0; i < len(dqGrp); i += 2 {
		dq1, dq2 := dqGrp[i], dqGrp[i+1]
		if jsonstr[dq1+1:dq2] == lookfor {
			dqV1, dqV2 = dqGrp[i+2], dqGrp[i+3]
			ok = true
			break
		}
	}
	if !ok {
		return "", ok
	}
	return jsonstr[dqV1+1 : dqV2], ok
}
