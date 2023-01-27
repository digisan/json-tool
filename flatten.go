package jsontool

import (
	"encoding/json"

	. "github.com/digisan/go-generics/v2"
)

func FlattenStr(jsonObj string) (map[string]any, error) {
	jsonMap := make(map[string]any)
	err := json.Unmarshal([]byte(jsonObj), &jsonMap)
	if err != nil {
		return nil, err
	}
	return MapNestedToFlat(jsonMap), nil
}

func Flatten(jsonObj []byte) (map[string]any, error) {
	jsonMap := make(map[string]any)
	err := json.Unmarshal(jsonObj, &jsonMap)
	if err != nil {
		return nil, err
	}
	return MapNestedToFlat(jsonMap), nil
}
