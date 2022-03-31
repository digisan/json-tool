package jsontool

import (
	"encoding/json"
	"fmt"
)

func dumpMap(pk string, jv any, mflat *map[string]any) {

	switch m := jv.(type) {
	case float64, string, bool, nil:
		(*mflat)[pk] = m

	case []any:
		for i, a := range m {
			idx := fmt.Sprintf("%s.%d", pk, i)
			dumpMap(idx, a, mflat)
		}

	case map[string]any:
		{
			for k, v := range m {
				if pk != "" {
					k = fmt.Sprintf("%s.%s", pk, k)
				}

				switch mv := v.(type) {
				case []any:
					for i, a := range v.([]any) {
						idx := fmt.Sprintf("%s.%d", k, i)
						dumpMap(idx, a, mflat)
					}
				default:
					dumpMap(k, mv, mflat)
				}
			}
		}
	}
}

func FlattenStr(jsonObj string) (map[string]any, error) {
	jsonMap := make(map[string]any)
	err := json.Unmarshal([]byte(jsonObj), &jsonMap)
	if err != nil {
		return nil, err
	}
	flatMap := make(map[string]any)
	dumpMap("", jsonMap, &flatMap)
	return flatMap, nil
}

func Flatten(jsonObj []byte) (map[string]any, error) {
	jsonMap := make(map[string]any)
	err := json.Unmarshal(jsonObj, &jsonMap)
	if err != nil {
		return nil, err
	}
	flatMap := make(map[string]any)
	dumpMap("", jsonMap, &flatMap)
	return flatMap, nil
}
