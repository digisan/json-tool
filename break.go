package jsontool

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"time"
)

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

// return after processing, Level & prev-obj tail & next-obj head & inline objects
func analyseJL(line string, L int) (Lout int, prevTail, nextHead string, objects []string) {

	var pc byte = 0
	quotes := false
	s, e := -1, -1
	gotPrevTail := false

	for i := 0; i < len(line); i++ {
		c := line[i]
		switch {
		case c == '"' && pc != '\\':
			quotes = !quotes
		case c == '{' && !quotes:
			L++
			if L == 1 {
				s, e = i, -1

				if !gotPrevTail {
					prevTail = sTrimRight(line[:i], "[, \t")
					gotPrevTail = true
				}
			}
		case c == '}' && !quotes:
			L--
			if L == 0 {
				e = i

				nextHead = sTrimLeft(line[i+1:], "], \t")
			}
		}
		pc = c

		// if got object in single line
		if s > -1 && e > s {
			objects = append(objects, line[s:e+1])
			s, e = -1, -1
		}
	}

	return L, prevTail, nextHead, objects
}

// detect left-curly-bracket '{', '{'->count++, '}'->count--
// func detectLCB(line string) (L int, objects []string) {

// 	var pc byte = 0
// 	quotes := false
// 	s, e := -1, -1

// 	for i := 0; i < len(line); i++ {
// 		c := line[i]
// 		switch {
// 		case c == '"' && pc != '\\':
// 			quotes = !quotes
// 		case c == '{' && !quotes:
// 			L++
// 			if L == 1 {
// 				s, e = i, -1
// 			}
// 		case c == '}' && !quotes:
// 			L--
// 			if L == 0 {
// 				e = i
// 			}
// 		}
// 		pc = c

// 		// if got object in single line
// 		if s > -1 && e > s {
// 			objects = append(objects, line[s:e+1])
// 			s, e = -1, -1
// 		}
// 	}
// 	return
// }

type (
	ResultOfScan struct {
		Obj string
		Err error
	}
	OutStyle int
)

const (
	OUT_ORI OutStyle = 0
	OUT_FMT OutStyle = 1
	OUT_MIN OutStyle = 2
)

// ScanObject : any format json array should be OK.
func ScanObject(r io.Reader, mustarray, check bool, style OutStyle) (<-chan ResultOfScan, bool) {

	var (
		chRst = make(chan ResultOfScan)
		ja    = true
	)

	go func() {
		defer func() {
			close(chRst)
		}()

		const (
			SCAN_STEP = bufio.MaxScanTokenSize
		)

		var (
			lbbChecked  = false
			N           = 0
			record      = false
			sbObject    = &strings.Builder{}
			partialLong = false
			sbLine      = &strings.Builder{}
			scanner     = bufio.NewScanner(r)
			scanBuf     = make([]byte, SCAN_STEP)
		)

		fillRst := func(object string) {

			object = sTrimLeft(object, "[ \t")
			object = sTrimRight(object, ",] \t")
			rst := ResultOfScan{}

			// if invalid json, report to error
			if check && !IsValid(object) {
				rst.Err = fEf("Error JSON @ \n%v\n", object)
			}

			// only record valid json
			if rst.Err == nil {
				switch style {
				case OUT_ORI:
					break
				case OUT_FMT:
					object = Fmt(object, "  ")
				case OUT_MIN:
					object = Minimize(object)
				}
				rst.Obj = object
			}

			chRst <- rst
		}

		lineToRst := func(line string) {

			// if partialLong, only inflate sbLine, return
			if partialLong {
				sbLine.WriteString(line)
				return
			}

			// if not partialLong, and sbLine has content, modify input line
			if sbLine.Len() > 0 {
				line = sbLine.String() + line
				defer sbLine.Reset()
			}

			L, prevTail, nextHead, objects := analyseJL(line, N)
			defer func() { N = L }()

			if len(prevTail) > 0 {
				sbObject.WriteString(prevTail)
				fillRst(sbObject.String())
				sbObject.Reset()
			}

			for _, object := range objects {
				fillRst(object)
			}

			if len(nextHead) > 0 {
				sbObject.WriteString(nextHead)
				record = true
				return
			}

			// object starts
			if N == 0 && L > 0 {
				record = true
			}

			if record {
				sbObject.WriteString(line)

				// object ends
				if L == 0 {
					fillRst(sbObject.String())
					sbObject.Reset()
					record = false
				}
			}
		}

		split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {

			// DEBUG++
			// if DEBUG >= 1 {
			// 	fPln("why?")
			// }

			////////////////////////////////////////////////////////////////

			// if atEOF && len(data) == 0 {
			// 	return 0, nil, nil
			// }
			// if i := bytes.IndexByte(data, '\n'); i >= 0 {
			// 	// We have a full newline-terminated line.
			// 	return i + 1, dropCR(data[0:i]), nil
			// }
			// // If we're at EOF, we have a final, non-terminated line. Return it.
			// if atEOF {
			// 	return len(data), dropCR(data), nil
			// }
			// // Request more data.
			// return 0, nil, nil

			////////////////////////////////////////////////////////////////

			partialLong = false
			advance = bytes.IndexByte(data, '\n')

			switch {

			case atEOF && len(data) == 0:
				return 0, nil, nil

			case advance >= 0: // found
				return advance + 1, dropCR(data[:advance]), nil

			case advance == -1 && cap(data) < SCAN_STEP: // didn't find, then expand to max cap
				return 0, nil, nil

			case advance == -1 && len(data) == SCAN_STEP: // didn't find, even if got max cap. ingest all
				partialLong = true
				return SCAN_STEP, dropCR(data), nil

			default: // case advance == -1 && len(data) < SCAN_STEP: // didn't find, got part when at max cap. ingest & close long line.
				return len(data), dropCR(data), nil
			}
		}

		scanner.Buffer(scanBuf, SCAN_STEP)
		scanner.Split(split)

		for scanner.Scan() {
			line := scanner.Text()
			if !lbbChecked {
				if s := sTrim(line, " \t"); len(s) > 0 {
					if s[0] != '[' {
						ja = false
						if mustarray {
							return // if not json array, do not ingest
						}
					}
					lbbChecked = true
				}
			}
			lineToRst(line)
		}
	}()

	time.Sleep(5 * time.Millisecond)
	return chRst, ja
}
