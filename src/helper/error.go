package helper

import (
	"regexp"
	"strings"
)

type CustomError struct {
	Error string `json:"error"`
}

func NewCustomError(err error) CustomError {
	return CustomError{
		Error: err.Error(),
	}
}

func ParseMongoError(err string) string {
	indexOfKeys := strings.Index(err, ": {") + 1
	regexp := regexp.MustCompile(` (\w+)`)
	matches := regexp.FindAllString(err[indexOfKeys:len(err) - 1], -1)
	message := "contact with "
	for i, match := range matches {
		if i == len(matches) - 1 {
			message += strings.Trim(match, " ")
		} else {
			message += strings.Trim(match, " ") + ", "
		}
	}
	return message + " already exists"
}