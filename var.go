package jsontool

import (
	"fmt"
	"regexp"
	"strings"

	lk "github.com/digisan/logkit"
)

var (
	fPln          = fmt.Println
	fSf           = fmt.Sprintf
	fEf           = fmt.Errorf
	sSplit        = strings.Split
	sJoin         = strings.Join
	sTrim         = strings.Trim
	sTrimLeft     = strings.TrimLeft
	sTrimRight    = strings.TrimRight
	sReplaceAll   = strings.ReplaceAll
	sIndex        = strings.Index
	sLastIndex    = strings.LastIndex
	sHasPrefix    = strings.HasPrefix
	sHasSuffix    = strings.HasSuffix
	rxMustCompile = regexp.MustCompile
	failOnErr     = lk.FailOnErr
	failOnErrWhen = lk.FailOnErrWhen
)

var (
	DEBUG = 0
)

// dropCR drops a terminal \r from the data.
var dropCR = func(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

func ZipInts(a1, a2 []int) (arr [][]int) {
	for i, a := range a1 {
		b := a2[i]
		arr = append(arr, []int{a, b})
	}
	return
}

func IndexAll(s, sub string) (starts, ends []int) {
	for i := sIndex(s, sub); i != -1; i = sIndex(s, sub) {
		last := 0
		if len(starts) > 0 {
			last = starts[len(starts)-1] + 1
		}
		start := last + i
		starts = append(starts, start)
		ends = append(ends, start+len(sub))
		s = s[i+1:]
	}
	return
}
