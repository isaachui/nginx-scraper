//Plan: Write a program that will open up the nginx access logs and parse out
//each log to a statsd compatible format
// Things to do:
//   Attempt file open
//   Go to end of the file. Every 5 seconds, read to end of the file and parse data
//   Append data to log

package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	fileLocation := "testfile"
	logFile, logFileErr := os.Open(fileLocation)
	if logFileErr != nil {
		log.Fatal(logFileErr)
	}
	defer logFile.Close()

	for {
		fmt.Println("sleeping 5 seconds")
		time.Sleep(5 * time.Second)
	}

}

func parseNginxLog(logLine string) string {
	return ""
}
