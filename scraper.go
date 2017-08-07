//Plan: Write a program that will open up the nginx access logs and parse out
//each log to a statsd compatible format
// Things to do:
//   Attempt file open
//   Go to end of the file. Every 5 seconds, read to end of the file and parse data
//   Append data to log

package main

import (
	"bufio"
	"fmt"
	"io"
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

	// might not need this #################
	logFileInfo, logFileInfoErr := logFile.Stat()
	if logFileInfoErr != nil {
		log.Fatal(logFileInfoErr)
	}
	fmt.Println(logFileInfo.Size())
	//#######################################

	var logLine string
	var logReadErr error

	logFile.Seek(0, 2) //moves to the bottom of the file to start reading

	logReader := bufio.NewReader(logFile)

	for {
		//start reading at end of line
		for {
			logLine, logReadErr = logReader.ReadString('\n')
			if logReadErr == io.EOF {
				break // once we reach the EOF break out of this forloop
			}
			fmt.Print("parse Logline here: ", logLine)

		}

		fmt.Println("sleeping 5 seconds")
		time.Sleep(5 * time.Second)
	}

}

func parseNginxLog(logLine string) string {
	return ""
}
