package jsontool

import (
	"fmt"
	"os"
	"testing"

	"github.com/digisan/gotk/slice/ts"
)

func TestNewSibling(t *testing.T) {
	fmt.Println(NewSibling("a.b.c.d", "c"))
	fmt.Println(NewSibling("a", "c"))
	fmt.Println(NewSibling("", "c"))
	fmt.Println(NewSibling(".", "c"))
}

func TestNewUncle(t *testing.T) {
	fmt.Println(NewUncle("a.0.b.1.c.2.d", "CC"))
	fmt.Println(NewUncle("a.b.c.d", "e"))
	fmt.Println(NewUncle("a.b.c", "e"))
	fmt.Println(NewUncle("a", "c"))
	fmt.Println(NewUncle("..", "c"))
	fmt.Println(NewUncle(".", "c"))
}

func TestFieldName(t *testing.T) {
	name := FieldName("children.1.children.2.children.2.children.0.dcterms_title")
	fmt.Println(name)
	name = FieldName("children")
	fmt.Println(name)
}

func TestHasSiblings(t *testing.T) {
	data, err := os.ReadFile("./data/FlattenTest.json")
	if err != nil {
		panic(err)
	}
	js := string(data)
	mSibling, mFamilyTree := FamilyTree(js)

	ok := HasSiblings("array", mSibling, "array1", "array2")
	fmt.Println(ok)
	ok = HasSiblings("array.0.subarray.0.c", mSibling, "g", "f", "h")
	fmt.Println(ok)
	ok = HasSiblings("array.0.subarray.1.aa", mSibling, "cc", "ee", "aa")
	fmt.Println(ok)
	ok = HasSiblings("array.0.subarray", mSibling, "c", "a", "e")
	fmt.Println(ok)

	fmt.Println(PathExists("array.0.subarray.1.aa", mFamilyTree))
	fmt.Println(PathExists("array.0.subarray.2.aaa", mFamilyTree))
	fmt.Println(PathExists("object.aa", mFamilyTree))
}

func TestGetFieldPaths(t *testing.T) {
	data, err := os.ReadFile("./data/sif.json")
	if err != nil {
		panic(err)
	}
	js := string(data)
	mSibling, _ := FamilyTree(js)
	paths, _ := GetLeavesPathOrderly(js)

	lookfor := "value"

	fpaths1 := GetFieldPaths(lookfor, mSibling)
	fmt.Println(len(fpaths1))
	for _, p := range fpaths1 {
		fmt.Println(p)
	}

	fmt.Println("--------------------------------")

	fpaths2 := GetLeafPathsOrderly(lookfor, paths)
	fmt.Println(len(fpaths2))
	for _, p := range fpaths2 {
		fmt.Println(p)
	}

	fmt.Println("--------------------------------")

	fmt.Println(ts.Equal(fpaths1, fpaths2))

	fmt.Println(ts.Minus(fpaths1, fpaths2))

	fmt.Println(ts.Minus(fpaths2, fpaths1))

}

func TestConditionalMod(t *testing.T) {

	data, err := os.ReadFile("./data/FlattenTest.json")
	if err != nil {
		panic(err)
	}
	js := string(data)
	mSibling, _ := FamilyTree(js)

	// paths := GetFieldPaths("dcterms_title", mSibling) // get all paths which contains field 'dcterms_title'
	// fmt.Println(paths)

	// mFS := GetSiblingPath("dcterms_title", "asn_statementLabel", mSibling) // get all valid siblings for each 'dcterms_title' path
	// fmt.Println(mFS)

	mFSs := GetSiblingsPath("array", mSibling, "array1")
	for fp, sps := range mFSs {
		fmt.Println(fp, sps)
	}

	mFSs = GetSiblingsPath("c", mSibling, "e", "h")
	for fp, sps := range mFSs {
		fmt.Println(fp, sps)
	}

	// for k, v := range mFS {
	// 	js, _ = sjson.Set(js, k, 1000)                   // modify existing field 1
	// 	js, _ = sjson.Set(js, v, "2002-09-13")           // modify existing field 2
	// 	js, _ = sjson.Set(js, NewSibling(k, "HELLO"), 1) // create new field
	// }

	// os.WriteFile("./data/test_out.json", []byte(js), os.ModePerm)
}

func TestFamilyTree(t *testing.T) {

	data, err := os.ReadFile("./data/Activity.json")
	if err != nil {
		panic(err)
	}

	js := string(data)
	mSibling, mFT := FamilyTree(js)

	for l := 0; l < 1024; l++ {
		if len(mSibling[l]) > 0 {
			fmt.Println(mSibling[l])
		}
	}

	fmt.Println()

	for k, v := range mFT {
		if len(v) > 0 {
			fmt.Println(k, v)
		}
	}

	// for i := 0; i < 1024; i++ {
	// 	if fields := mSibling[i]; len(fields) > 0 {
	// 		// fmt.Println(fields)
	// 		for _, path := range fields {
	// 			if r := gjson.Get(js, path); r.Exists() {
	// 				value := r.String()

	// 				literal := regexp.QuoteMeta(value) // ***

	// 				re := regexp.MustCompile(literal)
	// 				locPairsAll := re.FindAllStringIndex(js, -1)
	// 				fmt.Println(path, locPairsAll)

	// 				// simple
	// 				field := pathLeaf(path)
	// 				re = regexp.MustCompile(fmt.Sprintf(`"%s":\s*"?%s"?`, field, literal))
	// 				locPairsSimple := re.FindAllStringIndex(js, -1)
	// 				for i := 0; i < len(locPairsSimple); i++ {
	// 					if js[locPairsSimple[i][1]-1] == '"' {
	// 						locPairsSimple[i][1] -= 1
	// 					}
	// 				}
	// 				fmt.Println(path, locPairsSimple)

	// 				// array
	// 				// re = regexp.MustCompile(fmt.Sprintf(`"%s":\s*%s`, field, literal))

	// 				fmt.Println()
	// 			}
	// 		}
	// 	}
	// }
}

func TestParentPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				path: "",
			},
			want: "",
		},
		{
			name: "",
			args: args{
				path: "a",
			},
			want: "",
		},
		{
			name: "",
			args: args{
				path: "a.b.c.d",
			},
			want: "a.b.c",
		},
		{
			name: "",
			args: args{
				path: "a.b.c.d.3.e",
			},
			want: "a.b.c.d",
		},
		{
			name: "",
			args: args{
				path: "a.b.c.2.d.3.e",
			},
			want: "a.b.c.2.d",
		},
		{
			name: "",
			args: args{
				path: "c.1.c.2.c.3.c.4.c.5.f",
			},
			want: "c.1.c.2.c.3.c.4.c",
		},
		{
			name: "",
			args: args{
				path: "c.1.c.2.c.3.c.4.c.5",
			},
			want: "c.1.c.2.c.3.c.4.c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParentPath(tt.args.path); got != tt.want {
				t.Errorf("ParentPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllLeafPaths(t *testing.T) {
	data, err := os.ReadFile("./data/FlattenTest.json")
	if err != nil {
		panic(err)
	}
	js := string(data)
	paths, values := GetLeavesPathOrderly(js)
	fmt.Println(len(paths))
	for i, p := range paths {
		v := values[i]
		fmt.Printf("%02d --- %v: %v\n", i, p, v.String())
	}
}

func TestOPath2TPath(t *testing.T) {
	type args struct {
		op  string
		sep string
	}
	tests := []struct {
		name    string
		args    args
		wantTp  string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				op:  "a",
				sep: ".",
			},
			wantTp:  "a",
			wantErr: false,
		},
		{
			name: "",
			args: args{
				op:  "a.b.1.c.3.d",
				sep: ".",
			},
			wantTp:  "a.b.c.d",
			wantErr: false,
		},
		{
			name: "",
			args: args{
				op:  "a.b.1.c.3.0.d",
				sep: ".",
			},
			wantTp:  "a.b.c.d",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTp, err := OPath2TPath(tt.args.op, tt.args.sep)
			if (err != nil) != tt.wantErr {
				t.Errorf("OPath2TPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTp != tt.wantTp {
				t.Errorf("OPath2TPath() = %v, want %v", gotTp, tt.wantTp)
			}
		})
	}
}

func TestLastSegMod(t *testing.T) {
	type args struct {
		path string
		sep  string
		f    func(last string) string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				path: "a.b.c.d",
				sep:  ".",
				f: func(last string) string {
					return "@" + last
				},
			},
			want: "a.b.c.@d",
		},
		{
			name: "",
			args: args{
				path: "a",
				sep:  ".",
				f: func(last string) string {
					return "@" + last
				},
			},
			want: "@a",
		},
		{
			name: "",
			args: args{
				path: "",
				sep:  ".",
				f: func(last string) string {
					return "@" + last
				},
			},
			want: "@",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LastSegMod(tt.args.path, tt.args.sep, tt.args.f); got != tt.want {
				t.Errorf("LastSegMod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetProperties(t *testing.T) {
	data, err := os.ReadFile("./data/complex.json")
	if err != nil {
		panic(err)
	}
	js := string(data)
	properties, locPseudo, mPropLocs, mPropValues := GetProperties(js)
	fmt.Println()
	fmt.Println(properties)
	fmt.Println()
	fmt.Println(locPseudo)
	fmt.Println()
	fmt.Println(mPropLocs["Object"])
	fmt.Println()
	// fmt.Println(mPropValues)
	// fmt.Println()

	fmt.Println(mPropValues["Object"][0])
	fmt.Println(mPropValues["Object"][1])
	fmt.Println(mPropValues["Object1"][0])

	fmt.Println("----------------------")

	fmt.Println(GetOutPropBlock(js, 2558))

	// fmt.Println(mPropValues["problems"][0])

	// blocks := blockByPropName(js, "Object")
	// for _, block := range blocks {
	// 	fmt.Println(block)
	// }
}

func TestGetOutPropBlockByProp(t *testing.T) {
	data, err := os.ReadFile("./data/complex.json")
	if err != nil {
		panic(err)
	}
	js := string(data)
	ops, obs := GetOutPropBlockByProp(js, "name")
	for i, op := range ops {
		ob := obs[i]
		fmt.Println(op)
		fmt.Println(ob)
	}
}

func TestRemoveParent(t *testing.T) {
	data, err := os.ReadFile("./data/complex.json")
	if err != nil {
		panic(err)
	}
	js := string(data)
	_, _, mPropLocs, mPropValues := GetProperties(js)
	js = RemoveParent(js, "Object", mPropLocs, mPropValues)
	fmt.Println(js)
	os.WriteFile("./data/complex_out.json", []byte(js), os.ModePerm)
}

func TestGetSiblingProps(t *testing.T) {
	data, err := os.ReadFile("./data/complex.json")
	if err != nil {
		panic(err)
	}
	js := string(data)

	mLvlSiblings, _ := FamilyTree(js)

	ops, obs := GetOutPropBlockByProp(js, "name")
	for i := range ops {
		ob := obs[i]
		fmt.Println(ob)
		siblings := GetSiblings(ob, "name")
		fmt.Println("siblings:", siblings)
		for _, s := range siblings {
			for k, v := range GetSiblingsPath("name", mLvlSiblings, s) {
				fmt.Println(k, v)
			}	
			fmt.Println()		
		}
		fmt.Println()
	}
}
