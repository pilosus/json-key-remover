package remover

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

// FileExists returns true if given path exists and a file, otherwise false
func FileExists(path string) bool {
	stat, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}
	return !stat.IsDir()
}


// ParseJSONFile unmarshal data from JSON file to a map
func ParseJSONFile(path string) map[string]interface{} {
	jsonFile, errFileOpen := os.Open(path)

	if errFileOpen != nil {
		fmt.Println(errFileOpen)
		os.Exit(1)
	}

	defer jsonFile.Close()

	byteContent, _ := ioutil.ReadAll(jsonFile)

	var jsonParsed map[string]interface{}
	errUnmarshal := json.Unmarshal([]byte(byteContent), &jsonParsed)

	if errUnmarshal != nil {
		fmt.Printf("Error: %s\n", errUnmarshal)
		os.Exit(1)
	}

	return jsonParsed
}


// WriteJSONFile writes map to a file
func WriteJSONFile(path string, data map[string]interface{}) {
	encoded, errMarshal := json.Marshal(data)

	if errMarshal != nil {
		fmt.Printf("Error: %s\n", errMarshal)
		os.Exit(1)
	}

	errWrite := ioutil.WriteFile(path, encoded, 0644)

	if errWrite != nil {
		fmt.Printf("Error: %s\n", errWrite)
		os.Exit(1)
	}
}


// DeleteKey deletes a key from the map recursively
func DeleteKey(keyToDelete string, val reflect.Value) reflect.Value {
	// Indirect through pointers and interfaces
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			DeleteKey(keyToDelete, val.Index(i))
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if k.String() == keyToDelete {
				delete(val.Interface().(map[string]interface{}), k.String())
				continue
			}
			DeleteKey(keyToDelete, val.MapIndex(k))
		}
	default:
		// Do we need that case?
	}
	return val
}
