package main

import (
	"flag"
	"fmt"
)

func main() {
	minSize := flag.String("min-size", "2MB", "minimal size of the generated file")
	maxSize := flag.String("max-size", "10MB", "Maximum size of the generated file")
	output := flag.String("out", "/tmp/fileemitter", "The folder to put generated files")
	flag.Parse()

	fmt.Printf("minSize: %s, maxSize: %s, out: %s\n", *minSize, *maxSize, *output)
}
