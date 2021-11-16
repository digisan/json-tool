package jsontool

import (
	"fmt"
	"log"
	"regexp"

	"github.com/digisan/gotk"
	"github.com/digisan/gotk/slice/ts"
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
			lvl = ts.MkSet(lvl...)
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

// 'sibling' is valid 'field' path sibling
func GetSiblingPath(field, sibling string, mLvlSiblings map[int][]string) (mFieldSibling map[string]string) {

	mFieldSibling = make(map[string]string)
	sPathsCandidates := []string{}
	for _, p := range GetFieldPaths(field, mLvlSiblings) {
		sPathsCandidates = append(sPathsCandidates, NewSibling(p, sibling))
	}
	const MAX_LEVEL = 1024
	for l := 0; l < MAX_LEVEL; l++ {
		for _, sib := range mLvlSiblings[l] {
			if ts.In(sib, sPathsCandidates...) {
				mFieldSibling[NewSibling(sib, field)] = sib
			}
		}
	}
	return
}

// 'siblings' are all valid path in one fixed 'field' path sibling
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

func GetAllLeafPaths(js string) (paths []string, values []gjson.Result) {
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

func GetLeafPathsOrderly(field string, paths []string) []string {
	rField := regexp.MustCompile(fSf(`\.%s(\.\d+)*$`, field))
	return ts.FM(paths, func(i int, e string) bool {
		return rField.MatchString(e) || field == e
	}, nil)
}
