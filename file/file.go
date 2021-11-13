package file

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Generate generates file with given filename with of a given size with a content like "123 12345".
// sizeBytes is intentionally in int because we don't have any plans to generate Terabytes of files.
func Generate(filename string, sizeBytes int) error {
	isExist, err := fileExists(filename)
	if err != nil {
		return err
	}
	if isExist {
		return nil
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	for {
		consumerId := generateRand(1, 10000)
		resourceConsumed := generateRand(1, 100000)
		str := fmt.Sprintf("%d %d\n", consumerId, resourceConsumed)
		_, err = f.WriteString(str)

		sizeBytes = sizeBytes - len(str)
		if sizeBytes <= 0 {
			break
		}
	}

	return err
}

func generateRand(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func fileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
