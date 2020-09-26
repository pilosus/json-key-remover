package remover

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"sort"
	"testing"
)

func TestFileExists(t *testing.T) {
	var tests = []struct {
		path   string
		exists bool
	}{
		{"../../test/data/minimal.json", true},
		{"../../test/data/test.json", true},
		{"../../test/data", false},
		{"../../test/data/no_such_file.json", false},
	}

	for _, testCase := range tests {
		testName := fmt.Sprintf("%s", testCase.path)
		t.Run(testName, func(t *testing.T) {
			result := FileExists(testCase.path)

			if result != testCase.exists {
				t.Errorf("got %v, expected %v", result, testCase.exists)
			}
		})
	}
}

func TestWriteJSONFile(t *testing.T) {
	inputData := map[string]interface{}{
		"code":  200,
		"value": []map[string]interface{}{{"facebook": true, "twitter": false}, {"facebook": false, "twitter": true}},
	}

	tempFile, tempErr := ioutil.TempFile("", "test.*.json")
	if tempErr != nil {
		log.Fatal(tempErr)
	}
	defer os.Remove(tempFile.Name())

	var tests = []struct {
		desc string
		path string
		data map[string]interface{}
	}{
		{"Nested map", tempFile.Name(), inputData},
		{"Empty map", tempFile.Name(), map[string]interface{}{}},
	}

	for _, testCase := range tests {
		testName := fmt.Sprintf("%s", testCase.desc)
		t.Run(testName, func(t *testing.T) {
			WriteJSONFile(testCase.path, testCase.data)

			fileStat, fileErr := os.Stat(testCase.path)
			if fileErr != nil {
				log.Fatal(fileErr)
			}
			// get the fileSize
			fileSize := fileStat.Size()

			if fileSize == 0 {
				t.Errorf("got %v bytes, expected positive value", fileSize)
			}
		})
	}
}

func TestParseJSONFileSuccess(t *testing.T) {
	var tests = []struct {
		path string
		keys []string
	}{
		{"../../test/data/minimal.json", []string{"users"}},
		{"../../test/data/test.json", []string{"non-users", "users"}},
	}

	for _, testCase := range tests {
		testName := fmt.Sprintf("%s", testCase.path)
		t.Run(testName, func(t *testing.T) {
			actualMap := ParseJSONFile(testCase.path)

			actualKeys := []string{}
			for key := range actualMap {
				actualKeys = append(actualKeys, key)
			}
			sort.Strings(actualKeys)

			if !reflect.DeepEqual(actualKeys, testCase.keys) {
				t.Errorf("got %v, expected %v", actualKeys, testCase.keys)
			}
		})
	}
}

// See https://talks.golang.org/2014/testing.slide#23
func TestParseJSONFileFail(t *testing.T) {
	path := "no_such_file.json"

	if os.Getenv("BE_CRASHER") == "1" {
		ParseJSONFile(path)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestParseJSONFileFail")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, expected exit status 1", err)
}

func TestDeleteKey(t *testing.T) {
	inputData1 := map[string]interface{}{
		"code":  200,
		"value": []string{"one", "two"},
	}
	expectedData2 := map[string]interface{}{
		"code": 200,
	}
	inputData3 := map[string]interface{}{
		"code":  200,
		"value": map[string]interface{}{"facebook": true, "twitter": false},
	}
	expectedData3 := map[string]interface{}{
		"code":  200,
		"value": map[string]interface{}{"twitter": false},
	}
	inputData4 := map[string]interface{}{
		"code":  200,
		"value": []map[string]interface{}{{"facebook": true, "twitter": false}, {"facebook": false, "twitter": true}},
	}
	expectedData4 := map[string]interface{}{
		"code":  200,
		"value": []map[string]interface{}{{"twitter": false}, {"twitter": true}},
	}

	var tests = []struct {
		desc         string
		inputDict    map[string]interface{}
		keyToDelete  string
		exceptedDict map[string]interface{}
	}{
		{"Delete non existent key", inputData1, "no_such_key", inputData1},
		{"Delete root key", inputData1, "value", expectedData2},
		{"Delete key in nested map", inputData3, "facebook", expectedData3},
		{"Delete key in list of nested maps", inputData4, "facebook", expectedData4},
	}

	for _, testCase := range tests {
		testName := fmt.Sprintf("%s", testCase.desc)
		t.Run(testName, func(t *testing.T) {
			jsonParsed := reflect.ValueOf(testCase.inputDict)
			actualResult := DeleteKey(testCase.keyToDelete, jsonParsed)
			actualDict := actualResult.Interface().(map[string]interface{})

			if !reflect.DeepEqual(actualDict, testCase.exceptedDict) {
				t.Errorf("got %v, expected %v", actualDict, testCase.exceptedDict)
			}
		})
	}
}
