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
	"parsenginx"
	"strconv"
	"strings"
	"time"
)

const defaultNginxLogFormat = `$remote_addr - $http_x_forwarded_for - $http_x_realip - [$time_local]  $scheme $http_x_forwarded_proto $x_forwarded_proto_or_scheme "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent";`

func main() {

	testing := true
	readFileLocation := "sample.log"
	writeFileLocation := "stats.log"

	logFile, logFileErr := os.Open(readFileLocation)
	if logFileErr != nil {
		log.Fatal(logFileErr)
	}
	defer logFile.Close()

	statsFile, statsFileErr := os.OpenFile(writeFileLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if statsFileErr == os.ErrNotExist {
		if statsFileErr != nil {
			log.Fatal(statsFileErr)
		}
	}
	defer statsFile.Close()

	// might not need this #################
	logFileInfo, logFileInfoErr := logFile.Stat()
	if logFileInfoErr != nil {
		log.Fatal(logFileInfoErr)
	}
	fmt.Println(logFileInfo.Size())
	//#######################################

	var logLine string
	var logReadErr error
	statsDMap := make(map[string]int) //statsDMap will contain key of error codes
	//statsDMmap will contain value of number of instances

	if testing != true {
		logFile.Seek(0, 2) //moves to the bottom of the file to start reading
	}

	logReader := bufio.NewReader(logFile)

	//create Parser
	np := parsenginx.NewNginxParser(defaultNginxLogFormat)
	for {
		//start reading at end of line
		for {
			logLine, logReadErr = logReader.ReadString('\n')
			if logReadErr == io.EOF {
				break // once we reach the EOF break out of this forloop
			}

			//put status code into map
			httpstatus := parseStatus(np.ParseLine(logLine, "$status"))
			numberOccurences, exists := statsDMap[httpstatus]
			if exists {
				statsDMap[httpstatus] = numberOccurences + 1
			} else {
				statsDMap[httpstatus] = 1
			}
			//put 50X directories in
			if httpstatus == "50X" {
				httpRequest := parseRequest(np.ParseLine(logLine, "$request"))

				numberOccurences, exists := statsDMap[httpstatus]
				if exists {
					statsDMap[httpRequest] = numberOccurences + 1
				} else {
					statsDMap[httpRequest] = 1
				}

			}
			//put log parsing into map

		}
		// write map values to file
		for key, value := range statsDMap {
			statsFile.WriteString(key + ":" + strconv.Itoa(value) + "|s\n")
			fmt.Print(key + ":" + strconv.Itoa(value) + "|s\n")
		}
		//flush map
		statsDMap = make(map[string]int)
		fmt.Println("sleeping 5 seconds")
		time.Sleep(5 * time.Second)
	}

}

func parseStatus(status string) string {
	statusCode, _ := strconv.ParseInt(status, 0, 64)
	switch {
	case statusCode < 300:
		return "20X"
	case statusCode < 400:
		return "30X"
	case statusCode < 500:
		return "40X"
	case statusCode < 600:
		return "50X"
	}
	return "ERR"
}

func parseRequest(request string) string {
	startIndex := strings.Index(request, " ")
	endIndex := strings.LastIndex(request, " ")
	return request[startIndex+1 : endIndex]
}
