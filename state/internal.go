package state

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func read[T any](path string) []*T {
	data, err := os.ReadFile(path)
	switch {
	case os.IsNotExist(err):
		return []*T{}
	case err != nil:
		panic(err)
	}
	var out []*T
	if err := json.Unmarshal(data, &out); err != nil {
		panic(err)
	}
	return out
}

func write[T any](path string, data []*T) {
	if len(data) == 0 {
		err := os.Remove(path)
		switch {
		case os.IsNotExist(err):
			// Ignore if the file does not exist
		case err != nil:
			panic(err)
		}
		return
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(path, jsonData, 0644); err != nil {
		panic(err)
	}
}
