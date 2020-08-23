package main

import (
	"flag"
	"fmt"
	"github.com/pilosus/json-key-remover/pkg/remover"
	"os"
	"reflect"
)

var (
	version   string
	timestamp string
)

func main() {
	jsonFileFrom := flag.String("from", "", "Path to a JSON text file to read from")
	jsonFileTo := flag.String("to", "", "Path to a JSON text file to write data to")
	keyToRemovePtr := flag.String("key", "", "Key to remove recursively from the JSON")
	flag.Parse()

	if *jsonFileFrom == "" || *jsonFileTo == "" || *keyToRemovePtr == "" {
		fmt.Printf("json-key-remover (build: %s %s)\n", version, timestamp)
		flag.PrintDefaults()
		os.Exit(1)
	}

	// BUG (vrs) tilda for linux home directory not expanded properly
	if !remover.FileExists(*jsonFileFrom) {
		fmt.Printf("File %s not found/n", *jsonFileFrom)
		os.Exit(1)
	}

	fmt.Printf("Parsing file: %s\n", *jsonFileFrom)
	parsedJSON := remover.ParseJSONFile(*jsonFileFrom)
	resultJSON := remover.DeleteKey(*keyToRemovePtr, reflect.ValueOf(parsedJSON))

	fmt.Printf("Writing result to file: %s\n", *jsonFileTo)
	remover.WriteJSONFile(*jsonFileTo, resultJSON.Interface().(map[string]interface{}))
}
