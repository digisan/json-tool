package jsontool

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestTransform(t *testing.T) {

	RegisterRule(`^array\.`, func(path string, value interface{}) (ok bool, ps []string, vs []interface{}) {
		ok = true
		ps = append(ps, "")
		vs = append(vs, nil)
		return
	})

	RegisterRule(`.*`, func(path string, value interface{}) (ok bool, ps []string, vs []interface{}) {
		ok = GetStrVal(value) == "H"
		ps = append(ps, "")
		vs = append(vs, nil)
		return
	})

	data, err := os.ReadFile("./data/FlattenTest.json")
	if err != nil {
		panic(err)
	}
	fmt.Println(FmtStr(Transform(data), "\t"))
}

func TestComposite(t *testing.T) {
	data, err := os.ReadFile("./data/FlattenTest.json")
	if err != nil {
		panic(err)
	}
	m, err := Flatten(data)
	fmt.Println(len(m), err)

	str := Composite(m, nil)
	os.WriteFile("./data/FlattenTest_Composite_1.json", []byte(str), os.ModePerm)

	// str = Composite(m,
	// 	func(path string, value interface{}) (string, interface{}, bool) {
	// 		if strings.Contains(path, "object") {
	// 			if value == "b" {
	// 				return NewSibling(path, "ABC"), value.(string) + "ABC", false
	// 			}
	// 			return path, value, false
	// 		}
	// 		return "", nil, false
	// 	})
	// fmt.Println(str)

	str = Composite2(m, func(path string, value interface{}) (p []string, v []interface{}) {
		if strings.Contains(path, "object") {
			if value == "b" || value == "F" {
				return []string{
						NewChild(NewSibling(path, "DEF"), "abc"),
						NewChild(NewSibling(path, "DEF"), "def"),
					},
					[]interface{}{
						value.(string) + "ABC",
						value.(string) + "DEF",
					}
			}
			return []string{path}, []interface{}{value}
		}
		return nil, nil
	})
	fmt.Println(FmtStr(str, "  "))

	// str = CompositeIncl(m, "array.0.c", "object.a", "object1.object11.a")
	// os.WriteFile("./data/FlattenTest_CompositeIncl.json", []byte(str), os.ModePerm)

	// str = CompositeExcl(m, "array.0.c", "object.a", "object1.object11.a")
	// os.WriteFile("./data/FlattenTest_CompositeExcl.json", []byte(str), os.ModePerm)
}
