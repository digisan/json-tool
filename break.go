package jsontool

var (
	rName = rxMustCompile(`"[-\w]+":`)   // "name":
	rsVal = rxMustCompile(`^[-\d\.tfn]`) // non-string, simple value start
)

// BreakMulEleBlk : 'jsonstr' is LIKE {"1st-element": {...}, "2nd-element": {...}, "3rd-element": [...]}
// return one 'value' is like '{...}', OR like `[{...},{...},...]`
func BreakMulEleBlk(jsonstr string) (names, values []string) {
	jsonstr = sTrim(jsonstr, " \t\n")
	failOnErrWhen(jsonstr[0] != '{', "%v", fEf("error (format) json"))
	failOnErrWhen(jsonstr[len(jsonstr)-1] != '}', "%v", fEf("error (format) json"))

NEXT:
	if loc := rName.FindStringIndex(jsonstr); loc != nil { // find attr "name":
		s, e := loc[0], loc[1]
		root := jsonstr[s+1 : e-2]
		// fPln(root)
		names = append(names, root)
		jsonstr = sTrimLeft(jsonstr[e:], " ") // start @ "{" or "[" or simple...

		// Simple Non-String values
		if loc := rsVal.FindStringIndex(jsonstr); loc != nil {
			// fPln("non-string simple ele")
			for i := 1; i < len(jsonstr); i++ { // skip the 1st char
				c := jsonstr[i]
				if c == ',' || c == '\n' {
					values = append(values, jsonstr[:i])
					jsonstr = jsonstr[i+1:]
					goto NEXT
				}
			}
		}

		// Complex, Array or String value
		for i, mark := range []string{"{", "[", "\""} {
			if sHasPrefix(jsonstr, mark) {
				var m1, m2 byte
				switch i {
				case 0:
					m1, m2 = '{', '}'
				case 1:
					m1, m2 = '[', ']'
				default:
					m1, m2 = '"', '"'
				}
				L := 0
				for i := 0; i < len(jsonstr); i++ {
					c := jsonstr[i]
					if m1 != m2 { // Complex, Array
						if c == m1 { // { or [
							L++
						}
						if c == m2 { // } or ]
							L--
							if L == 0 {
								values = append(values, jsonstr[:i+1])
								jsonstr = jsonstr[i+1:]
								goto NEXT
							}
						}
					} else { // String
						if c == m1 { // "***"
							L++
							if L == 2 {
								// values = append(values, jsonstr[1:i]) // remove '"' at start&end (string & other types mixed)
								values = append(values, jsonstr[:i+1]) // remove '"' at start&end
								jsonstr = jsonstr[i+1:]
								goto NEXT
							}
						}
					}
				}
			}
		}
	}
	return
}

// BreakArr : 'jsonstr' is like [{...},{...}]
// i.e. [{...},{...}] => {...} AND {...}
// NO ele name could get here
func BreakArr(jsonstr string) (values []string, ok bool) {
	jsonstr = sTrim(jsonstr, " ")
	if jsonstr[0] != '[' || jsonstr[len(jsonstr)-1] != ']' {
		return values, false
	}
	L, S := 0, -1
	for i := 0; i < len(jsonstr); i++ {
		c := jsonstr[i]
		if c == '{' {
			L++
			if L == 1 {
				S = i
			}
		}
		if c == '}' {
			L--
			if L == 0 {
				values = append(values, jsonstr[S:i+1])
			}
		}
	}
	return values, true
}

// BreakMulEleBlkV2 : 'jsonstr' LIKE { "1st-element": {...}, "2nd-element": {...}, "3rd-element": [...] }
// in return 'values', array types are broken into duplicated names & its single value block
// one 'value' is like '{...}', 'names' may have duplicated names
func BreakMulEleBlkV2(jsonstr string) (names, values []string) {
	mIndEles := make(map[int][]string)
	Names, Values := BreakMulEleBlk(jsonstr)
	for i, Val := range Values {
		if elements, ok := BreakArr(Val); ok {
			mIndEles[i] = elements
		}
	}
	for i, Val := range Values {
		if elements, ok := mIndEles[i]; ok {
			for _, ele := range elements {
				names = append(names, Names[i])
				values = append(values, ele)
			}
		} else {
			names = append(names, Names[i])
			values = append(values, Val)
		}
	}
	return
}
