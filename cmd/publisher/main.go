package main

import (
	"flag"

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
		linesCh := make(chan string)
		outputCh := make(chan file.ProcessingOutput)
		// yes, workaround for spawning goroutines in a for-loop
		filename := *output + "/" + filename

		go file.ProcessLineByLine(filename, linesCh, outputCh)
		for line := range linesCh {
			log.Info(line)
		}
		out := <-outputCh
		log.Info("dkdkdkdkdkdkdkkdkdkdkdkdk", out.Filename)
		if out.Err != nil {
			log.Fatal(out.Err)
		}

		// if err := <-readerrCh; err != nil {
		// 	log.Fatal(err)
		// }
	}
}
