package jsontool

import (
	"strings"
)

// MergeSgl :
func MergeSgl(jsonGrp ...string) string {
	switch {
	case len(jsonGrp) >= 3:
		var sb strings.Builder
		for i, json := range jsonGrp {
			if i == 0 {
				p := sLastIndex(json, "}")
				sb.WriteString(sTrimRight(json[:p], " \t\r\n"))
			} else if i == len(jsonGrp)-1 {
				p := sIndex(json, "{")
				sb.WriteString(",")
				sb.WriteString(json[p+1:])
			} else {
				p1, p2 := sIndex(json, "{"), sLastIndex(json, "}")
				sb.WriteString(",")
				sb.WriteString(sTrimRight(json[p1+1:p2], " \t\r\n"))
			}
		}
		return sb.String()
	case len(jsonGrp) == 2:
		json1, json2 := jsonGrp[0], jsonGrp[1]
		p1 := sLastIndex(json1, "}")
		p2 := sIndex(json2, "{")
		return sTrimRight(json1[:p1], " \t\n") + "," + json2[p2+1:]
	case len(jsonGrp) == 1:
		return jsonGrp[0]
	}
	return ""
}

// merge4chan :
func merge4chan(chGrp ...<-chan string) string {
	var jsonGrp []string
	for _, ch := range chGrp {
		jsonGrp = append(jsonGrp, <-ch)
	}
	return MergeSgl(jsonGrp...)
}

// asyncScalarSel :
func asyncScalarSel(json, attr string) <-chan string {
	c := make(chan string)
	go func() {
		var sb strings.Builder
		sb.WriteString(fSf("{\n  \"%s\": [\n", attr))
		tag := fSf("\"%s\": ", attr)
		offset := len(tag)
		r := rxMustCompile(fSf(`%s.+,?\n`, tag))
		for _, l := range r.FindAllString(json, -1) {
			sb.WriteString("    ")
			l = sTrimRight(l, ",\r\n")[offset:]
			sb.WriteString(l)
			sb.WriteString(",\n")
		}
		sb.WriteString("  ]\n}")
		ret := sb.String()

		r = rxMustCompile(`,\n[ ]+\]`)
		pairs := r.FindAllStringIndex(ret, -1)
		failOnErrWhen(len(pairs) > 1, "%v", fEf("fatal!"))
		if len(pairs) == 1 {
			rmPos := pairs[0][0]
			ret = ret[:rmPos] + ret[rmPos+1:]
		}
		c <- ret
	}()
	return c
}

// ScalarSelX :
func ScalarSelX(json string, attrGrp ...string) string {
	chanGrp := make([]<-chan string, len(attrGrp))
	for i, attr := range attrGrp {
		chanGrp[i] = asyncScalarSel(json, attr)
	}
	return merge4chan(chanGrp...)
}

// ---------------------------------------------------- //

// L1Attrs : Level-1 attributes
// func L1Attrs(json string) (attrs []string) {
// 	failP1OnErrWhen(!isJSON(json), "%v", n3err.PARAM_INVALID_JSON)
// 	json = Fmt(json, 2)
// 	r := rxMustCompile(`\n  "[^"]+": [\[\{"-1234567890ntf]`)
// 	found := r.FindAllString(json, -1)
// 	for _, a := range found {
// 		attrs = append(attrs, a[4:len(a)-4])
// 	}
// 	return
// }

// Join :
// func Join(jsonL, fKey, jsonR, pKey, name string) (string, bool) {
// 	if name == "" {
// 		if hasAnySuffix(pKey, "-ID", "-id", "-Id", "_ID", "_id", "_Id") {
// 			name = pKey[:len(pKey)-3]
// 		}
// 	}

// 	inputs, keys, keyTypes := []string{jsonL, jsonR}, []string{fKey, pKey}, []string{"foreign", "primary"}
// 	starts, ends := []int{0, 0}, []int{0, 0}
// 	keyLines, keyValues := []string{"", ""}, []string{"", ""}
// 	posGrp := [][]int{}

// 	for i := 0; i < 2; i++ {
// 		lsAttr := toGeneralSlc(L1Attrs(inputs[i]))
// 		failOnErrWhen(!exist(keys[i], lsAttr...), "%v: NO %s key attribute [%s]", n3err.INTERNAL, keyTypes[i], keys[i])

// 		r := rxMustCompile(fSf(`\n  "%s": .+[,]?\n`, keys[i]))
// 		pSEs := r.FindAllStringIndex(inputs[i], 1)
// 		failOnErrWhen(len(pSEs) == 0, "%v: %s key's value error", n3err.INTERNAL, keyTypes[i])
// 		starts[i], ends[i] = pSEs[0][0], pSEs[0][1]
// 		keyLines[i] = sTrim(inputs[i][starts[i]:ends[i]], ", \t\r\n")
// 		keyValues[i] = keyLines[i][len(fKey)+4:]

// 		if i == 0 {
// 			posGrp = pSEs
// 			failOnErrWhen(exist(name, lsAttr...), "%v: [%s] already exists in left json", n3err.INTERNAL, name)
// 		}
// 	}

// 	if keyValues[0] == keyValues[1] {
// 		comma := ","
// 		if jsonL[posGrp[0][1]] == '}' {
// 			comma = ""
// 		}
// 		insert := fSf(`"%s": %s%s`, name, jsonR, comma)
// 		str := replByPosGrp(jsonL, posGrp, []string{insert})
// 		return Fmt(str, 2), true
// 	}

// 	return jsonL, false
// }

// ArrJoin :
// func ArrJoin(jsonarrL, fKey, jsonarrR, pKey, name string) (ret string, pairs [][2]int) {
// 	jsonLarr := SplitArr(jsonarrL, 2)
// 	jsonRarr := SplitArr(jsonarrR, 2)
// 	joined := []string{}
// 	for i, jsonL := range jsonLarr {
// 		for j, jsonR := range jsonRarr {
// 			if join, ok := Join(jsonL, fKey, jsonR, pKey, name); ok {
// 				// fPln(ok, i, j)
// 				pairs = append(pairs, [2]int{i, j})
// 				joined = append(joined, join)
// 			}
// 		}
// 	}
// 	return MakeArr(joined...), pairs
// }
