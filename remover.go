package main

import "os"

// FileExists returns true if given path exists and a file, otherwise false
func FileExists(path string) bool {
	stat, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}
	return !stat.IsDir()
}


//func ParseJSONFile(path string) {
//
//}