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
	daemonize := flag.Bool("d", false, "This value is used to generate files infinitely")
	period := flag.Int64("period", 2, "This value is used to set up the period of waiting until generating new file")
	flag.Parse()

	for {
		waitingPeriod := int(*period)
		time.Sleep(time.Duration(waitingPeriod) * time.Minute)
		minSizeVal, err := units.FromHumanSize(*minSize)
		if err != nil {
			log.Error("Cannot convert min-size parameter", err.Error)
		}
		maxSizeVal, err := units.FromHumanSize(*maxSize)
		if err != nil {
			log.Error("Cannot convert max-size parameter", err.Error)
		}

		if minSizeVal > maxSizeVal {
			log.Error("min-size cannot be higher than maxsize")
		}

		err = os.MkdirAll(*output, os.ModePerm)
		if err != nil {
			log.Error("Cannot create directory for the files", err.Error)
		}

		rand.Seed(time.Now().UnixNano())
		fileSize := rand.Intn(int(maxSizeVal)-int(minSizeVal)) + int(minSizeVal)

		filePath := buildFullPath(*output)
		log.Info("Creating %s size of %s \n", filePath, units.HumanSize(float64(fileSize)))
		err = file.Generate(filePath, fileSize)
		if err != nil {
			log.Error("Error while generating new file", err)
		}
		log.Info("File succesfully generated")

		if !(*daemonize) {
			break
		}
	}

}

func buildFullPath(outputDir string) (filepath string) {
	currentTime := time.Now()
	return outputDir + "/" + petname.Generate(2, "-") + "_" + currentTime.Format("2006-01-02_15_04_05")
}
