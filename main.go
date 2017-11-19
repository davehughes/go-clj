package main

import (
	"fmt"
	"os"

	"github.com/davehughes/lparse/lang"
)

func main() {
	parseFiles := os.Args[1:]

	for _, file := range parseFiles {
		fmt.Printf("Parsing file: %v\n", file)
		parseResult, err := lang.ParseFile(file)
		if err != nil {
			fmt.Printf("Encountered error: %v\n", err)
		}
		fmt.Printf("Successfully parsed, result: %v\n", parseResult)
	}
}
