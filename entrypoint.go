// https://tutorialedge.net/golang/parsing-json-with-golang/
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	jsonFilePathPtr := flag.String("path", "", "Path to a JSON text file")
	flag.Parse()

	if *jsonFilePathPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// BUG (vrs) tilda for linux home directory not expanded properly
	if !FileExists(*jsonFilePathPtr) {
		fmt.Printf("File %s not found/n", *jsonFilePathPtr)
		os.Exit(1)
	}

	fmt.Printf("Parsing file: %s\n", *jsonFilePathPtr)
}
