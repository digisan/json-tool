package jsontool

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/digisan/gotk"
	"github.com/digisan/go-generics/str"
	strs "github.com/digisan/gotk/strings"
	"github.com/tidwall/gjson"
)

func LastSegMod(path, sep string, f func(last string) string) string {
	ss := sSplit(path, sep)
	ss[len(ss)-1] = f(ss[len(ss)-1])
	return sJoin(ss, sep)
}

func OPath2TPath(op, sep string) (tp string, err error) {
	iNumGrp := []int{}
	ss := []string{}
	for i, s := range sSplit(op, sep) {
		if !gotk.IsNumeric(s) {
			ss = append(ss, s)
		} else {
			iNumGrp = append(iNumGrp, i)
		}
	}
	if len(iNumGrp) > 1 {
		for i := 1; i < len(iNumGrp); i++ {
			prev, curr := iNumGrp[i-1], iNumGrp[i]
			if curr-prev == 1 {
				err = fmt.Errorf("array as another array's element cannot be converted to TypePath")
			}
		}
	}
	return sJoin(ss, sep), err
}

// for json path sep by dot(.)
func ParentPath(path string) string {
	ss := sSplit(path, ".")
	if len(ss) >= 2 {
		if gotk.IsNumeric(ss[len(ss)-2]) {
			return sJoin(ss[:len(ss)-2], ".")
		}
	}
	return sJoin(ss[:len(ss)-1], ".")
}

func FieldName(path string) string {
	ss := sSplit(path, ".")
	return ss[len(ss)-1]
}

func NewChild(fieldPath, childName string) string {
	return fieldPath + "." + childName
}

// NewSibling : return a new created sibling path,
// empty fieldPath return empty,
// "." fieldPath creates a new field as sibName
func NewSibling(fieldPath, sibName string) string {
	if fieldPath == "" {
		return ""
	}
	ss := sSplit(fieldPath, ".")
	sibPath := sJoin(ss[:len(ss)-1], ".") + "." + sibName
	return sTrimLeft(sibPath, ".")
}

// NewUncle : return a new created uncle path
// empty fieldPath return empty,
// ".." fieldPath creates a new field as uncleName
func NewUncle(fieldPath, uncleName string) string {
	if fieldPath == "" {
		return ""
	}
	pp := ParentPath(fieldPath)
	if pp == "" {
		return ""
	}
	return NewSibling(pp, uncleName)
}

func FamilyTree(js string) (mLvlSiblings map[int][]string, mFamilyTree map[string][]string) {

	const MAX_LEVEL = 1024

	mData, err := FlattenStr(js)
	if err != nil {
		log.Fatalln(err)
	}

	lvls := make([][]string, MAX_LEVEL)
	for path := range mData {
		ss := sSplit(path, ".")
		for i := 0; i < len(ss); i++ {
			lvls[i] = append(lvls[i], sJoin(ss[:i+1], "."))
		}
	}

	mLvlSiblings = make(map[int][]string)

	for I, lvl := range lvls {
		if len(lvl) > 0 {
			lvl = str.MkSet(lvl...)
			mLvlSiblings[I] = lvl
		}
	}

	mFamilyTree = make(map[string][]string)

	for i := 0; i < MAX_LEVEL; i++ {
		this := mLvlSiblings[i]
		for _, f := range this {
			mFamilyTree[f] = []string{}
			if end := sLastIndex(f, "."); end > 0 {
				parent := f[:end]
				if _, ok := mFamilyTree[parent]; ok {
					mFamilyTree[parent] = append(mFamilyTree[parent], f)
				}
			}
		}
	}

	return
}

func GetFieldPaths(field string, mLvlSiblings map[int][]string) (paths []string) {
	const MAX_LEVEL = 1024
	rField := regexp.MustCompile(fSf(`\.%s(\.\d+)*$`, field))
	for l := 0; l < MAX_LEVEL; l++ {
		if len(mLvlSiblings[l]) > 0 {
			for _, sib := range mLvlSiblings[l] {
				if rField.MatchString(sib) || sib == field {
					paths = append(paths, sib)
				}
			}
		}
	}
	return
}

// 'sibling' is valid 'field' path sibling, return m[each/path/field]"valid/path/sibling"
func GetSiblingPath(field, sibling string, mLvlSiblings map[int][]string) (mFieldSibling map[string]string) {

	mFieldSibling = make(map[string]string)
	sPathsCandidates := []string{}
	for _, p := range GetFieldPaths(field, mLvlSiblings) {
		sPathsCandidates = append(sPathsCandidates, NewSibling(p, sibling))
	}
	const MAX_LEVEL = 1024
	for l := 0; l < MAX_LEVEL; l++ {
		for _, sib := range mLvlSiblings[l] {
			if str.In(sib, sPathsCandidates...) {
				mFieldSibling[NewSibling(sib, field)] = sib
			}
		}
	}
	return
}

