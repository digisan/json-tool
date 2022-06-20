package jsontool

import (
	"encoding/json"
	"fmt"
	"strings"

	. "github.com/digisan/go-generics/v2"
	"github.com/digisan/gotk/strs"
	lk "github.com/digisan/logkit"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func GetLvlChunkByPos(js string, stop int) (Lvl int, mLvlNChunk map[int]int) {
	var (
		pc     byte = 0
		quotes bool = false
	)
	mLvlNChunk = make(map[int]int)
	for i := 1; i <= 10; i++ {
		mLvlNChunk[i] = 0
	}
	for i := 0; i < len(js); i++ {
		if i > stop {
			break
		}
		c := js[i]
		switch {
		case c == '"' && pc != '\\':
			quotes = !quotes
		case c == '{' && !quotes:
			Lvl++
		case c == '}' && !quotes:
			mLvlNChunk[Lvl]++
			Lvl--
		}
		pc = c
	}
	return
}

func getPrefixChunkTrait(mLvlNChunk map[int]int, prefLvl int) string {
	sb := strings.Builder{}
	for i := 1; i <= prefLvl; i++ {
		sb.WriteString(fmt.Sprintf("%02d", mLvlNChunk[i]))
		if i < prefLvl {
			sb.WriteString("-")
		}
	}
	return sb.String()
}

func prefix2key(prefix string) string {
	return strs.SplitPartFromLast(strings.TrimSuffix(prefix, `":`), `"`, 1)
}

func markPrefix(prefix, trait, sp string, n int) (idPrefixGrp []string) {
	for i := 0; i < n; i++ {
		suffix := fmt.Sprintf("%s%s%s%02d", sp, trait, sp, i)
		prefix := strings.TrimSuffix(prefix, `":`)
		prefix = prefix + suffix + `":`
		idPrefixGrp = append(idPrefixGrp, prefix)
	}
	return
}

// [js] MUST be formatted with duplicate keys
// [prefix] MUST be like `\n    (n-space)"Field":`
func FixOneDupKeyOnce(js, prefix string) (string, bool) {

	starts, ends := strs.IndexAllByReg(js, prefix)

	if len(starts) == 0 {
		lk.Log("NO duplicate error at [%s]", prefix)
		return js, true
	}

	mTraitPosGrp := make(map[string][][2]int)
	for i, start := range starts {
		end := ends[i]
		lvl, mLvlNChunk := GetLvlChunkByPos(js, start)
		trait := getPrefixChunkTrait(mLvlNChunk, lvl)
		mTraitPosGrp[trait] = append(mTraitPosGrp[trait], [2]int{start, end})
	}

	N := 0
	for _, posGrp := range mTraitPosGrp {
		if len(posGrp) > 1 {
			N++
		}
	}
	lk.Warn("total duplicate error: %d", N)

	if N == 0 {
		lk.Log("NO duplicate error")
		return js, true
	}

	var (
		idPrefixes []string
		jsx        string
	)
	for trait, posGrp := range mTraitPosGrp {
		if n := len(posGrp); n > 1 {
			idPrefixes = markPrefix(prefix, trait, "^", n)
			jsx = strs.RangeReplace(js, posGrp, idPrefixes)
			break
		}
	}

	data, err := FlattenStr(jsx)
	lk.FailOnErr("%v", err)

	rmPaths := []string{}

	I := 0
	for path := range data {
		for _, idPrefix := range idPrefixes {
			idkey := prefix2key(idPrefix)
			if strings.Contains(path, idkey) {
				path = strs.SplitPart(path, idkey, 0) + idkey
				val := gjson.Get(jsx, path).Raw

				// fmt.Println(path, val)

				val = strings.TrimSpace(val)
				val = strings.TrimPrefix(val, "[")
				val = strings.TrimSuffix(val, "]")

				rmPaths = append(rmPaths, path)
				path = strs.SplitPart(path, "^", 0)
				jsx, _ = sjson.SetRaw(jsx, fmt.Sprintf("%s.%d", path, I), val)
				I++
				break
			}
		}
	}

	for _, path := range rmPaths {
		var err error
		jsx, err = sjson.Delete(jsx, path)
		lk.FailOnErr("%v", err)
	}

	return jsx, false
}

// [js] MUST be formatted with duplicate keys
// [prefix] MUST be like `\n    (n-space)"Field":`
func FixOneDupKey(js, prefix string) string {
	var (
		ok bool
		n  int
	)
	for !ok {
		js, ok = FixOneDupKeyOnce(js, prefix)
		// os.WriteFile(fmt.Sprintf("./debug_%02d.json", n), []byte(js), os.ModePerm)
		n++
	}
	return js
}

func RmDupEleOnce(js, path string) string {
	if r := gjson.Get(js, path); r.IsArray() {

		indices, array := []int{}, []string{}
		temp := []string{}

		for i, ele := range r.Array() {
			switch {
			case ele.IsObject():
				m := make(map[string]any)
				lk.FailOnErr("%v", json.Unmarshal([]byte(ele.Raw), &m))
				sm := fmt.Sprint(m)
				if NotIn(sm, temp...) {
					indices = append(indices, i)
					temp = append(temp, sm)
				}
			default:
				if NotIn(ele.Raw, temp...) {
					indices = append(indices, i)
					temp = append(temp, ele.Raw)
				}
			}
		}

		for i, ele := range r.Array() {
			if In(i, indices...) {
				array = append(array, ele.Raw)
			}
		}

		js, err := sjson.SetRaw(js, path, "["+strings.Join(array, ",")+"]")
		lk.FailOnErr("%v", err)
		return js
	}
	return js
}

// 'allPaths' from 'GetLeavesPathOrderly'
func RmDupEle(js, samplePath string, allPaths []string) string {
	if allPaths == nil {
		allPaths, _ = GetLeavesPathOrderly(js)
	}
	for _, path := range SimilarPaths(allPaths, samplePath) {
		js = RmDupEleOnce(js, path)
	}
	return js
}
