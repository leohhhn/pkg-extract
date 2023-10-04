package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

// see how many msgs there are
// create x go routines
// each will invoke jq
// get packages parsed out into an array

func processFile(fileName string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Print("Error opening file: ", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Could not close file %s properly.", file.Name())
		}
	}(file)

	pkgMap := make(map[string]LogLine) // path -> package

	scanner := bufio.NewScanner(file)

	for i := 0; scanner.Scan(); i++ {
		var logLine LogLine

		line := scanner.Text()

		if err := json.Unmarshal([]byte(line), &logLine); err != nil {
			fmt.Printf("Error parsing JSON at line %d: %v\n", i, err)
			continue
		}

		if logLine.Msg[0].Type == "/vm.m_addpkg" {
			path := logLine.Msg[0].Package.Path

			// do not add duplicates
			_, ok := pkgMap[path]
			if !ok {
				writePkg(logLine)
				pkgMap[path] = logLine // possibly not needed
			}
		}
	}
}

func writePkg(logLine LogLine) {
	msg := logLine.Msg[0]
	path := msg.Package.Path

	// do routines for each file write
	trimmedPath := strings.TrimLeft(path, "gno.land/")

	// write dirs needed to write package
	if err := os.MkdirAll("extracted/"+trimmedPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	metadata, err := logLine.MarshalMetadata()

	if err != nil {
		log.Fatal("Failed to marshal metadata: " + trimmedPath)
	}
	// write metadata
	err = os.WriteFile("extracted/"+trimmedPath+"/metadata.json", metadata, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// write files

}

func main() {

	entries, err := os.ReadDir("./logs")
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for _, e := range entries {
		wg.Add(1)
		go processFile("./logs/"+e.Name(), &wg)
	}

	wg.Wait()

}
