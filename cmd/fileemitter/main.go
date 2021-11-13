package main

import (
	"flag"
	"math/rand"
	"os"
	"time"

	units "github.com/docker/go-units"
	petname "github.com/dustinkirkland/golang-petname"

	"github.com/valerykalashnikov/streaming-pipeline/file"
	"github.com/valerykalashnikov/streaming-pipeline/log"
)

func main() {
	minSize := flag.String("min-size", "2MB", "minimal size of the generated file")
	maxSize := flag.String("max-size", "10MB", "Maximum size of the generated file")
	output := flag.String("out", "/tmp/fileemitter", "The folder to put generated files")
	flag.Parse()

	minSizeVal, err := units.FromHumanSize(*minSize)
	if err != nil {
		log.Error("Cannot convert min-size parameter", err.Error)
	}
	maxSizeVal, err := units.FromHumanSize(*maxSize)
	if err != nil {
		log.Error("Cannot convert max-size parameter", err.Error)
	}

	err = os.MkdirAll(*output, os.ModePerm)
	if err != nil {
		log.Error("Cannot create directory for the files", err.Error)
	}

	for {
		rand.Seed(time.Now().UnixNano())
		fileSize := rand.Intn(int(maxSizeVal)-int(minSizeVal)) + int(minSizeVal)
		log.Info("Creating file size of %d bytes \n", fileSize)

		filePath := buildFullPath(*output)
		err := file.Generate(filePath, fileSize)
		if err != nil {
			log.Error("Error while generating new file", err)
		}
		log.Info("File succesfully generated")
	}
}

func buildFullPath(outputDir string) (filepath string) {
	currentTime := time.Now()
	return outputDir + "/" + petname.Generate(2, "-") + "_" + currentTime.Format("2006-01-02_15_04_05")
}
