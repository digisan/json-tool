package validator

import (
	"fmt"
	"os"
	"testing"

	dt "github.com/digisan/gotk/data-type"
	"github.com/tidwall/gjson"
)

func report(js string, fields ...string) {
	for _, field := range fields {
		fmt.Printf("===========> %v\n", field)
		r := gjson.Get(js, field)
		fmt.Printf("\n%s\n\n", r.Raw)
		if IsMissing(r) {
			fmt.Printf("'%s' is missing\n", field)
		}
		if IsNull(r) {
			fmt.Printf("'%s' is null\n", field)
		}
		if HasNotNullValue(r) {
			fmt.Printf("'%s' HasNotNullValue\n", field)
		}
		if HasEmptyValue(r) {
			fmt.Printf("'%s' HasEmptyValue\n", field)
		}
		if HasSomeValue(r) {
			fmt.Printf("'%s' HasSomeValue\n", field)
		}
		if IsBlankStr(r) {
			fmt.Printf("'%s' IsBlankStr\n", field)
		}
		if IsNullElemArray(r) {
			fmt.Printf("'%s' IsNullElemArray\n", field)
		}
		if IsEmptyElemArray(r) {
			fmt.Printf("'%s' IsEmptyElemArray\n", field)
		}
		fmt.Println()
	}
}

func TestTemp(t *testing.T) {

	data, err := os.ReadFile("./data.json")
	if err != nil {
		fmt.Println("reading json file error")
		return
	}
	js := string(data)

	if !dt.IsJSON(data) {
		fmt.Printf("not valid json")
		return
	}

	report(
		js,
		"Missing",
		"Null",
		"EmptyStr",
		"BlankStr",
		"EmptyObject",
		"EmptyArray",
		"NullElemArray",
		"EmptyStrElemArray",
		"EmptyArrElemArray",
		"EmptyObjElemArray",
	)

	// if IsNull(r) {
	// 	fmt.Printf("'%s' is null\n", field)
	// }

	// if IsEmptyStr(r) {
	// 	fmt.Printf("'%s' is empty string\n", field)
	// }

	// if IsHTMLStr(r) {
	// 	fmt.Printf("'%s' is html string\n", field)
	// }

	// if IsPlainStr(r) {
	// 	fmt.Printf("'%s' is plain string\n", field)
	// }

	// if IsStr(r) {
	// 	fmt.Printf("'%s' is string\n", field)
	// }

	// // fmt.Printf("not nil,", r.Num, r.Raw, r.Str, r.Type)
	// fmt.Println("----------------------------------------")

	// field = "OtherNames"
	// // field = "SIF"
	// r = gjson.Get(js, field)

	// if IsMissing(r) {
	// 	fmt.Printf("'%s' is missing\n", field)
	// 	return
	// }

	// if IsArr(r) {
	// 	fmt.Printf("'%s' is array\n", field)
	// }

	// if IsEmptyArr(r) {
	// 	fmt.Printf("'%s' is empty array\n", field)
	// }

	// if IsStrArr(r) {
	// 	fmt.Printf("'%s' is string array\n", field)
	// }

	// if IsObjArr(r) {
	// 	fmt.Printf("'%s' is object array\n", field)
	// }

	// r1 := gjson.Get(js, "SIF.#.Definition")
	// fmt.Printf("SIF.#.Definition: %+v\n", r1.Array())
}
