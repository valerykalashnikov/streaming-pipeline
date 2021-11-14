package main

import (
	"bufio"
	"context"
	"flag"
	"os"

	"github.com/jasonlvhit/gocron"
	"github.com/valerykalashnikov/streaming-pipeline/file"
	"github.com/valerykalashnikov/streaming-pipeline/log"
)

func main() {
	output := flag.String("out", "/tmp/fileemitter", "The folder to put generated files")
	flag.Parse()

	publishing(output)

	gocron.Every(1).Hour().Do(func() { publishing(output) })
	<-gocron.Start()
}

func publishing(output *string) {
	fileList, err := file.IOReadDir(*output)
	if err != nil {
		log.Error("Unable to read file list to process", err.Error())
	}

	var ctx = context.Background()
	rdb := NewRedisClient()

	alreadyProcessedList, err := rdb.GetProcessedFilesList(ctx)
	if err != nil {
		log.Fatal(err)
	}

	filesToProcess := difference(fileList, alreadyProcessedList)
	log.Info("start processing %v", filesToProcess)

	for _, filename := range filesToProcess {
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
		err = rdb.AddFilename(ctx, filename)
	}
	log.Info("successfully processed %v", filesToProcess)
}

func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
