package parsenginx

import (
	"regexp"
	"strings"
)

type nginxParser struct {

	//logFormat will store the nginx log_format string in one line
	logFormat      string
	reference      []string
	validVariables map[string]int //key stores nginx variable; value stores reference index
}

func NewNginxParser(format string) *nginxParser {
	nP := new(nginxParser)
	nP.logFormat = format
	nP.validVariables = make(map[string]int)

	var referenceArray []string
	var referenceElement string
	var isVariable bool
	arrayIndex := 0 //stores the array index of nginxParser.reference
	i := 0
	for i < len(format) {
		referenceElement = ""
		if string(format[i]) == "$" {
			isVariable = true
			referenceElement += string(format[i])
			i++

			//convert format[i] which is a byte into rune
			byteArray := make([]byte, 0)
			byteArray = append(byteArray, format[i])
			isValid, _ := regexp.Match("^[A-Za-z_]", byteArray)
			//check if oneRune is alphanumeric or is an _
			for isValid {
				referenceElement += string(format[i])
				i++

				byteArray = make([]byte, 0)
				byteArray = append(byteArray, format[i])
				isValid, _ = regexp.Match("^[A-Za-z_]", byteArray) // check if its valid
			}

		} else {
			//encountered a non variable
			isVariable = false
			for i < len(format) && string(format[i]) != "$" {
				referenceElement += string(format[i])
				i++
			}
		}
		if isVariable {
			isVariable = false
			nP.validVariables[referenceElement] = arrayIndex
		}
		referenceArray = append(referenceArray, referenceElement)
		arrayIndex++

	}
	nP.reference = referenceArray

	return nP
}
func (p nginxParser) isValidVariable(input string) bool {
	_, ok := p.validVariables[input]
	return ok
}

func (p nginxParser) ParseLine(input, search string) string {
	if !p.isValidVariable(search) {
		return "ERR invalid search variable"
	}
	//	i := 0 //i will contain the index of string
	//j := 0 //j will contain the index of which reference we're on
	mutableInputString := input
	indexOfSearch := 0
	foundData := false
	for _, value := range p.reference {
		if value == search {
			foundData = true
		}
		if !strings.Contains(value, "$") {
			indexOfSearch = strings.Index(mutableInputString, value)
			if foundData {
				return mutableInputString[:indexOfSearch]
			}
			mutableInputString = mutableInputString[indexOfSearch+len(value):]
		}
	}

	//looking for status code
	return "hellO"
}

const defaultNginxLogFormat = `$remote_addr - $http_x_forwarded_for - $http_x_realip - [$time_local]  $scheme $http_x_forwarded_proto $x_forwarded_proto_or_scheme "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"`

//func main() {
//
//	parser := NewNginxParser(defaultNginxLogFormat)
//
//	testLine := "50.112.166.232 - 50.112.166.232, 192.33.28.238, 50.112.166.232,127.0.0.1 - - - [03/Aug/2015:08:34:40 +0000]  http https,http https,http \"GET /api/v1/user HTTP/1.1\" 200 3350 \"https://release.dollarshaveclub.com/login\" \"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:39.0) Gecko/20100101 Firefox/39.0\""
//	//	splitTestLine := strings.Split(testLine, " - ")
//	fmt.Println(parser.parseLine(testLine, "$status"))
//
//}
