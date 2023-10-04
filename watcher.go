package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
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

	i := 0
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var logLine LogLine

		line := scanner.Text()

		if err := json.Unmarshal([]byte(line), &logLine); err != nil {
			fmt.Printf("Error parsing JSON at line %d: %v\n", i, err)
			continue
		}

		msg := logLine.Msg[0]
		if msg.Type == "/vm.m_addpkg" {

			path := msg.Package.Path
			// do not add duplicates
			_, ok := pkgMap[path]

			if !ok {
				pkgMap[path] = logLine
			}
			i++
		}

	}
}

func writePkg(pkgPath string) {
	fmt.Println("")

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
