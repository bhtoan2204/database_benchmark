package xtypes

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSON map[string]interface{}

func (m *JSON) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	return json.Unmarshal(bytes, m)
}

func isJSON(s []byte) bool {
	var js map[string]interface{}
	return json.Unmarshal(s, &js) == nil
}

func (m *JSON) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		*m = JSON{}
	}

	if string(b) == "[]" {
		*m = JSON{}
		return nil
	}

	b = bytes.TrimPrefix(b, []byte(`"`))
	b = bytes.TrimSuffix(b, []byte(`"`))

	if transformBytes := bytes.ReplaceAll(b, []byte(`\"`), []byte(`"`)); isJSON(transformBytes) {
		b = transformBytes
	}

	var tmp map[string]interface{}
	jd := json.NewDecoder(bytes.NewReader(b))
	if err := jd.Decode(&tmp); err != nil {
		return err
	}

	*m = JSON(tmp)

	return nil
}

func (m JSON) GetValuePath(items ...string) (interface{}, error) {
	value := m

	for idx, item := range items {
		nestedValue, ok := value[item]
		if !ok {
			return nil, fmt.Errorf("could not get field at %s", item)
		}

		if idx == len(items)-1 {
			return nestedValue, nil
		}

		nestedMap, ok := nestedValue.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("could not get nested item after field %s", item)
		}

		value = nestedMap
	}

	return nil, nil
}

func (m JSON) Value() (driver.Value, error) {
	b, err := json.Marshal(m)
	if len(b) == 0 {
		return "{}", nil
	}
	return string(b), err
}

func (m JSON) String() string {
	str, _ := json.Marshal(m)
	return string(str)
}

func Struct2JSON(obj interface{}) (JSON, error) {
	if obj == nil {
		return JSON{}, nil
	}

	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("convert struct to json failed=%w", err)
	}

	var result JSON
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return nil, fmt.Errorf("convert json data to JSON failed=%w", err)
	}

	return result, nil
}
