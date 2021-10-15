package jsontool

import (
	"log"

	"github.com/digisan/gotk/slice/ts"
)

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
