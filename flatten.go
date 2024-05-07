package jsontool

import (
	"encoding/json"

	. "github.com/digisan/go-generics"
)

// empty object i.e. {} will be ignored, excluded from returned map
func FlattenStr(jsonObj string) (map[string]any, error) {
	jsonMap := make(map[string]any)
	err := json.Unmarshal([]byte(jsonObj), &jsonMap)
	if err != nil {
		return nil, err
	}
	return MapNestedToFlat(jsonMap), nil
}

// empty object i.e. {} will be ignored, excluded from returned map
func Flatten(jsonObj []byte) (map[string]any, error) {
	jsonMap := make(map[string]any)
	err := json.Unmarshal(jsonObj, &jsonMap)
	if err != nil {
		return nil, err
	}
	return MapNestedToFlat(jsonMap), nil
}