// 'siblings' are all valid path in one fixed 'field' path sibling, return m[one fixed/path/field]["valid/path/sibling1", "valid/path/sibling2",...]
func GetSiblingsPath(field string, mLvlSiblings map[int][]string, siblings ...string) (mFieldSiblings map[string][]string) {

	mFieldSiblingsCand := make(map[string][]string)
	for _, sib := range siblings {
		for fp, sp := range GetSiblingPath(field, sib, mLvlSiblings) {
			mFieldSiblingsCand[fp] = append(mFieldSiblingsCand[fp], sp)
		}
	}

	mFieldSiblings = make(map[string][]string)
	for fp, sps := range mFieldSiblingsCand {
		if len(sps) == len(siblings) {
			mFieldSiblings[fp] = sps
		}
	}

	return
}

// if return 'false', with the first uncovered sibling name
func HasSiblings(fieldPath string, mLvlSiblings map[int][]string, siblings ...string) bool {
	mFSs := GetSiblingsPath(FieldName(fieldPath), mLvlSiblings, siblings...)
	if len(mFSs) == 0 {
		return false
	}

	sibpaths := mFSs[fieldPath]
NEXT:
	for _, sib := range siblings {
		for _, sibpath := range sibpaths {
			if sib == FieldName(sibpath) {
				continue NEXT
			}
		}
		return false
	}
	return true
}

func PathExists(fieldPath string, mFamilyTree map[string][]string) bool {
	_, ok := mFamilyTree[fieldPath]
	return ok
}

func GetLeavesPathOrderly(js string) (paths []string, values []gjson.Result) {
	iteratePath(js, "", true, false, &paths, &values)
	return
}

func iteratePath(js, ppath string, first, array bool, paths *[]string, values *[]gjson.Result) {

	path := ""
	idx := 0

	gjson.Get(js, "@this").ForEach(func(key, value gjson.Result) bool {

		kstr := key.String()

		if first {
			path = kstr
		} else {
			if kstr == "" {
				if array {
					path = fSf(`%s.%d`, ppath, idx)
					idx++
				} else {
					path = ppath
				}
			} else {
				path = fSf(`%s.%s`, ppath, kstr)
			}
		}

		switch {
		case value.IsArray():
			for i, ele := range value.Array() {
				elestr := ele.Raw
				ipath := fSf("%s.%d", path, i)
				iteratePath(elestr, ipath, false, elestr[0] == '[', paths, values)
			}
		case value.IsObject():
			iteratePath(value.String(), path, false, false, paths, values)
		default:
			// fmt.Println(path, value)
			*paths = append(*paths, path)
			*values = append(*values, value)
		}
		return true
	})
}

func GetLeafPathsOrderly(field string, allPaths []string) []string {
	rField := regexp.MustCompile(fSf(`\.%s(\.\d+)*$`, field))
	return str.FM(allPaths, func(i int, e string) bool {
		return rField.MatchString(e) || field == e
	}, nil)
}

////////////////////////////////////////////////////////////////////////////

func getNearPos4OA(js string, start int) (start4val, end int) {
	var open, close byte = '*', '*'
	inDQ := false
	n := 0
	s4c := -1
	for p := start; p < len(js); p++ {
		c := js[p]
		if !inDQ && c == '"' && js[p-1] != '\\' {
			inDQ = true
			continue
		} else if inDQ && c == '"' && js[p-1] != '\\' {
			inDQ = false
			continue
		}
		if !inDQ {
			switch c {

			case '{', '[':
				if open == '*' {
					open = c
				}
				if close == '*' {
					if open == '{' {
						close = '}'
					} else {
						close = ']'
					}
				}
				if c == open {
					n++
				}
				if s4c == -1 {
					s4c = p
				}

			case close:
				n--
				if n == 0 {
					return s4c, p + 1
				}
			}
		}
	}
	return s4c, -1
}

func GetOutPropBlock(js string, start4prop int) (prop, block string) {

	start, end := -1, -1
	inDQ := false
	n := 0

	if js[start4prop] == '"' {
		start4prop--
	}

FOR_PREV_OPEN_BRACKET:
	for p := start4prop; p >= 0; p-- {
		c := js[p]
		if !inDQ && c == '"' && js[p-1] != '\\' {
			inDQ = true
			continue
		} else if inDQ && c == '"' && js[p-1] != '\\' {
			inDQ = false
			continue
		}
		if !inDQ {
			switch c {
			case '{':
				n++
				if n == 1 {
					start = p
					break FOR_PREV_OPEN_BRACKET
				}
			case '}':
				n--
			}
		}
	}

	start, end = getNearPos4OA(js, start)
	block = js[start:end]

	inDQ = false
	for p := start; p >= 0; p-- {
		c := js[p]
		if !inDQ && c == '"' && js[p-1] != '\\' {
			inDQ = true
			end = p
		} else if inDQ && c == '"' && js[p-1] != '\\' {
			start = p
			break
		}
	}
	prop = js[start+1 : end] // without '" "'

	return
}

func GetOutPropBlockByProp(js, prop string) (props, blocks []string) {
	_, _, mPropLocs, _ := GetProperties(js)
	for _, loc := range mPropLocs[prop] {
		prop, block := GetOutPropBlock(js, loc[0])
		props = append(props, prop)
		blocks = append(blocks, block)
	}
	return
}

