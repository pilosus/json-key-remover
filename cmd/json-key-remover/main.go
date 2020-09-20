package main

import (
	"flag"
	"github.com/pilosus/json-key-remover/pkg/remover"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
)

var (
	version   string
	timestamp string
)

func initLogging(debug bool) {
	logLevel := log.InfoLevel
	if debug {
		logLevel = log.DebugLevel
	}
	log.SetOutput(os.Stdout)
	log.SetLevel(logLevel)
}

func main() {
	jsonFileFrom := flag.String("from", "", "Path to a JSON text file to read from")
	jsonFileTo := flag.String("to", "", "Path to a JSON text file to write data to")
	keyToRemovePtr := flag.String("key", "", "Key to remove recursively from the JSON")
	debugMode := flag.Bool("debug", false, "Show debugging info")
	flag.Parse()

	initLogging(*debugMode)

	if *jsonFileFrom == "" || *jsonFileTo == "" || *keyToRemovePtr == "" {
		log.Info("json-key-remover (build: ", version, " ", timestamp, ")")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// BUG (vrs) tilda for linux home directory not expanded properly
	if !remover.FileExists(*jsonFileFrom) {
		log.Error("File not found: ", *jsonFileFrom)
		os.Exit(1)
	}

	log.Debug("Parsing file: ", *jsonFileFrom)
	parsedJSON := remover.ParseJSONFile(*jsonFileFrom)
	resultJSON := remover.DeleteKey(*keyToRemovePtr, reflect.ValueOf(parsedJSON))

	log.Debug("Writing result to file: ", *jsonFileTo)
	remover.WriteJSONFile(*jsonFileTo, resultJSON.Interface().(map[string]interface{}))
}
