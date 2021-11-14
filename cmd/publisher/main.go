package main

import (
	"bufio"
	"flag"
	"os"

	"github.com/valerykalashnikov/streaming-pipeline/file"
	"github.com/valerykalashnikov/streaming-pipeline/log"
)

func main() {
	output := flag.String("out", "/tmp/fileemitter", "The folder to put generated files")
	fileList, err := file.IOReadDir(*output)
	if err != nil {
		log.Error("Unable to read file list to process", err.Error())
	}

	for _, filename := range fileList {
		filepath := *output + "/" + filename
		file, err := os.Open(filepath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			log.Info(scanner.Text())
		}
		log.Info(filename)
	}
}
