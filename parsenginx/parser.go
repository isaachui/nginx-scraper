package parsenginx

import (
	"errors"
	"regexp"
	"strings"
)

type nginxParser struct {

	//logFormat will store the nginx log_format string in one line
	logFormat      string
	reference      []string
	validVariables map[string]int //key stores nginx variable; value stores reference index
}

//defaultNginxLogFormat contains a default nginx log_format to start parsing
const defaultNginxLogFormat = `$remote_addr - $http_x_forwarded_for - $http_x_realip - [$time_local]  $scheme $http_x_forwarded_proto $x_forwarded_proto_or_scheme "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"`

//NewDefaultParser creates a NewNginxParser with defaultNginxLogFormat as the parameter
func NewDefaultParser() *nginxParser {
	return NewNginxParser(defaultNginxLogFormat)
}

//NewNginxParser creats a nginxParser with all values filled in depending on the format string
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

				if i == len(format) {
					break
				}
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

//isValidVariable wil check if the map validVariables contains the input string.
func (p nginxParser) isValidVariable(input string) bool {
	_, ok := p.validVariables[input]
	return ok
}

//ParseLine will take in an nginx variable as a string "$example" and return the parameter as a string
func (p nginxParser) ParseLine(input, search string) (string, error) {
	if !p.isValidVariable(search) {
		return "", errors.New("not a valid search string")
	}
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
				return mutableInputString[:indexOfSearch], nil
			}
			mutableInputString = mutableInputString[indexOfSearch+len(value):]
		}

	}
	return "", errors.New("not found")
}
