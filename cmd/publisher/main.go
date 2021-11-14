package main

import (
	"bufio"
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/adjust/rmq/v4"
	"github.com/jasonlvhit/gocron"
	"github.com/valerykalashnikov/streaming-pipeline/file"
	"github.com/valerykalashnikov/streaming-pipeline/log"
)

func main() {
	output := flag.String("out", "/tmp/fileemitter", "The folder to put generated files")
	forceScan := flag.Bool("force-scan", false, "!!!Use on your own risk!!! This value is used to remove the state from redis and start scanning from the beginning.")
	daemonize := flag.Bool("d", false, "This value is used to process files and then daemonize the process to rescan the folder once an hour")
	flag.Parse()
	go func() {
		publishing(output, *&forceScan)

		if *daemonize == true {
			log.Info("Files processing will be running then once an hour")
			gocron.Every(1).Hour().Do(func() {
				forceScan := false
				publishing(output, &forceScan)
			})
			<-gocron.Start()
		}
	}()

	connection, err := rmq.OpenConnection("dashboard", "tcp", "localhost:6379", 2, nil)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/overview", NewDashboardHandler(connection))
	log.Info("Dashboard is rendered on http://localhost:3333/overview")
	if err := http.ListenAndServe(":3333", nil); err != nil {
		panic(err)
	}
}

func publishing(output *string, forceScan *bool) {

	connection, err := rmq.OpenConnection("publisher", "tcp", "localhost:6379", 2, nil)
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
	rdb := NewRedisClient()

	if *forceScan == true {
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
