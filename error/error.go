package error

import (
	"fmt"
	"strings"
)

type ChakiError struct {
	ClientName string
	StatusCode int
	RawBody    []byte
	ParsedBody interface{}
}

func (e *ChakiError) Error() string {
	msg := fmt.Sprintf("Error on client %s (Status %d)", e.ClientName, e.StatusCode)
	if details := e.extractErrorDetails(); details != "" {
		msg += ": " + details
	}

	return msg
}

func (e *ChakiError) Status() int {
	return e.StatusCode
}

type RandomError interface {
	Status() int
}

func (e *ChakiError) extractErrorDetails() string {
	var details []string

	var extract func(interface{})
	extract = func(v interface{}) {
		switch value := v.(type) {
		case string:
			details = append(details, strings.TrimSpace(value))
		case map[string]interface{}:
			for _, v := range value {
				extract(v)
			}
		case []interface{}:
			for _, v := range value {
				extract(v)
			}
		}
	}

	extract(e.ParsedBody)

	if len(details) == 0 && len(e.RawBody) > 0 {
		return strings.TrimSpace(string(e.RawBody))
	}

	return strings.Join(details, "; ")
}
