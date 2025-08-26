package util

import (
	"github.com/goccy/go-json"
)

//nolint:ireturn
func JSONToStruct[T any](src any) (T, error) {
	var res T
	result, err := json.Marshal(src)
	if err != nil {
		return res, err
	}

	if err := json.Unmarshal(result, &res); err != nil {
		return res, err
	}

	return res, nil
}

//nolint:ireturn
func BytesToStruct[T any](data []byte) (T, error) {
	var res T
	if err := json.Unmarshal(data, &res); err != nil {
		return res, err
	}
	return res, nil
}
