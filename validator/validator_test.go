package validator

import (
	"fmt"
	"testing"

	dt "github.com/digisan/gotk/data-type"
	"github.com/tidwall/gjson"
)

const js = `{
	"Entity": "<p>abc</p>",
	"OtherNames": [
	  "DEF",
	  "SES"
	],
	"Definition": "",
	"SIF": [
	  {
		"XPath": [],
		"Definition": "abc",
		"Commentary": "",
		"Datestamp": ""
	  },
	  {
		"XPath": [],
		"Definition": "def",
		"Commentary": "",
		"Datestamp": ""
	  }
	]
}`

func TestTemp(t *testing.T) {

	if !dt.IsJSON([]byte(js)) {
		fmt.Printf("not valid json")
		return
	}

	field := "Entity"
	r := gjson.Get(js, field)

	if IsMissing(r) {
		fmt.Printf("'%s' is missing\n", field)
		return
	}

	if IsNull(r) {
		fmt.Printf("'%s' is null\n", field)
	}

	if IsEmptyStr(r) {
		fmt.Printf("'%s' is empty string\n", field)
	}

	if IsHTMLStr(r) {
		fmt.Printf("'%s' is html string\n", field)
	}

	if IsPlainStr(r) {
		fmt.Printf("'%s' is plain string\n", field)
	}

	if IsStr(r) {
		fmt.Printf("'%s' is string\n", field)
	}

	// fmt.Printf("not nil,", r.Num, r.Raw, r.Str, r.Type)
	fmt.Println("----------------------------------------")

	field = "OtherNames"
	// field = "SIF"
	r = gjson.Get(js, field)

	if IsMissing(r) {
		fmt.Printf("'%s' is missing\n", field)
		return
	}

	if IsArr(r) {
		fmt.Printf("'%s' is array\n", field)
	}

	if IsEmptyArr(r) {
		fmt.Printf("'%s' is empty array\n", field)
	}

	if IsStrArr(r) {
		fmt.Printf("'%s' is string array\n", field)
	}

	if IsObjArr(r) {
		fmt.Printf("'%s' is object array\n", field)
	}

	r1 := gjson.Get(js, "SIF.#.Definition")
	fmt.Printf("SIF.#.Definition: %+v\n", r1.Array())
}
