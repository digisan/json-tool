package jsontool

import (
	"fmt"
	"testing"
)

var jsonStr = `
	{
		"array1": [1,null,false,"abc"],
		"array2": [ [ 2, "foo", true, null ], [1,null,false,"abc", {"obj" : { "attr": "innervalue" }} ]],
	  	"array": [
		{
			"a": "b",
			"c": "d",
			"e": "f",
			"subarray": [
				{
					"aa": "BB",
					"cc": "DD",
					"ee": "FF"
				},
				{
					"aaa": "BBB",
					"ccc": "DDD",
					"eee": "FFF"
				}
			]
		},
		{
			"a": "B",
			"c": "D",
			"e": "F",
			"subarray": [
				{
					"foo": "FOO",
					"bar": "BAR",
					"baz": "BAZ"
				},
				{
					"foofoo": "FOOFOO",
					"barbar": "BARBAR",
					"bazbaz": "BAZBAZ"
				}
			]
		}
	  ],
	  "object1": {
		"object11": {
			"a": "b",
			"c": "d",
			"e": "f"
		},
		"object12": {
			"a": "B",
			"c": "D",
			"e": "F"
		}
		},
	  "boolean": true,
	  "null": null,
	  "number": 123,
	  "object": {
		"a": "b",
		"c": "d",
		"e": "f"
	  },
	  "string": "Hello World"
	}
	`

func TestFlattenObject(t *testing.T) {
	m, err := FlattenObject(jsonStr)
	fmt.Println(len(m), err)
	I := 0
	for k, v := range m {
		fmt.Printf("%d --- %v: %v\n", I, k, v)
		I++
	}
}
