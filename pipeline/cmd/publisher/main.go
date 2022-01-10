package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/adjust/rmq/v4"
	"github.com/jasonlvhit/gocron"
	"github.com/valerykalashnikov/streaming-pipeline/pipeline/db"
	"github.com/valerykalashnikov/streaming-pipeline/pipeline/file"
	"github.com/valerykalashnikov/streaming-pipeline/pipeline/log"
)

func main() {
	output := flag.String("out", "/tmp/fileemitter", "The folder to put generated files")
	forceScan := flag.Bool("force-scan", false, "!!!Use on your own risk!!! This value is used to remove the state from redis and start scanning from the beginning.")
	daemonize := flag.Bool("d", false, "This value is used to process files and then daemonize the process to rescan the folder once an hour")
	period := flag.Int64("period", 10, "This value is used to set up the period of scanning")
	flag.Parse()

	if *daemonize {
		log.Info(fmt.Sprintf("Files processing will be running then once a %d minutes", *period))
		gocron.Every(uint64(*period)).Minute().Do(func() {
			forceScan := false
			publishing(output, &forceScan)
		})
		<-gocron.Start()
	} else {
		publishing(output, forceScan)
	}

}

func publishing(output *string, forceScan *bool) {

	redisAddr := db.GetRedisConnection()

	connection, err := rmq.OpenConnection("publisher", "tcp", redisAddr, 2, nil)
	if err != nil {
		log.Fatal(err)
	}

	data, err := connection.OpenQueue("data")
	if err != nil {
		log.Fatal(err)
	}

	fileList, err := file.IOReadDir(*output)
	if err != nil {
		log.Error("Unable to read file list to process", err.Error())
	}

	var ctx = context.Background()

	rdb := NewRedisClient(redisAddr)

	if *forceScan {
		err := rdb.RemoveState(ctx)
		if err != nil {
			log.Error("unable to remove state,", err)
		}
		log.Info("!!!processing files with --force-scan!!!")
	}

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
			if err := data.Publish(scanner.Text()); err != nil {
				log.Info("Failed to publish: %s", err)
			}
		}
		err = rdb.AddFilename(ctx, filename)
		if err != nil {
			log.Error("error while adding filename to redis %v", err)
		}
	}
	log.Info("Successfully processed %v", filesToProcess)
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
