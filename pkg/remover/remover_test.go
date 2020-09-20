package remover

import (
	"fmt"
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
