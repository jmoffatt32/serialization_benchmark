package jsonmethod

import (
	"encoding/json"
)

var NAME = "JSON"

func Encode(obj any) ([]byte, error) {

	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return b, nil

}

func Decode(data []byte, ptr any) error {

	err := json.Unmarshal(data, &ptr)
	if err != nil {
		return err
	}
	return nil
}
