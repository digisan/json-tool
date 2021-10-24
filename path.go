package jsontool

import (
	"log"

	"github.com/digisan/gotk/slice/ts"
)

// for json path sep by dot(.)
func ParentPath(path string) string {
	ss := sSplit(path, ".")
	return sJoin(ss[:len(ss)-1], ".")
}

func FieldName(path string) string {
	ss := sSplit(path, ".")
	return ss[len(ss)-1]
}

// NewSibling : return a new created sibling path
func NewSibling(fieldPath, sibName string) string {
	if pp := ParentPath(fieldPath); pp != "" {
		return ParentPath(fieldPath) + "." + sibName
	}
	return sibName
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

func GetFieldPaths(js, field string, mLvlSiblings map[int][]string) (paths []string) {

	if mLvlSiblings == nil {
		mLvlSiblings, _ = FamilyTree(js)
	}

	const MAX_LEVEL = 1024
	for l := 0; l < MAX_LEVEL; l++ {
		if len(mLvlSiblings[l]) > 0 {
			for _, sib := range mLvlSiblings[l] {
				switch {
				case sib == field:
					paths = append(paths, sib)
				case sHasSuffix(sib, "."+field):
					paths = append(paths, sib)
				default:
					// r := regexp.MustCompile(fmt.Sprintf(`.+\.%s\.\d+$`, field))
					// if r.MatchString(sib) {
					// 	paths = append(paths, sib)
					// }
				}
			}
		}
	}
	return
}

// 'sibling' is valid 'field' path sibling
func GetSiblingPath(js, field, sibling string, mLvlSiblings map[int][]string) (mFieldSibling map[string]string) {

	if mLvlSiblings == nil {
		mLvlSiblings, _ = FamilyTree(js)
	}

	mFieldSibling = make(map[string]string)
	sPathsCandidates := []string{}
	for _, p := range GetFieldPaths(js, field, mLvlSiblings) {
		if sContains(p, ".") {
			sPathsCandidates = append(sPathsCandidates, NewSibling(p, sibling))
		}
	}
	const MAX_LEVEL = 1024
	for l := 0; l < MAX_LEVEL; l++ {
		if len(mLvlSiblings[l]) > 0 {
			for _, sib := range mLvlSiblings[l] {
				if ts.In(sib, sPathsCandidates...) {
					mFieldSibling[NewSibling(sib, field)] = sib
				}
			}
		}
	}
	return
}

// 'siblings' are all valid path in one fixed 'field' path sibling
func GetSiblingsPath(js, field string, mLvlSiblings map[int][]string, siblings ...string) (mFieldSiblings map[string][]string) {

	if mLvlSiblings == nil {
		mLvlSiblings, _ = FamilyTree(js)
	}

	mFieldSiblingsCand := make(map[string][]string)
	for _, sib := range siblings {
		for fp, sp := range GetSiblingPath(js, field, sib, mLvlSiblings) {
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
func HasSiblings(js, fieldPath string, mLvlSiblings map[int][]string, siblings ...string) bool {
	mFSs := GetSiblingsPath(js, FieldName(fieldPath), mLvlSiblings, siblings...)
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

func PathExists(js, fieldPath string, mFamilyTree map[string][]string) bool {
	if mFamilyTree == nil {
		_, mFamilyTree = FamilyTree(js)
	}
	_, ok := mFamilyTree[fieldPath]
	return ok
}
