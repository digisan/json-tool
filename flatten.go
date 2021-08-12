package jsontool

import (
	"encoding/json"
	"fmt"
)

func dumpMap(pk string, jv interface{}, mflat *map[string]interface{}) {

	switch m := jv.(type) {
	case float64, string, bool, nil:
		(*mflat)[pk] = m

	case []interface{}:
		for i, a := range m {
			idx := fmt.Sprintf("%s.%d", pk, i)
			dumpMap(idx, a, mflat)
		}

	case map[string]interface{}:
		{
			for k, v := range m {
				if pk != "" {
					k = fmt.Sprintf("%s.%s", pk, k)
				}

				switch mv := v.(type) {
				case []interface{}:
					for i, a := range v.([]interface{}) {
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

func FlattenObject(jsonObj string) (map[string]interface{}, error) {
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonObj), &jsonMap)
	if err != nil {
		return nil, err
	}
	flatMap := make(map[string]interface{})
	dumpMap("", jsonMap, &flatMap)
	return flatMap, nil
}
