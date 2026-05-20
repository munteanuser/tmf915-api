package repository

import "encoding/json"

func marshalJSON(v interface{}) (*string, error) {
	if v == nil {
		return nil, nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	s := string(b)
	return &s, nil
}

func unmarshalJSON(raw *string, v interface{}) error {
	if raw == nil || *raw == "" {
		return nil
	}
	return json.Unmarshal([]byte(*raw), v)
}
