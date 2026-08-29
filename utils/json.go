package utils

import (
	"bufio"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"os"
	"path/filepath"
	"reflect"
)

var notPointer = errors.New("must be a pointer")

func WriteJSON(path string, val any) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp(dir, "json-tmp-*")
	if err != nil {
		return err
	}

	tmpName := tmpFile.Name()
	var success bool
	defer func() {
		_ = tmpFile.Close()
		if !success {
			_ = os.Remove(tmpName)
		}
	}()

	bufWriter := bufio.NewWriter(tmpFile)
	if err := json.MarshalWrite(bufWriter, val, jsontext.WithIndent("  ")); err != nil {
		return err
	}

	if err := bufWriter.Flush(); err != nil {
		return err
	}

	if err := tmpFile.Sync(); err != nil {
		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpName, path); err != nil {
		return err
	}

	success = true
	return nil
}

func ReadJSON(path string, data any) error {
	if !isPointer(data) {
		return notPointer
	}

	fileData, err := os.ReadFile(path)
	if err != nil {
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
