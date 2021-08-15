package jsontool

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"
)

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
	ScanResult struct {
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
func ScanObject(ctx context.Context, r io.Reader, mustarray, check bool, style OutStyle) (<-chan ScanResult, bool) {

	var (
		cOut = make(chan ScanResult)
		ja   = true
	)

	const (
		SCAN_STEP = bufio.MaxScanTokenSize
	)

	go func() {
		defer close(cOut)

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

		fillRst := func(object string) (next bool) {

			object = sTrimLeft(object, "[ \t")
			object = sTrimRight(object, ",] \t")
			rst := ScanResult{}

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

			select {
			case cOut <- rst:
			case <-ctx.Done():
				return false
			}

			return true
		}

		lineToRst := func(line string) (next bool) {

			// if partialLong, only inflate sbLine, return
			if partialLong {
				sbLine.WriteString(line)
				return true
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
				if !fillRst(sbObject.String()) {
					return false
				}
				sbObject.Reset()
			}

			for _, object := range objects {
				if !fillRst(object) {
					return false
				}
			}

			if len(nextHead) > 0 {
				sbObject.WriteString(nextHead)
				record = true
				return true
			}

			// object starts
			if N == 0 && L > 0 {
				record = true
			}

			if record {
				sbObject.WriteString(line)

				// object ends
				if L == 0 {
					if !fillRst(sbObject.String()) {
						return false
					}
					sbObject.Reset()
					record = false
				}
			}

			return true
		}

		split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {

			// DEBUG++
			// if DEBUG >= 1 {
			// 	fPln("why?")
			// }

			////////////////////////////////////////////////////////////////

			partialLong = false
			advance = bytes.IndexByte(data, '\n')

			switch {

			case atEOF && len(data) == 0:
				return 0, nil, nil

			case advance >= 0: // found
				return advance + 1, dropCR(data[:advance]), nil

			case advance == -1 && len(data) == cap(data) && cap(data) < SCAN_STEP: // didn't find, then expand to max cap
				return 0, nil, nil

			case advance == -1 && len(data) < cap(data) && !atEOF: // didn't find AND only take part this time
				return 0, nil, nil

			case advance == -1 && len(data) == SCAN_STEP: // didn't find, even if got max cap. ingest all
				partialLong = true
				return len(data), dropCR(data), nil

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
			if !lineToRst(line) {
				return
			}
		}
	}()

	time.Sleep(20 * time.Millisecond)
	return cOut, ja
}

func ScanObjectInArray(ctx context.Context, r io.Reader, check bool) (<-chan string, <-chan error, error) {
	var (
		cOut = make(chan string, 1)
		cErr = make(chan error, 1)
	)

	scanOut, ok := ScanObject(ctx, r, true, check, OUT_MIN)
	if !ok {
		return nil, nil, fmt.Errorf("not a valid JSON array")
	}
	go func() {
		for o := range scanOut {
			cOut <- o.Obj
			cErr <- o.Err
		}
		close(cOut)
		close(cErr)
	}()
	return cOut, cErr, nil
}