func GetProperties(js string) (
	properties []string,
	loc [][2]int,
	mPropLocs map[string][][3]int, // start, start4value, end
	mPropValues map[string][]interface{},
) {

	var (
		rProperty = regexp.MustCompile(`"[^"]*"\s*:\s*[\{\["\-\dtfn]`)
		rKVstr    = regexp.MustCompile(`^"[^"]*"\s*:\s*"[^"]*"\s*,?`)
		rKVnum    = regexp.MustCompile(`^"[^"]*"\s*:\s*\-?\d+\.?\d*\s*,?`)
		rKVbool   = regexp.MustCompile(`^"[^"]*"\s*:\s*[(true)|(false)]\s*,?`)
		rKVnull   = regexp.MustCompile(`^"[^"]*"\s*:\s*null\s*,?`)
		rKVobj    = regexp.MustCompile(`^"[^"]*"\s*:\s*\{`)
		rKVarr    = regexp.MustCompile(`^"[^"]*"\s*:\s*\[`)
	)

	Idx := 0
	mPropIdx := make(map[string][]int)

	rProperty.ReplaceAllStringFunc(js, func(s string) string {
		s = s[:sLastIndex(s, ":")]
		s = sTrimPrefix(s, "\"")
		s = sTrimRight(s, " \t")
		s = sTrimSuffix(s, "\"")
		properties = append(properties, s)
		mPropIdx[s] = append(mPropIdx[s], Idx)
		Idx++
		return ""
	})

	mPropLocs = make(map[string][][3]int)
	for i, p := range rProperty.FindAllStringIndex(js, -1) {
		loc = append(loc, [2]int{p[0], -1})

		for prop, indices := range mPropIdx {
			for _, idx := range indices {
				if i == idx {
					mPropLocs[prop] = append(mPropLocs[prop], [3]int{p[0], -1, -1})
				}
			}
		}
	}

	mPropValues = make(map[string][]interface{})
NEXT_PROP:
	for prop, locs := range mPropLocs {
		for _, loc := range locs {
			temp := js[loc[0]:]

			if s := rKVstr.FindString(temp); s != "" {

				ps, _ := strs.IndexAll(s, "\"")
				mPropValues[prop] = append(mPropValues[prop], s[ps[2]+1:ps[3]])

			} else if s := rKVnum.FindString(temp); s != "" {

				s = sTrimSuffix(s, ",")
				p := sLastIndex(s, ":")
				numstr := sTrim(s[p+1:], " \t\n")
				num, err := strconv.ParseFloat(numstr, 64)
				if err != nil {
					panic(err)
				}
				mPropValues[prop] = append(mPropValues[prop], num)

			} else if s := rKVbool.FindString(temp); s != "" {

				s = sTrimSuffix(s, ",")
				p := sLastIndex(s, ":")
				boolstr := sTrim(s[p+1:], " \t\n")
				b, err := strconv.ParseBool(boolstr)
				if err != nil {
					panic(err)
				}
				mPropValues[prop] = append(mPropValues[prop], b)

			} else if s := rKVnull.FindString(temp); s != "" {

				mPropValues[prop] = append(mPropValues[prop], nil)

			} else if rKVobj.MatchString(temp) {

				mPropValues[prop] = append(mPropValues[prop], "#OBJECT")

			} else if rKVarr.MatchString(temp) {

				mPropValues[prop] = append(mPropValues[prop], "#ARRAY")

			} else {

				continue NEXT_PROP

			}
		}
	}

	////////////////////////////////////////////////////////////////

	// update end position for [object] & [array]
	for prop, vals := range mPropValues {

		// if prop == "labs" {
		// 	fmt.Println("DEBUG")
		// }

		for i, val := range vals {
			switch sv := val.(type) {
			case string:
				if str.In(sv, "#OBJECT", "#ARRAY") {

					start := mPropLocs[prop][i][0]
					s4c, end := getNearPos4OA(js, start)
					mPropLocs[prop][i][1] = s4c
					mPropLocs[prop][i][2] = end

				} else {
					// simple string
				}
			default: // other type
			}
		}
	}

	// update value content for [object] & [array]
	for prop, locs := range mPropLocs {
		for i, loc := range locs {
			s4c, e := loc[1], loc[2]
			if e != -1 {
				mPropValues[prop][i] = js[s4c:e]
			}
		}
	}

	return
}

func RemoveParent(js, prop string, mPropLocs map[string][][3]int, mPropValues map[string][]interface{}) string {
	locs := [][2]int{}
	for _, loc := range mPropLocs[prop] {
		locs = append(locs, [2]int{loc[0], loc[2]})
	}
	vals := []string{}
	for _, val := range mPropValues[prop] {
		sval := val.(string)
		sval = sTrim(sval, "{}")
		vals = append(vals, sval)
	}
	return strs.RangeReplace(js, locs, vals)
}

func GetSiblings(block, prop string) (siblings []string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(block), &m); err == nil {
		for k := range m {
			if k != prop {
				siblings = append(siblings, k)
			}
		}
	} else {
		log.Fatalln(err)
	}
	return
}
