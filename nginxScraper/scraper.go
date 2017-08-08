//Plan: Write a program that will open up the nginx access logs and parse out
//each log to a statsd compatible format
// Things to do:
//   Attempt file open
//   Go to end of the file. Every 5 seconds, read to end of the file and parse data
//   Append data to log
//   Account for logrotate

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"nginxscraper/parsenginx"
	"os"
	"strconv"
	"strings"
	"time"
)

const combinedNginxLogFormat = `$remote_addr - $http_x_forwarded_for - $http_x_realip - [$time_local]  $scheme $http_x_forwarded_proto $x_forwarded_proto_or_scheme "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent";`

func main() {

	testing := true
	//switch directory to /var to read logs
	os.Chdir("/var")
	readFileLocation := "var/log/nginx/access.log"
	writeFileLocation := "var/log/stats.log"

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

	// This is needed for logrotate checking
	logFileInfo, logFileInfoErr := logFile.Stat()
	if logFileInfoErr != nil {
		log.Fatal(logFileInfoErr)
	}

	//Logrotate values
	fileMoved := false
	fileTruncated := false

	//logLine will contain the lines from the file
	var logLine string
	var logReadErr error
	statsDMap := make(map[string]int)
	//statsDMap will contain key of error codes
	//statsDMmap will contain value of number of instances

	if testing != true {
		logFile.Seek(0, 2) //moves to the bottom of the file to start reading
	}

	logReader := bufio.NewReader(logFile)

	//create Parser
	np := parsenginx.NewNginxParser(combinedNginxLogFormat)
	for {

		//start reading at end of the file
		for {
			logLine, logReadErr = logReader.ReadString('\n')

			// #### Account for logrotate ######
			//check if file has been moved
			watchFileInfo, watchFileError := os.Stat(readFileLocation)
			if watchFileError != nil {
				// the file at readFileLocation does not exist anymore. Break and wait 5 seconds
				fmt.Println("file not found")
				break
			}

			//compares watchFile and logFile
			if os.SameFile(watchFileInfo, logFileInfo) != true {
				fileMoved = true
				//file has been moved finish reading, then handle after loop
			}

			currentPosition, _ := logFile.Seek(0, 1)
			if watchFileInfo.Size() < currentPosition {
				fileTruncated = true
				//file has been truncated. Finish readigng, then handle after loop
			}

			//once the EOF is reached, break the loop and sleep 5 seconds
			if logReadErr == io.EOF {
				break
			}

			//put status code into map
			httpStatus, httpStatusError := np.ParseLine(logLine, "$status")
			if httpStatusError != nil {
				fmt.Println(httpStatusError)
				continue
			}

			parsedHttpStatus := parseStatus(httpStatus)
			//check if parsedHttpStatus is valid if its not within 200-599
			if validStatus(parsedHttpStatus) == false {
				continue
			}
			numberOccurences, exists := statsDMap[parsedHttpStatus]
			if exists {
				statsDMap[parsedHttpStatus] = numberOccurences + 1
			} else {
				statsDMap[parsedHttpStatus] = 1
			}

			//put 50x routes in
			if parsedHttpStatus == "50x" {
				httpRequest, httpRequestError := np.ParseLine(logLine, "$request")
				if httpRequestError != nil {
					fmt.Println(httpRequestError)
				}
				parsedHttpRequest := parseRequest(httpRequest)
				numberOccurences, exists := statsDMap[parsedHttpRequest]
				if exists {
					statsDMap[parsedHttpRequest] = numberOccurences + 1
				} else {
					statsDMap[parsedHttpRequest] = 1
				}

			}

		}
		//  output map values to file
		for key, value := range statsDMap {
			statsFile.WriteString(key + ":" + strconv.Itoa(value) + "|s\n")
			fmt.Print(key + ":" + strconv.Itoa(value) + "|s\n")
		}
		//flush map
		statsDMap = make(map[string]int)

		//Account for file moved and truncated: reopen logfile
		if fileMoved {
			fileMoved = false
		}

		//Account for file truncated
		if fileTruncated {
			logFile.Seek(0, 2) // move to end of file
			fileTruncated = false
		}
		time.Sleep(5 * time.Second)
	}

}

func validStatus(status string) bool {
	statusCode, err := strconv.ParseInt(status, 0, 64)
	if err != nil {
		return false
	}
	if statusCode <= 200 || statusCode >= 600 {
		return false
	}
	return true
}

func parseStatus(status string) string {
	statusCode, _ := strconv.ParseInt(status, 0, 64)
	switch {
	case statusCode < 300:
		return "20x"
	case statusCode < 400:
		return "30x"
	case statusCode < 500:
		return "40x"
	case statusCode < 600:
		return "50x"
	}
	return "ERR"
}

func parseRequest(request string) string {
	startIndex := strings.Index(request, " ")
	endIndex := strings.LastIndex(request, " ")
	return request[startIndex+1 : endIndex]
}
