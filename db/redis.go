package db

import (
	"fmt"
	"os"
)

func GetRedisConnection() (addr string) {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	return fmt.Sprintf("%s:%s", host, port)
}
