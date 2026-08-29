package utils

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
)

var notPointer = errors.New("must be a pointer")

func WriteJSON(path string, val any) error {
	if !isPointer(val) {
		return notPointer
	}

	data, err := json.MarshalIndent(val, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func ReadJSON(path string, data any) error {
	if !isPointer(data) {
		return notPointer
	}

	fileData, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		}
		return err
	}

	if len(fileData) == 0 {
		return nil
	}

	return json.Unmarshal(fileData, data)
}

func isPointer(v any) bool {
	return reflect.ValueOf(v).Kind() == reflect.Pointer
}
